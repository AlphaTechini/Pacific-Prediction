package balance

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type Snapshot struct {
	PlayerID         domain.PlayerID
	AvailableBalance string
	LockedBalance    string
	UpdatedAt        time.Time
}

type Service interface {
	GetBalance(ctx context.Context, playerID domain.PlayerID) (Snapshot, error)
	LockStake(ctx context.Context, playerID domain.PlayerID, amount string) error
	UnlockStake(ctx context.Context, playerID domain.PlayerID, amount string) error
	Credit(ctx context.Context, playerID domain.PlayerID, amount string) error
	Debit(ctx context.Context, playerID domain.PlayerID, amount string) error
}

type Controller interface {
	GetBalance(ctx context.Context, playerID domain.PlayerID) (Snapshot, error)
}
