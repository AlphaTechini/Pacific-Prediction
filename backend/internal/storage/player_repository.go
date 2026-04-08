package storage

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type Player struct {
	ID          domain.PlayerID
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreatePlayerInput struct {
	ID          domain.PlayerID
	DisplayName string
}

type PlayerRepository interface {
	Create(ctx context.Context, input CreatePlayerInput) (Player, error)
	GetByID(ctx context.Context, playerID domain.PlayerID) (Player, error)
}
