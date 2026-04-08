package market

import (
	"context"

	"prediction/internal/domain"
)

type Controller interface {
	Create(ctx context.Context, input CreateInput) (Record, error)
	List(ctx context.Context, filter ListFilter) ([]Record, error)
	GetByID(ctx context.Context, marketID domain.MarketID) (Record, error)
	SupportedValidationModels() []ValidationModel
}
