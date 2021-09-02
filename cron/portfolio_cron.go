package cron

import (
	"context"
	"fmt"
	"portfolio/pkg"
	"portfolio/pricing"
	"time"

	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func LogValue(db *gorm.DB, pm *pricing.PriceManager, logger *zap.SugaredLogger) {
	var portfolios []pkg.Portfolio

	if err := db.Preload(clause.Associations).Find(&portfolios).Error; err != nil {
		logger.Error("error fetching portfolios from DB", err)
		return
	}
	for i := 0; i < len(portfolios); i++ {
		req := portfolios[i].GeneratePricingRequest()
		prices, err := pm.GetPricing(req)
		if err != nil {
			logger.Error(err)
			priceRecovery(context.Background(), pm, &portfolios[i], logger)
		}
		portfolios[i].UpdateTotalValue(prices)

		hist := pkg.PortfolioHistory{
			ID:          ksuid.New().String(),
			Date:        time.Now(),
			Value:       portfolios[i].TotalValue,
			SnapShot:    portfolios[i].String(),
			PortfolioID: portfolios[i].ID,
		}

		if err := db.Create(&hist).Error; err != nil {
			logger.Errorw("error storing portfolio history", "err", err)
		}
		if portfolios[i].TotalValue > portfolios[i].AllTimeHigh {
			portfolios[i].AllTimeHigh = portfolios[i].TotalValue
			if err := db.Model(&pkg.Portfolio{}).Where("id = ?", portfolios[i].ID).Update("all_time_high", portfolios[i].AllTimeHigh).Error; err != nil {
				logger.Errorw("error storing all time high to portfolio", "err", err)
			}

		}
		fmt.Printf("Run Date: %s\n", time.Now().String())
		fmt.Println(portfolios[i].String())
	}
}

func priceRecovery(ctx context.Context, pm *pricing.PriceManager, p *pkg.Portfolio, logger *zap.SugaredLogger) {
	for {
		select {
		case <-time.After(5 * time.Minute):
			logger.Infow("running price recovery", "time", time.Now().Format(time.Kitchen))
			prices, err := pm.GetPricing(p.GeneratePricingRequest())
			if err == nil {
				p.UpdateTotalValue(prices)
				logger.Info("price recovery complete")
				break
			}
		case <-ctx.Done():
			logger.Info("ctx halted")
		}
	}
}
