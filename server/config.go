package server

type Config struct {
	CoinAPIToken   string `required:"true"  envconfig:"coin_api_token"`
	CoinAPIBaseURL string `default:"https://rest.coinapi.io/v1/"`
	ADDR           string `default:"127.0.0.1:5000"`
	DailyCronTime  string `default:"17:00" envconfig:"daily_cron_time"`
}
