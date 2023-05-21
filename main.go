package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patrickmn/go-cache"

	"server/config"
	C "server/controller"
	M "server/middleware"
)

func main() {
	config.InitConfig()

	client := &http.Client{}
	coinlist := C.GetCoins(client)               // map[symbol]name
	supported_coinlist := C.GetSupported(client) // []string

	c := cache.New(time.Duration(config.Config.Cache.DefaultExpiration)*time.Minute,
		time.Duration(config.Config.Cache.CleanupInterval)*time.Minute)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Route("/rate/{base}-{quote}", func(r chi.Router) {
		r.Use(M.Caching(c))
		r.Get("/", convert(c, coinlist, supported_coinlist))
	})

	http.ListenAndServe(config.Config.App.Port, r)
}

func convert(c *cache.Cache, coinlist map[string]string, supported_coinlist []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		base := chi.URLParam(r, "base")
		quote := chi.URLParam(r, "quote")

		if C.Contains(supported_coinlist, quote) == false {
			http.Error(w, fmt.Sprint("Unsupported second cryptocurrency, supported: ", supported_coinlist), http.StatusBadRequest)
			return
		}
		if C.Contains(coinlist, base) == false {
			http.Error(w, "Unknown cryptocurrency", http.StatusBadRequest)
			return
		}

		price := C.GetPrice(coinlist[base], quote)
		c.Set(fmt.Sprintf(base+"-"+quote), []byte(price), cache.DefaultExpiration)
		w.Header().Set("Content-Type", "application/json")
		w.Write(price)
	}
}
