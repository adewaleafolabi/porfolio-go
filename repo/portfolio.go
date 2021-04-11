package repo

import (
	"github.com/segmentio/ksuid"
	"portfolio/pkg"
)

func CreatePortfolio(p *pkg.Portfolio) (pkg.Portfolio, error) {
	id := ksuid.New().String()
	p.ID = id
	return *p, nil
}

func GetPortfolio(id string) (*pkg.Portfolio, error) {
	return nil, nil
}

func ListPortfolio(id string) ([]*pkg.Portfolio, error) {
	return nil, nil
}

func UpdatePortfolio(p *pkg.Portfolio) (*pkg.Portfolio, error) {
	return p, nil
}

func UpsertPortfolioItem(p *pkg.PortfolioItem) (*pkg.PortfolioItem, error) {
	return p, nil
}
