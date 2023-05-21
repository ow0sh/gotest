package main

import "C"
import (
	"github.com/go-chi/chi/v5"
	"github.com/ow0sh/gotest/coingecko"
	"github.com/pkg/errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/patrickmn/go-cache"

	"github.com/ow0sh/gotest/config"
	middlint "github.com/ow0sh/gotest/middleware"
)

func main() {
	// TODO: return config, don't use global variables

	config.InitConfig()

	// TODO: add logger https://github.com/sirupsen/logrus
	httpCli := &http.Client{}
	coinCli := coingecko.NewClient(httpCli)
	bases, err := coinCli.GetCoins()
	if err != nil {
		panic(err)
	}
	// TODO: add error processing
	quotes := coinCli.GetSupported()

	c := cache.New(time.Duration(config.Config.Cache.DefaultExpiration)*time.Minute,
		time.Duration(config.Config.Cache.CleanupInterval)*time.Minute)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Route("/rate", func(r chi.Router) {
		r.Use(middlint.Caching(c))
		r.Get("/", newHandler(coinCli, bases, mapToSet(quotes)).convert)
	})

	if err := http.ListenAndServe(config.Config.App.Port, r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
