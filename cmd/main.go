package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/segmentio/ksuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"portfolio/pkg"
	"portfolio/server"
	"runtime"
	"time"
)

// It will add the specified files.
//go:embed dist/favicon.svg dist/index.html
// It will add all non-hidden file in images, css, and js.
//go:embed dist/*
var static embed.FS

func main() {
	dbPath := flag.String("db", "portfolio.db", "path to db file")
	doHistoryNow := flag.Bool("runHistory", false, "run history portfolio cron")
	doMigration := flag.Bool("runMigration", false, "run db migration")

	flag.Parse()

	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if *doHistoryNow {
		logValue(db)
	}

	if *doMigration {
		if err:= migrate(db);err!=nil{
			log.Fatal(err)
		}
	}

	s1 := gocron.NewScheduler(time.Local)
	_, err = s1.Every(1).Day().At("17:00").Do(logValue, db)
	if err != nil {
		log.Fatal(err)
	}
	s1.StartAsync()
	s := server.Server{DB: db}

	app := fiber.New()
	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080, http://localhost:5000",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, X-CSRF-Token, X-Requested-With, Accept-Language, Cache-Control, User-Agent",
		AllowMethods: "GET, POST, OPTIONS , PUT, DELETE",
		ExposeHeaders: "Link",
		AllowCredentials: true,
	}))


	subFS, _ := fs.Sub(static, "dist")
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   http.FS(subFS),
		Browse: true,
	}))

	v1 := app.Group("/api").Group("/v1")
	v1.Get("/portfolios/", s.GetPortfolios)
	v1.Get("/portfolios/:portfolio_id", s.GetPortfolio)

	log.Fatal(app.Listen(":5000"))
}



func logValue(db *gorm.DB) {
	p := pkg.Portfolio{}
	if err := db.Preload(clause.Associations).First(&p).Error; err != nil {
		log.Println("error fetching portfolio from DB")
		return
	}
	if err := p.UpdateTotalValue(); err != nil {
		PriceRecovery(context.Background(), &p)
	}
	response, err := json.Marshal(p)
	if err != nil {
		log.Printf("error saving portfolio %s value. %.2f\n", p.ID, p.TotalValue)
	}
	hist := pkg.PortfolioHistory{
		ID:          ksuid.New().String(),
		Date:        time.Now(),
		Value:       p.TotalValue,
		SnapShot:    string(response),
		PortfolioID: p.ID,
	}
	db.Create(&hist)
	log.Println(p.String())
}



func PriceRecovery(ctx context.Context, p *pkg.Portfolio) {
	for {
		select {
		case <-time.After(5 * time.Minute):
			fmt.Printf("calling price recovery %s \n", time.Now().Format(time.Kitchen))
			if err := p.UpdateTotalValue(); err == nil {
				fmt.Println("price recovery complete")
				break
			}
		case <-ctx.Done():
			fmt.Println("ctx halted")
		}
	}
}

func migrate(db *gorm.DB) error {
	p := pkg.Portfolio{
		ID:           ksuid.New().String(),
		Name:         "PortfolioName",
		BaseCurrency: "CAD",
		TotalValue:   0,
		Items:        nil,
		History:      nil,
	}
	//Add portfolio items
	//p.Items = append(p.Items, pkg.PortfolioItem{
	//	PortfolioID: p.ID,
	//	Symbol:      "ETH-CAD",
	//	Quantity:    0,
	//})


	if err := db.AutoMigrate(&pkg.Portfolio{}, &pkg.PortfolioItem{}, &pkg.PortfolioHistory{}); err != nil {
		return err
	}
	db.Create(&p)
	os.Exit(1)
	return nil
}


func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}