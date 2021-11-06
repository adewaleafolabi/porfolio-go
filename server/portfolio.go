package server

import (
	"net/http"
	"portfolio/cron"
	"portfolio/pkg"
	"portfolio/pricing"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Server struct {
	DB             *gorm.DB
	Logger         *zap.SugaredLogger
	PricingManager *pricing.PriceManager
}

func (h *Server) GetPortfolio(c *fiber.Ctx) error {
	id := c.Params("id")

	p := pkg.Portfolio{}

	if err := h.DB.Preload(clause.Associations).First(&p, "id=?", id).Error; err != nil {
		return c.Status(500).SendString("No Portfolio Found with ID")
	}

	return c.Status(200).JSON(h.enrichPortfolio(p))
}

func (h *Server) GetPortfolios(c *fiber.Ctx) error {
	var portfolios []pkg.Portfolio

	if err := h.DB.Preload(clause.Associations).Find(&portfolios).Error; err != nil {
		h.Logger.Error(err)
		return c.Status(500).SendString("No Portfolios Found")
	}
	for i := range portfolios {
		portfolios[i] = h.enrichPortfolio(portfolios[i])
	}
	return c.Status(200).JSON(portfolios)
}

func (h *Server) LogPortfoliosValue(c *fiber.Ctx) error {
	go func() {
		cron.LogValue(h.DB, h.PricingManager, h.Logger)
	}()
	return c.SendStatus(http.StatusOK)
}

func (h *Server) enrichPortfolio(p pkg.Portfolio)pkg.Portfolio  {
	var items []pricing.GetPricingRequestItem
	for _, item := range p.Items {
		if item.Symbol == p.BaseCurrency {
			continue
		}
		items = append(items, pricing.GetPricingRequestItem{
			Symbol:   item.Symbol,
			Provider: item.PriceProvider,
		})
	}
	request := pricing.GetPricingRequest{
		BaseCurrency: p.BaseCurrency,
		Items:        items,
	}

	prices, err := h.PricingManager.GetPricing(request)

	if err != nil {
		//do not error out if pricing fails
		h.Logger.Error(err)
	}

	p.UpdateTotalValue(prices)
	if p.TotalValue > p.AllTimeHigh {
		p.AllTimeHigh = p.TotalValue
		if err := h.DB.Model(&pkg.Portfolio{}).Where("id = ?", p.ID).Update("all_time_high", p.AllTimeHigh).Error; err != nil {
			h.Logger.Errorw("error storing all time high to portfolio", "err", err)
		}
	}
	return p
}