package balance

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

func (c *controller) GetBalance(ctx context.Context, playerID domain.PlayerID) (Snapshot, error) {
	return c.service.GetBalance(ctx, playerID)
}
