package sqlx

import (
	"context"
	sqlerr "database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ow0sh/gotest/repos"
	"github.com/pkg/errors"
)

type baseRepo[T any] struct {
	db *sqlx.DB
	q  querier
}

func newBaseRepo[T any](db *sqlx.DB, tableName string, columnsName ...string) baseRepo[T] {
	return baseRepo[T]{
		db: db,
		q:  newQuerier(tableName, columnsName...),
	}
}

func (s baseRepo[T]) Update(ctx context.Context, column string, value interface{}) ([]T, error) {
	var prices []T
	s.q.sqlUpdate = s.q.sqlUpdate.Set(column, value)
	return prices, s.q.QUpdate(ctx, &prices, s.db)
}

func (s baseRepo[T]) WhereID(ids ...int64) baseRepo[T] {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"id": ids})
	return s
}

func (s baseRepo[T]) WhereBase(base string) baseRepo[T] {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"base": base})
	return s
}

func (s baseRepo[T]) OrderBy(by string, order string) baseRepo[T] {
	s.q.sqlSelect = s.q.sqlSelect.OrderBy(fmt.Sprintf("%s %s", by, order))
	return s
}

func (s baseRepo[T]) Limit(u uint64) baseRepo[T] {
	s.q.sqlSelect = s.q.sqlSelect.Limit(u)
	return s
}

func (s baseRepo[T]) Select(ctx context.Context) ([]T, error) {
	var result []T
	return result, s.q.QSelect(ctx, &result, s.db)
}

func (s baseRepo[T]) Get(ctx context.Context) (*T, error) {
	var result T
	return &result, s.q.QGet(ctx, &result, s.db)
}

func (s baseRepo[T]) QDelete(ctx context.Context, q repos.DB) error {
	return s.q.QDelete(ctx, s.db)
}

func (s baseRepo[T]) QCreate(ctx context.Context, q repos.DB) ([]T, error) {
	var result []T
	return result, s.q.QCreate(ctx, &result, s.db)
}

func (s baseRepo[T]) Create(ctx context.Context) ([]T, error) {
	var result []T
	return result, s.q.QCreate(ctx, &result, s.db)
}

func (s baseRepo[T]) QUpdate(ctx context.Context, q repos.DB) ([]T, error) {
	var result []T
	return result, s.q.QSelect(ctx, &result, s.db)
}

func (s baseRepo[T]) QSelect(ctx context.Context, q repos.DB) ([]T, error) {
	var result []T
	return result, s.q.QSelect(ctx, &result, s.db)
}

func (s baseRepo[T]) QGet(ctx context.Context, q repos.DB) (*T, error) {
	var result T
	err := s.q.QGet(ctx, &result, s.db)
	if err != nil {
		if errors.Is(err, sqlerr.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get tender")
	}

	return &result, nil
}
