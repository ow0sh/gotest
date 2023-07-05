package repos

import "context"

type PricesTransactor interface {
	QDelete(ctx context.Context, q DB) error
	QCreate(ctx context.Context, q DB) ([]Price, error)
	QUpdate(ctx context.Context, q DB) ([]Price, error)
	QGet(ctx context.Context, q DB) (*Price, error)
}

type PricesInserter interface {
	Create(ctx context.Context) ([]Price, error)
	SetCreate(...CreatePrice) PricesInserter
}

type PricesSelector interface {
	FilterByPricesId(...int64) PricesSelector
	OrderBy(string, string) PricesSelector
	Limit(uint64) PricesSelector

	Select(ctx context.Context) ([]Price, error)
	Get(ctx context.Context) (*Price, error)
}

type PricesUpdater interface {
	WhereId(...int64) PricesUpdater
	WhereBase(base string) PricesUpdater
	Update(ctx context.Context, column string, value interface{}) ([]Price, error)
}

type PricesRepo interface {
	Updater() PricesUpdater
	Selector() PricesSelector
	Inserter() PricesInserter
	Tx() PricesTransactor
}

type CreatePrice struct {
	Base  string
	Quote string
	Rate  float64
}

type Price struct {
	Id    int64   `db:"id"`
	Base  string  `db:"base"`
	Quote string  `db:"quote"`
	Rate  float64 `db:"rate"`
}
