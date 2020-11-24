package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	TelegramApiKey string
	GoogleSafe     string
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
	return nil
}
