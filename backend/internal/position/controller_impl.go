package position

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

func (c *controller) Create(ctx context.Context, playerID domain.PlayerID, input CreateInput) (Record, error) {
	return c.service.Create(ctx, playerID, input)
}

func (c *controller) ListByPlayerID(ctx context.Context, playerID domain.PlayerID, filter ListFilter) ([]Record, error) {
	return c.service.ListByPlayerID(ctx, playerID, filter)
}

func (c *controller) ValidateCreateInput(ctx context.Context, playerID domain.PlayerID, input CreateInput) error {
	return c.service.ValidateCreateInput(ctx, playerID, input)
}
