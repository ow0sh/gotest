package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/ow0sh/gotest/config"
	middlint "github.com/ow0sh/gotest/middleware"
	"github.com/ow0sh/gotest/models"
	"github.com/ow0sh/gotest/pkg/coingecko"
	sqlx2 "github.com/ow0sh/gotest/repos/sqlx"
	"github.com/ow0sh/gotest/usecases"
	"github.com/pkg/errors"
)

const defaultConfigPath = "./config.json"

func main() {
	cfg, err := config.NewConfig(defaultConfigPath)
	if err != nil {
		panic(err)
	}

	log := cfg.Log()
	db := cfg.DB()
	c := cfg.C()

	ctx, cancel := ctxWithSig()
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			cancel()
		}
	}()

	coinCli := coingecko.NewClient()

	bases, err := coinCli.GetCoinList(ctx)
	if err != nil {
		panic(err)
	}
	quotes, err := coinCli.GetSupported(ctx)
	if err != nil {
		panic(err)
	}

	handler := NewHandler(ctx, coinCli, bases, MapToSet(quotes))

	r := chi.NewRouter()
	r.Use(middlint.Logger(log))
	r.Route("/rate/{base}-{quote}", func(r chi.Router) {
		r.Use(middlint.Caching(c))
		r.Get("/", handler.convert)
	})

	pricesUse := usecases.NewPricesUseCase(sqlx2.NewPricesRepo(db))

	params := models.BQ{Base: []string{"bitcoin", "ethereum", "solana", "binancecoin"},
		Quote: []string{"usd", "usd", "usd", "usd"}}
	go UpdateDB(log, pricesUse, coinCli, params, ctx)

	log.Info("Server started on port: 3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func ctxWithSig() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	go func() {
		select {
		case <-ch:
			cancel()
		}
	}()

	return ctx, cancel
}
