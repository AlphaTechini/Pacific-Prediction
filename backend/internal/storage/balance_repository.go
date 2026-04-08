package storage

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type Balance struct {
	PlayerID         domain.PlayerID
	AvailableBalance string
	LockedBalance    string
	UpdatedAt        time.Time
}

type CreateBalanceInput struct {
	PlayerID         domain.PlayerID
	AvailableBalance string
	LockedBalance    string
}

type LockStakeInput struct {
	PlayerID domain.PlayerID
	Amount   string
}

type BalanceRepository interface {
	Create(ctx context.Context, input CreateBalanceInput) (Balance, error)
	GetByPlayerID(ctx context.Context, playerID domain.PlayerID) (Balance, error)
	LockStake(ctx context.Context, input LockStakeInput) (Balance, error)
}
