package auth

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type Session struct {
	ID               domain.SessionID
	PlayerID         domain.PlayerID
	SessionTokenHash string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

type CreateGuestSessionInput struct {
	DisplayName string
}

type AuthContext struct {
	PlayerID  domain.PlayerID
	SessionID domain.SessionID
}

type Service interface {
	CreateGuestSession(ctx context.Context, input CreateGuestSessionInput) (Session, error)
	ValidateSession(ctx context.Context, rawToken string) (AuthContext, error)
}

type Controller interface {
	CreateGuestSession(ctx context.Context, input CreateGuestSessionInput) (Session, error)
}
