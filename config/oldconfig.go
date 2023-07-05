package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type TypeAppConfig struct {
	Port string `json:"PORT"`
}

type TypeCacheConfig struct {
	DefaultExpiration int `json:"defaultExpiration"`
	CleanupInterval   int `json:"cleanupInterval"`
}

type TypePSQLConfig struct {
	Dsn      string `json:"dsn"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Dbname   string `json:"dbname"`
	Driver   string `json:"driver"`
}

type TypeConfig struct {
	App   TypeAppConfig   `json:"app"`
	Cache TypeCacheConfig `json:"cache"`
	PSQL  TypePSQLConfig  `json:"psql"`
}

func InitConfig(log *logrus.Logger) TypeConfig {
	configFilename := "default.json"

	configFile, err := ioutil.ReadFile("./config/" + configFilename)
	if err != nil {
		errors.Wrap(err, "failed to read config file")
	}

	var config TypeConfig
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		errors.Wrap(err, "failed to unmarshal json")
	}

	log.Info("Config initialized")
	return config
}
