package player

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

func (c *controller) GetMe(ctx context.Context, playerID domain.PlayerID) (Profile, error) {
	return c.service.GetProfile(ctx, playerID)
}
