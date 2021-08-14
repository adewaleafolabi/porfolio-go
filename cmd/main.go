package main

import (
	"embed"
	"flag"
	"github.com/pkg/browser"
	"io/fs"
	"net/http"
	"portfolio/cron"
	"portfolio/pkg"
	"portfolio/pricing"
	"portfolio/server"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// It will add the specified files.
//go:embed dist/favicon.svg dist/index.html
// It will add all non-hidden file in images, css, and js.
//go:embed dist/*
var static embed.FS

func main() {
	prodLogger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(prodLogger)

	sugar := prodLogger.Sugar()

	dbPath := flag.String("db", "portfolio.db", "path to db file")
	doHistoryNow := flag.Bool("runHistory", false, "run history portfolio cron")
	doMigration := flag.Bool("runMigration", false, "run db migration")

	flag.Parse()

	conf := server.Config{}
	if err := envconfig.Process("", &conf); err != nil {
		sugar.Fatal(err)
	}
	conf.CoinAPIBaseURL = strings.TrimRight(conf.CoinAPIBaseURL, "/")

	pm := pricing.NewPriceManager(conf.CoinAPIToken, conf.CoinAPIBaseURL)

	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})
	if err != nil {
		sugar.Fatalw("failed to connect database", "db", *dbPath, "err", err)
	}

	if *doHistoryNow {
		cron.LogValue(db, pm, sugar)
	}

	if *doMigration {
		sugar.Info("Starting DB migration")

		if err := migrate(db); err != nil {
			sugar.Fatal("DB migration failed", err)
		}
		sugar.Info("DB migration successful")
	}

	s1 := gocron.NewScheduler(time.Local)
	_, err = s1.Every(1).Day().At(conf.DailyCronTime).Do(cron.LogValue, db, pm, sugar)
	if err != nil {
		sugar.Fatal(err)
	}
	s1.StartAsync()

	app := setupWebServerApp(db, pm, sugar)
	go func() {
		browser.OpenURL(conf.ADDR)
	}()
	sugar.Fatal(app.Listen(conf.ADDR))
}

func setupWebServerApp(db *gorm.DB, pm *pricing.PriceManager, sugar *zap.SugaredLogger) *fiber.App {
	s := server.Server{DB: db, Logger: sugar, PricingManager: pm}

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	//app.Use(cache.New(cache.Config{
	//	Next: func(c *fiber.Ctx) bool {
	//		return c.Query("refresh") == "true"
	//	},
	//	Expiration:   5 * time.Minute,
	//	CacheControl: true,
	//}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080, http://localhost:5000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-CSRF-Token, X-Requested-With, Accept-Language, Cache-Control, User-Agent",
		AllowMethods:     "GET, POST, OPTIONS , PUT, DELETE",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
	}))

	subFS, _ := fs.Sub(static, "dist")
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   http.FS(subFS),
		Browse: true,
	}))

	v1 := app.Group("/api").Group("/v1")
	v1.Get("/portfolios/", s.GetPortfolios)
	v1.Post("/portfolios::logValue", s.LogPortfoliosValue)
	v1.Get("/portfolios/:portfolio_id", s.GetPortfolio)

	return app
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

	if err := db.AutoMigrate(&pkg.Portfolio{}, &pkg.PortfolioItem{}, &pkg.PortfolioHistory{}); err != nil {
		return err
	}

	return db.Create(&p).Error
}
