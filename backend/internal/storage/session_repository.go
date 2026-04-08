package storage

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

type CreateSessionInput struct {
	ID               domain.SessionID
	PlayerID         domain.PlayerID
	SessionTokenHash string
	ExpiresAt        time.Time
}

type SessionRepository interface {
	Create(ctx context.Context, input CreateSessionInput) (Session, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (Session, error)
}
