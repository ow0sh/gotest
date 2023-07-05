package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ow0sh/gotest/coingecko"
	"github.com/pkg/errors"
)

type handler struct {
	coinCli *coingecko.Client
	bases   map[string]string
	quotes  map[string]struct{}
}

func NewHandler(coinCli *coingecko.Client, bases map[string]string, quotes map[string]struct{}) handler {
	return handler{
		coinCli: coinCli,
		bases:   bases,
		quotes:  quotes,
	}
}

func (h *handler) convert(w http.ResponseWriter, r *http.Request) {
	req, err := h.newRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	price, err := h.coinCli.GetPrice(req.base, req.quote)
	if err != nil {
		errors.Wrap(err, "failed to get price")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(price)
}

// func (h *handler) getCoinList(coinlist []coingecko.Coin) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tmp, err := json.Marshal(coinlist)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(tmp)
// 	}
// }

// func (h *handler) getSupportedList(supportedList []string) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tmp, err := json.Marshal(supportedList)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(tmp)
// 	}
// }

type request struct {
	base  string
	quote string
}

func (h *handler) newRequest(r *http.Request) (*request, error) {
	base := chi.URLParam(r, "base")
	quote := chi.URLParam(r, "quote")

	if _, ok := h.quotes[quote]; !ok {
		return nil, errors.New(fmt.Sprintf("such quote coin %s is not supported", quote))
	}

	baseLong, ok := h.bases[base]
	if !ok {
		return nil, errors.New(fmt.Sprintf("such base coin %s is not supported", base))
	}

	return &request{
		base:  baseLong,
		quote: quote,
	}, nil
}
