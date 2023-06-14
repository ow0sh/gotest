package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ow0sh/gotest/coingecko"
	"github.com/ow0sh/gotest/postgres"

	"github.com/pkg/errors"

	"github.com/patrickmn/go-cache"

	"github.com/ow0sh/gotest/config"
	middlint "github.com/ow0sh/gotest/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	config := config.InitConfig(log)

	httpCli := &http.Client{}
	coinCli := coingecko.NewClient(httpCli)

	PSQLconn, _ := postgres.NewConn(config.PSQL)
	defer PSQLconn.CloseConn()

	bases, err := coinCli.GetCoins()
	if err != nil {
		panic(err)
	}
	quotes, err := coinCli.GetSupported()
	if err != nil {
		panic(err)
	}

	handler := NewHandler(coinCli, bases, MapToSet(quotes))

	c := cache.New(time.Duration(config.Cache.DefaultExpiration)*time.Minute,
		time.Duration(config.Cache.CleanupInterval)*time.Minute)
	r := chi.NewRouter()

	r.Use(middlint.Logger(log))
	r.Route("/rate/{base}-{quote}", func(r chi.Router) {
		r.Use(middlint.Caching(c))
		r.Get("/", handler.convert)
	})

	PSQLconn.InsertInfo(bases, MapToSet(quotes))
	log.Info("Server started on port: " + config.App.Port)
	if err := http.ListenAndServe(config.App.Port, r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
