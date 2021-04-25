package server

import (
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

	if err := h.DB.Preload(clause.Associations).Find(&p, id).Error; err != nil {
		return c.Status(500).SendString("No Portfolio Found with ID")
	}

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

	return c.Status(200).JSON(p)
}

func (h *Server) GetPortfolios(c *fiber.Ctx) error {
	var portfolios []pkg.Portfolio

	if err := h.DB.Find(&portfolios).Error; err != nil {
		h.Logger.Error(err)
		return c.Status(500).SendString("No Portfolios Found")
	}

	return c.Status(200).JSON(portfolios)
}
