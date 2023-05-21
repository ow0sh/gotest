package main

import "C"
import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ow0sh/gotest/coingecko"
	"github.com/pkg/errors"
	"net/http"
)

type handler struct {
	coinCli *coingecko.Client
	bases   map[string]string
	quotes  map[string]struct{}
}

func newHandler(coinCli *coingecko.Client, bases map[string]string, quotes map[string]struct{}) handler {
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

	// TODO: make error handling
	price := h.coinCli.GetPrice(req.base, req.quote)
	w.Header().Set("Content-Type", "application/json")
	w.Write(price)
}

type request struct {
	base  string
	quote string
}

func (h handler) newRequest(r *http.Request) (*request, error) {
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
