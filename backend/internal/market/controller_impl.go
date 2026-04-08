package market

import (
	"context"

	"prediction/internal/domain"
)

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service: service}
}

func (c *controller) Create(ctx context.Context, input CreateInput) (Record, error) {
	return c.service.Create(ctx, input)
}

func (c *controller) List(ctx context.Context, filter ListFilter) ([]Record, error) {
	return c.service.List(ctx, filter)
}

func (c *controller) ListCatalog(ctx context.Context, limitPerStatus int) (Catalog, error) {
	return c.service.ListCatalog(ctx, limitPerStatus)
}

func (c *controller) GetByID(ctx context.Context, marketID domain.MarketID) (Record, error) {
	return c.service.GetByID(ctx, marketID)
}

func (c *controller) ValidateCreateInput(ctx context.Context, input CreateInput) error {
	return c.service.ValidateCreateInput(ctx, input)
}

func (c *controller) SupportedValidationModels() []ValidationModel {
	return c.service.SupportedValidationModels()
}
