package main

import (
	"os"
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	TelegramApiKey string
	GoogleSafe string
	SkypeAppId string
	SkypeSecret string
	SkypeToken string
	ConfigGoogle
	SkypeConfig
}

type ConfigGoogle struct {
	GoogleSearchSafe string
}

type SkypeConfig struct {
	TokenType string 		`json:"token_type"`
	ExpiresIn int			`json:"expires_in"`
	ExtExpiresIn int		`json:"ext_expires_in"`
	AccessToken string		`json:"access_token"`
}


func (conf *Config) SetConfig(configfile string) error {
	_, err := os.Stat(configfile)
	if err != nil {
		fmt.Println("Config file is missing: ", configfile)
		os.Exit(0)
	}

	if _, err := toml.DecodeFile(configfile, &conf); err != nil {
		fmt.Printf("Failed to parse config file %s \n", configfile)
	}
	//log.Print(config.Index)
	c.ConfigGoogle.GoogleSearchSafe = "active"
	return nil
}