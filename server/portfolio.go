package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"portfolio/pkg"
)

type Server struct {
	DB *gorm.DB
}

 func (h *Server)  GetPortfolio(c *fiber.Ctx) error {
	id := c.Params("id")

	p := pkg.Portfolio{}

	if err := h.DB.Preload(clause.Associations).Find(&p,id).Error; err != nil {

		return c.Status(500).SendString("No Portfolio Found with ID")
	}

	p.UpdateTotalValue()

	return c.Status(200).JSON(p)
}

 func (h *Server)  GetPortfolios(c *fiber.Ctx) error {
	 var portfolios []pkg.Portfolio

	if err := h.DB.Find(&portfolios).Error; err != nil {

		return c.Status(500).SendString("No Portfolios Found")
	}

	return c.Status(200).JSON(portfolios)
}
