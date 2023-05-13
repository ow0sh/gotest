package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Coin struct {
	Id     string
	Symbol string
	Name   string
}

type Supported_Coin struct {
	Symbol string
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/convert/{base}-{quote}", convert)

	http.ListenAndServe(":3000", r)
}

func convert(w http.ResponseWriter, r *http.Request) {
	coinlist := getCoins()
	supported_coinlist := getSupported()

	base := chi.URLParam(r, "base")
	quote := chi.URLParam(r, "quote")

	for i, v := range coinlist {
		if v.Symbol == base {
			base = strings.ToLower(v.Name)
			break
		}

		if i == len(coinlist) {
			w.Write([]byte(fmt.Sprint("I don't know about that cryptocurrency")))
		}
	}

	for i, v := range supported_coinlist {
		if v == quote {
			break
		}

		if i == len(supported_coinlist) {
			w.Write([]byte(fmt.Sprint("Please choose only supported coins as second parameter: ", supported_coinlist)))
			return
		}
	}

	price := getPrice(base, quote)
	w.Write([]byte(fmt.Sprint("One " + base + " is " + fmt.Sprint(price) + " " + quote)))
}

func getCoins() []Coin {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/list", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var coinlist []Coin
	json.NewDecoder(resp.Body).Decode(&coinlist)

	return coinlist
}

func getSupported() []string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/supported_vs_currencies", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var supported_coinlist []string
	json.NewDecoder(resp.Body).Decode(&supported_coinlist)

	return supported_coinlist
}

func getPrice(base, quote string) float64 {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%v&vs_currencies=%v", base, quote), nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var price map[string]map[string]float64
	json.NewDecoder(resp.Body).Decode(&price)
	fmt.Println(price)

	return price[base][quote]
}
