package market

import (
	"context"

	"prediction/internal/domain"
)

type Service interface {
	Create(ctx context.Context, input CreateInput) (Record, error)
	List(ctx context.Context, filter ListFilter) ([]Record, error)
	ListCatalog(ctx context.Context, limitPerStatus int) (Catalog, error)
	GetByID(ctx context.Context, marketID domain.MarketID) (Record, error)
	GetCreateContext(ctx context.Context) (CreateContext, error)
	ValidateCreateInput(ctx context.Context, input CreateInput) error
	SupportedValidationModels() []ValidationModel
}
