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

type SettleWonPositionInput struct {
	PlayerID     domain.PlayerID
	StakeAmount  string
	PayoutAmount string
}

type SettleLostPositionInput struct {
	PlayerID    domain.PlayerID
	StakeAmount string
}

type BalanceRepository interface {
	Create(ctx context.Context, input CreateBalanceInput) (Balance, error)
	GetByPlayerID(ctx context.Context, playerID domain.PlayerID) (Balance, error)
	LockStake(ctx context.Context, input LockStakeInput) (Balance, error)
	SettleWonPosition(ctx context.Context, input SettleWonPositionInput) (Balance, error)
	SettleLostPosition(ctx context.Context, input SettleLostPositionInput) (Balance, error)
}
