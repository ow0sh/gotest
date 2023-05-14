package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patrickmn/go-cache"
)

type Coin struct {
	Id     string
	Symbol string
	Name   string
}

var coinlist []Coin
var supported_coinlist []string
var c *cache.Cache

func main() {
	c = cache.New(5*time.Minute, 10*time.Minute)
	c.Set("coinlist", getCoins(), cache.DefaultExpiration)
	c.Set("supported_coinlist", getSupported(), cache.DefaultExpiration)

	if coinlistI, ok := c.Get("coinlist"); ok {
		coinlist = coinlistI.([]Coin)
	}

	if supported_coinlistI, ok := c.Get("supported_coinlist"); ok {
		supported_coinlist = supported_coinlistI.([]string)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/convert/{base}-{quote}", convert)

	http.ListenAndServe(":3000", r)
	fmt.Println("Listening on port 3000")
}

func convert(w http.ResponseWriter, r *http.Request) {
	if _, ok := c.Get("coinlist"); !ok {
		c.Set("coinlist", getCoins(), cache.DefaultExpiration)
	}

	if _, ok := c.Get("supported_coinlist"); !ok {
		c.Set("supported_coinlist", getSupported(), cache.DefaultExpiration)
	}

	base := chi.URLParam(r, "base")
	quote := chi.URLParam(r, "quote")

	for i, v := range coinlist {
		if v.Symbol == base {
			base = strings.ToLower(v.Name)
			break
		}

		if i == len(coinlist)-1 {
			http.Error(w, "Unknown cryptocurrency", http.StatusBadRequest)
			return
		}
	}

	for i, v := range supported_coinlist {
		if v == quote {
			break
		}

		if i == len(supported_coinlist)-1 {
			http.Error(w, fmt.Sprint("Unsupported second cryptocurrency, supported: ", supported_coinlist), http.StatusBadRequest)
			return
		}
	}

	price := getPrice(base, quote)
	w.Header().Set("Content-Type", "application/json")
	w.Write(price)
}

func getCoins() []Coin {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/list", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var tmpCoinlist []Coin
	json.NewDecoder(resp.Body).Decode(&tmpCoinlist)

	return tmpCoinlist
}

func getSupported() []string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/supported_vs_currencies", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var tmpSupportedCoinlist []string
	json.NewDecoder(resp.Body).Decode(&tmpSupportedCoinlist)

	return tmpSupportedCoinlist
}

func getPrice(base, quote string) json.RawMessage {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%v&vs_currencies=%v", base, quote), nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var jsonPrice json.RawMessage
	json.NewDecoder(resp.Body).Decode(&jsonPrice)

	return jsonPrice
}
