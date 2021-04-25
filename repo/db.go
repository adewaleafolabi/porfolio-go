package repo

import "portfolio/pkg"

type PortfolioStore interface {
	GetPortfolio(id string) (*pkg.Portfolio, error)
	GetPortfolios() ([]*pkg.Portfolio, error)

	CreatePortfolio(p pkg.Portfolio) (pkg.Portfolio, error)
	UpdatePortfolio(p pkg.Portfolio) (pkg.Portfolio, error)

	DeletePortfolio(id string) (pkg.Portfolio, error)

	UpdatePortfolioHistory(ph pkg.PortfolioHistory) (pkg.PortfolioHistory, error)
}
