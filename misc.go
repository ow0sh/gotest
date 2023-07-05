package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ow0sh/gotest/coingecko"
	"github.com/ow0sh/gotest/models"
	"github.com/ow0sh/gotest/usecases"
	"github.com/sirupsen/logrus"
)

func MapToSet(inter interface{}) map[string]struct{} {
	var result = make(map[string]struct{})
	switch val := inter.(type) {
	case []string:
		for _, str := range val {
			result[str] = struct{}{}
		}
	case map[string]string:
		for str := range val {
			result[str] = struct{}{}
		}
	default:
		panic("such type is not impl")
	}

	return result
}

func UpdateDB(log *logrus.Logger, pricesUse *usecases.PricesUseCase, cli *coingecko.Client, bq models.BQ, ctx context.Context) {
	for {
		for i := range bq.Base {
			jsonRate, err := cli.GetPrice(bq.Base[i], bq.Quote[i])
			if err != nil {
				log.Error(err, "failed to get rate")
			}
			var data map[string]map[string]float64
			if err := json.Unmarshal(jsonRate, &data); err != nil {
			}
			var rate float64
			for _, quote := range data {
				for _, r := range quote {
					rate = float64(r)
				}
			}

			test, err := pricesUse.UpdatePrices(ctx, bq.Base[i], rate)
			if err != nil {
				log.Error(err)
				return
			}
			if len(test) == 0 {
				_, err = pricesUse.CreatePrices(ctx, models.Price{Base: bq.Base[i], Quote: bq.Quote[i], Rate: rate})
				if err != nil {
					log.Error(err)
					return
				} else {
					log.Info("Pairs created and updated successfully")
				}
			} else {
				log.Info("Pairs updated successfully")
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
