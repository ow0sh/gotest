package repos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ow0sh/gotest/models"
)

var (
	ErrSuchObjectAlreadyExist = errors.New("such object is already exist")
	ErrNothingUpdated         = errors.New("nothing updated")
)

type DB interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Transactor interface {
	QDelete(ctx context.Context, q DB) error
	QCreate(ctx context.Context, q DB) ([]int64, error)
	QUpdate(ctx context.Context, q DB) error
}

type OrderBy string

func PricesToCreateRepo(prices ...models.Price) []CreatePrice {
	result := make([]CreatePrice, len(prices))

	for i, price := range prices {
		result[i] = CreatePrice{
			Base:  price.Base,
			Quote: price.Quote,
			Rate:  price.Rate,
		}
	}

	return result
}
