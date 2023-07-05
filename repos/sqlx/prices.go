package sqlx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ow0sh/gotest/repos"
)

type pricesRepo struct {
	baseRepo[repos.Price]
}

func NewPricesRepo(db *sqlx.DB) repos.PricesRepo {
	return &pricesRepo{
		newBaseRepo[repos.Price](db, "prices", "base", "quote", "rate"),
	}
}

func (s pricesRepo) Inserter() repos.PricesInserter {
	return s
}

func (s pricesRepo) Create(ctx context.Context) ([]repos.Price, error) {
	return s.baseRepo.Create(ctx)
}

func (s pricesRepo) SetCreate(prices ...repos.CreatePrice) repos.PricesInserter {
	for _, price := range prices {
		s.q.sqlInsert = s.q.sqlInsert.Values(price.Base, price.Quote, price.Rate)
	}

	return s
}

func (s pricesRepo) Selector() repos.PricesSelector {
	return s
}

func (s pricesRepo) FilterByPricesId(ids ...int64) repos.PricesSelector {
	s.q.sqlSelect = s.q.sqlSelect.Where(squirrel.Eq{"id": ids})
	return s
}

func (s pricesRepo) OrderBy(by string, order string) repos.PricesSelector {
	s.baseRepo = s.baseRepo.OrderBy(by, order)
	return s
}

func (s pricesRepo) Limit(u uint64) repos.PricesSelector {
	s.baseRepo = s.baseRepo.Limit(u)
	return s
}

func (s pricesRepo) Tx() repos.PricesTransactor {
	return s
}

func (s pricesRepo) Updater() repos.PricesUpdater {
	return s
}

func (s pricesRepo) WhereId(ids ...int64) repos.PricesUpdater {
	s.baseRepo = s.baseRepo.WhereID(ids...)
	return s
}

func (s pricesRepo) WhereBase(base string) repos.PricesUpdater {
	s.baseRepo = s.baseRepo.WhereBase(base)
	return s
}
