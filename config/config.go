package config

import (
	"encoding/json"
	"os"

	"github.com/jmoiron/sqlx"
	gocache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Config interface {
	DB() *sqlx.DB
	Log() *logrus.Logger
	C() *gocache.Cache
}

type config struct {
	db
	logger
	c
}

func NewConfig(cfgPath string) (Config, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	cfg := config{}
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
