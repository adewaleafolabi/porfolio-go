package pricing

import (
	"fmt"
	"portfolio/pricing/providers"

	"github.com/piquette/finance-go/quote"
)

type PriceProvider string

const (
	PriceProviderCoinAPI PriceProvider = "CoinAPI"
	PriceProviderNone    PriceProvider = ""
)

type PriceManager struct {
	priceCache *PriceCache
	coinAPI    *providers.CoinAPI
}

func NewPriceManager(coinAPIToken, coinAPIBaseURL string) *PriceManager {
	return &PriceManager{
		priceCache: NewPriceCache(300),
		coinAPI: &providers.CoinAPI{
			APIKey:  coinAPIToken,
			BaseURL: coinAPIBaseURL,
		},
	}
}

type GetPricingRequest struct {
	BaseCurrency string
	Items        []GetPricingRequestItem
}

type GetPricingRequestItem struct {
	Symbol   string
	Provider PriceProvider
}

func getBaseCurrencyUSDRate(baseCurrency string) (float64, error) {
	q, err := quote.Get(fmt.Sprintf("%s=X", baseCurrency))
	if err != nil {
		return 0.0, err
	}
	return q.RegularMarketPrice, nil
}

func (p *PriceManager) GetPricing(request GetPricingRequest) (map[string]float64, error) {
	output := make(map[string]float64)

	usdBaseCurrencyRate, err := getBaseCurrencyUSDRate(request.BaseCurrency)
	if err != nil {
		return output, err
	}
	output[fmt.Sprintf("%s=X", request.BaseCurrency)] = usdBaseCurrencyRate

	symbols := make(map[PriceProvider][]string)

	for _, item := range request.Items {
		symbols[item.Provider] = append(symbols[item.Provider], item.Symbol)
	}

	if items, ok := symbols[PriceProviderNone]; ok {
		quotes := quote.List(items)
		for quotes.Next() {
			q := quotes.Quote()
			output[q.Symbol] = q.RegularMarketPrice
			if q.CurrencyID == "USD" {
				price := q.RegularMarketPrice * usdBaseCurrencyRate
				output[q.Symbol] = price
			}
		}
	}

	for _, s := range symbols[PriceProviderCoinAPI] {
		price := p.priceCache.Get(s)
		if price == 0 {
			price, err = p.coinAPI.GetPrice(s)
			if err != nil {
				return output, err
			}
			p.priceCache.Set(s, price)
		}

		val := price * usdBaseCurrencyRate
		output[s] = val
	}
	return output, nil
}
