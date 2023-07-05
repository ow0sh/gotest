package main

import (
	"encoding/json"
	"time"

	"github.com/ow0sh/gotest/coingecko"
	"github.com/ow0sh/gotest/postgres"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type BQ struct {
	base  []string
	quote []string
}

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

func UpdateDB(log *logrus.Logger, conn *postgres.PSQLConn, cli *coingecko.Client, bq BQ) error {
	for {
		for i := range bq.base {
			jsonRate, err := cli.GetPrice(bq.base[i], bq.quote[i])
			if err != nil {
				log.Error(err, "failed to get rate")
				return errors.Wrap(err, "failed to get rate")
			}
			var data map[string]map[string]float64
			if err := json.Unmarshal(jsonRate, &data); err != nil {
				log.Error(err, "failed to unmarshal jsonprice")
				return errors.Wrap(err, "failed to unmarshal jsonprice")
			}
			var rate float64
			for _, quote := range data {
				for _, r := range quote {
					rate = float64(r)
				}
			}

			if conn.Exist(bq.base[i]) {
				err = conn.UpdateInfo(log, bq.base[i], rate)
				if err != nil {
					log.Error("failed to update info")
					return errors.Wrap(err, "failed to update info")
				}
			} else {
				err = conn.InsertInfo(log, bq.base[i], bq.quote[i], rate)
				if err != nil {
					log.Error("failed to updateDB")
					return errors.Wrap(err, "failed to insert info")
				}
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
