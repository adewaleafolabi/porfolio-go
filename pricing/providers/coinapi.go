package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CoinAPI struct {
	APIKey  string
	BaseURL string
}

func (c *CoinAPI) GetPrice(code string) (float64, error) {
	url := fmt.Sprintf("%s/quotes/%s/current", c.BaseURL, code)

	client := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(url)
		return 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-CoinAPI-Key", c.APIKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(url)
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	data := map[string]interface{}{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err)
		fmt.Println(url)
		return 0, err
	}

	p, ok := data["ask_price"].(float64)
	if !ok {
		fmt.Println(err)
		fmt.Println(url)
		return 0, fmt.Errorf("coinapi ask_price for %s not valid", code)
	}

	return p, nil
}
