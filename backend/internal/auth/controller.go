package auth

import "context"

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service: service}
}

func (c *controller) CreateGuestSession(ctx context.Context, input CreateGuestSessionInput) (IssuedSession, error) {
	return c.service.CreateGuestSession(ctx, input)
}

func (c *controller) ValidateSession(ctx context.Context, rawToken string) (AuthContext, error) {
	return c.service.ValidateSession(ctx, rawToken)
}
