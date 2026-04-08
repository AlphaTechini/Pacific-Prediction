package player

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type Profile struct {
	ID          domain.PlayerID
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Service interface {
	GetProfile(ctx context.Context, playerID domain.PlayerID) (Profile, error)
}

type Controller interface {
	GetMe(ctx context.Context, playerID domain.PlayerID) (Profile, error)
}
