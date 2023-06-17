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

	tickerChannel := make(chan bitflyer.Ticker)
	go apiClient.ReacTimeTicker(config.Config.ProductCode, tickerChannel)
	for ticker := range tickerChannel {
		fmt.Println(ticker)
		fmt.Println(ticker.MidPrice())
		fmt.Println(ticker.DateTime())
		fmt.Println(ticker.TruncateDateTime(time.Hour))
	}

	// order := &bitflyer.Order{
	// 	ProductCode:     config.Config.ProductCode,
	// 	ChildOrderType:  "MARKET",
	// 	Side:            "SELL",
	// 	Size:            0.0001,
	// 	MinuteToExpires: 1,
	// 	TimeInForce:     "GTC",
	// }
	// res, _ := apiClient.SendOrder(order)
	// fmt.Println(res.ChildOrderAcceptanceID)
}
