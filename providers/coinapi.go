package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CoinAPI struct {
}

func (c *CoinAPI) GetPrice(code string) (float64, error) {

	url := fmt.Sprintf("https://rest.coinapi.io/v1/quotes/%s/current", code)
	method := "GET"

	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-CoinAPI-Key", "7C321061-36E6-4953-95F8-E4FE65489450")

	res, err := client.Do(req)
	if err != nil {

		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	p, ok := data["ask_price"].(float64)
	if !ok {
		fmt.Println(string(body))
		return 0, fmt.Errorf("coinapi ask_price for %s not valid", code)
	}

	return p, nil
}
