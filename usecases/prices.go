package usecases

import (
	"context"

	"github.com/ow0sh/gotest/models"
	"github.com/ow0sh/gotest/repos"
	"github.com/pkg/errors"
)

type PricesUseCase struct {
	repo repos.PricesRepo
}

func NewPricesUseCase(repo repos.PricesRepo) *PricesUseCase {
	return &PricesUseCase{repo: repo}
}

func (use PricesUseCase) CreatePrices(ctx context.Context, prices ...models.Price) ([]repos.Price, error) {
	pricesDB := make([]repos.Price, 0, len(prices))
	for _, price := range prices {
		pricesDB, err := use.repo.Inserter().SetCreate(repos.PricesToCreateRepo(price)...).Create(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create price")
		}
		pricesDB = append(pricesDB, pricesDB[0])
	}

	return pricesDB, nil
}

func (use PricesUseCase) UpdatePrices(ctx context.Context, base string, rate float64) ([]repos.Price, error) {
	pricesDB, err := use.repo.Updater().WhereBase(base).Update(ctx, "rate", rate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update price")
	}
	return pricesDB, nil
}
