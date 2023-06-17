package main

import (
	"fmt"
	"trade/bitflyer"
	"trade/config"
	"trade/utils"
)

func init() {
	utils.LoggingSettigns(config.Config.LogFile)
}

func main() {
	apiClient := bitflyer.New(config.Config.APIKey, config.Config.APISecret)

	order := &bitflyer.Order{
		ProductCode:     config.Config.ProductCode,
		ChildOrderType:  "MARKET",
		Side:            "SELL",
		Size:            0.0001,
		MinuteToExpires: 1,
		TimeInForce:     "GTC",
	}
	res, _ := apiClient.SendOrder(order)
	fmt.Println(res.ChildOrderAcceptanceID)
}
