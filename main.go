package main

import (
	"fmt"
	"time"
	"trade/bitflyer"
	"trade/config"
	"trade/utils"
)

func init() {
	utils.LoggingSettigns(config.Config.LogFile)
}

func main() {
	apiClient := bitflyer.New(config.Config.APIKey, config.Config.APISecret)
	ticker, err := apiClient.Ticker("BTC_JPY")
	if err != nil {
		panic(err)
	}
	fmt.Println(ticker)
	fmt.Println(ticker.MidPrice())
	fmt.Println(ticker.DateTime())
	fmt.Println(ticker.TruncateDateTime(time.Hour))
}
