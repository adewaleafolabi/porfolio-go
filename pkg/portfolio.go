package pkg

import (
	"fmt"
	"portfolio/pricing"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"gorm.io/gorm"
)

const (
	Cash   AssetType = "CASH"
	Stock  AssetType = "STOCK"
	Crypto AssetType = "CRYPTO"
	Debt   AssetType = "DEBT"
)

type AssetType string

type Portfolio struct {
	gorm.Model
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	BaseCurrency  string             `json:"base_currency"`
	TotalValue    float64            `json:"total_value"`
	Goal          float64            `json:"goal"`
	TotalValueUSD float64            `json:"total_value_usd" gorm:"-"`
	Items         []PortfolioItem    `json:"items"`
	History       []PortfolioHistory `json:"history"`
}

type PortfolioItem struct {
	gorm.Model
	PortfolioID   string                `json:"portfolio_id"`
	Symbol        string                `json:"symbol"`
	AssetType     string                `json:"asset_type"`
	Icon          string                `json:"icon"`
	Quantity      float64               `json:"quantity"`
	TotalValue    float64               `json:"total_value"`
	UnitPrice    float64               `json:"unit_price" gorm:"-"`
	Label         string                `json:"label"`
	PriceProvider pricing.PriceProvider `json:"price_provider"`
}

func (pi *PortfolioItem) DisplayLabel() string {
	if pi.Label != "" {
		return pi.Label
	}
	return pi.Symbol
}

func (pi *PortfolioItem) String() string {
	return fmt.Sprintf("%s\t%.2f\t%.2f\n", pi.DisplayLabel(), pi.Quantity, pi.TotalValue)
}

type PortfolioHistory struct {
	gorm.Model
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Value       float64   `json:"value"`
	SnapShot    string    `json:"-"`
	PortfolioID string    `json:"portfolio_id"`
}

func (p *Portfolio) UpdateTotalValue(prices map[string]float64) {
	p.TotalValue = 0
	for i := range p.Items {
		if p.Items[i].Symbol == p.BaseCurrency {
			p.Items[i].TotalValue = p.Items[i].Quantity
			p.TotalValue += p.Items[i].TotalValue
			continue
		}
		if price, ok := prices[p.Items[i].Symbol]; ok {
			p.Items[i].TotalValue = price * p.Items[i].Quantity
			p.Items[i].UnitPrice = price
			p.TotalValue += p.Items[i].TotalValue
		}
	}

	p.TotalValueUSD = p.TotalValue / prices[fmt.Sprintf("%s=X", p.BaseCurrency)]

	sort.Slice(p.Items, func(i, j int) bool {
		return p.Items[i].TotalValue > p.Items[j].TotalValue
	})
}

func (p Portfolio) GeneratePricingRequest() pricing.GetPricingRequest {
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
	return pricing.GetPricingRequest{
		BaseCurrency: p.BaseCurrency,
		Items:        items,
	}

}

func (p *Portfolio) String() string {
	buffer := strings.Builder{}
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{"Item", "Quantity", "Value", "Weight (%)"})
	table.SetFooter([]string{"", "", "Total", fmt.Sprintf("$%s", humanize.CommafWithDigits(p.TotalValue, 2))})
	table.SetBorder(false)
	for _, item := range p.Items {
		weight := 0.0
		if item.TotalValue != 0 && p.TotalValue != 0 {
			weight = (item.TotalValue / p.TotalValue) * 100
		}
		table.Append([]string{item.DisplayLabel(), fmt.Sprintf("%v", item.Quantity), fmt.Sprintf("%.2f", item.TotalValue), fmt.Sprintf("%.2f", weight)})
	}
	table.Render()
	return fmt.Sprintf("Portfolio Name: %s\nTotal Value: %s %s\n%s", p.Name, p.BaseCurrency, fmt.Sprintf("%.2f", p.TotalValue), buffer.String())
}
