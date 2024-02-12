package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	APIKey      string
	APISecret   string
	LogFile     string
	ProductCode string

	TradeDuration time.Duration
	Durations     map[string]time.Duration
	DBName        string
	SQLDriver     string
	Port          int
}

var Config ConfigList

func (c ConfigList) String() string {
	return fmt.Sprintf("APIKey: XXXXX, APISecret: XXXXX")
}

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	durations := map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}

	Config = ConfigList{
		APIKey:        cfg.Section("bitflyer").Key("api_key").String(),
		APISecret:     cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:       cfg.Section("trade").Key("log_file").String(),
		ProductCode:   cfg.Section("trade").Key("product_code").String(),
		TradeDuration: durations[cfg.Section("trade").Key("trade_duration").String()],
		DBName:        cfg.Section("db").Key("name").String(),
		SQLDriver:     cfg.Section("db").Key("driver").String(),
		Port:          cfg.Section("web").Key("port").MustInt(),
	}
}
