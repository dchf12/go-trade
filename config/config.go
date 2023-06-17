package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	APIKey    string
	APISecret string
	LogFile   string
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

	Config = ConfigList{
		APIKey:    cfg.Section("bitflyer").Key("api_key").String(),
		APISecret: cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:   cfg.Section("trade").Key("log_file").String(),
	}
}
