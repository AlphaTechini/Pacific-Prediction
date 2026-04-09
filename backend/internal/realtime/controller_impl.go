package realtime

import "context"

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service: service}
}

func (c *controller) Subscribe(ctx context.Context) (Subscription, error) {
	return c.service.Subscribe(ctx)
}
