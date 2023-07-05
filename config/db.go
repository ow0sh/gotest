package config

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type db struct {
	dbParams `json:"db"`
	db       *sqlx.DB
}

type dbParams struct {
	URL    string `json:"url"`
	Driver string `json:"driver"`
}

func (d *db) validate() error {
	return errors.Wrap(d.check(), "failed to validate db")
}

func (d *db) check() error {
	if _, err := pq.ParseURL(d.URL); err != nil {
		return err
	}

	return nil
}

func (d *db) DB() *sqlx.DB {
	if d.db == nil {
		db, err := sqlx.Open(d.Driver, d.URL)
		if err != nil {
			panic(err)
		}

		if err = db.Ping(); err != nil {
			panic(err)
		}

		d.db = db
	}

	return d.db
}
