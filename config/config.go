package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Coin struct {
	Id     string
	Symbol string
	Name   string
}

type TypeAppConfig struct {
	Port string `json:"PORT"`
}

type TypeCacheConfig struct {
	DefaultExpiration int `json:"defaultExpiration"`
	CleanupInterval   int `json:"cleanupInterval"`
}

type TypeConfig struct {
	App   TypeAppConfig   `json:"app"`
	Cache TypeCacheConfig `json:"cache"`
}

var Config TypeConfig

func InitConfig() {
	configFilename := "default.json"

	configFile, err := ioutil.ReadFile("./config/" + configFilename)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(configFile, &Config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Config initialized")
}
