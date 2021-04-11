package pkg

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/piquette/finance-go/quote"
	"gorm.io/gorm"
	"portfolio/providers"
	"sort"
	"strings"
	"time"
)

type PriceProvider string

const (
	CoinAPI PriceProvider = "CoinAPI"
	Yahoo   PriceProvider = "Yahoo"
)

type Portfolio struct {
	gorm.Model
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	BaseCurrency string             `json:"base_currency"`
	TotalValue   float64            `json:"total_value"`
	Items        []PortfolioItem    `json:"items"`
	History      []PortfolioHistory `json:"history"`
}

type PortfolioItem struct {
	gorm.Model
	PortfolioID   string        `json:"portfolio_id"`
	Symbol        string        `json:"symbol"`
	Icon          string        `json:"icon"`
	Quantity      float64       `json:"quantity"`
	TotalValue    float64       `json:"total_value"`
	Label         string        `json:"label"`
	PriceProvider PriceProvider `json:"price_provider"`
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

func (p *Portfolio) UpdateTotalValue() error {
	usdRate, err := p.GetBaseCurrencyUSDRate()
	if err != nil {
		return err
	}

	var symbols []string
	var symbolsCoinApi []string
	for _, item := range p.Items {
		if item.Symbol == p.BaseCurrency {
			continue
		}
		if item.PriceProvider == CoinAPI {
			symbolsCoinApi = append(symbolsCoinApi, item.Symbol)
		}
		symbols = append(symbols, item.Symbol)
	}

	coinProvider := providers.CoinAPI{}
	registry := map[string]*float64{}
	for _, s := range symbolsCoinApi {
		p, err := coinProvider.GetPrice(s)
		if err != nil {
			return err
		}
		val := p * usdRate
		registry[s] = &val
	}
	quotes := quote.List(symbols)

	for quotes.Next() {
		q := quotes.Quote()
		registry[q.Symbol] = &q.RegularMarketPrice
		if q.CurrencyID == "USD" {
			price := q.RegularMarketPrice * usdRate
			registry[q.Symbol] = &price
		}
	}
	p.TotalValue = 0
	for i := range p.Items {
		if p.Items[i].Symbol == p.BaseCurrency {
			p.Items[i].TotalValue = p.Items[i].Quantity
			p.TotalValue += p.Items[i].TotalValue
			continue
		}
		price := registry[p.Items[i].Symbol]
		if price != nil {
			p.Items[i].TotalValue = *price * p.Items[i].Quantity
			p.TotalValue += p.Items[i].TotalValue
		}
	}
	sort.Slice(p.Items, func(i, j int) bool {
		return p.Items[i].TotalValue > p.Items[j].TotalValue
	})

	return nil
}

func (p *Portfolio) GetBaseCurrencyUSDRate() (float64, error) {
	q, err := quote.Get(fmt.Sprintf("%s=X", p.BaseCurrency))
	if err != nil {
		return 0.0, err
	}
	return q.RegularMarketPrice, nil
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
