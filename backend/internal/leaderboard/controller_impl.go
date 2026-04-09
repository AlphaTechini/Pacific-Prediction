package leaderboard

import "context"

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service: service}
}

func (c *controller) GetSnapshot(ctx context.Context, limit int) (Snapshot, error) {
	return c.service.GetSnapshot(ctx, limit)
}
