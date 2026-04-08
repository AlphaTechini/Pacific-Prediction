package storage

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type Position struct {
	ID              domain.PositionID
	PlayerID        domain.PlayerID
	MarketID        domain.MarketID
	Side            domain.PositionSide
	StakeAmount     string
	PotentialPayout string
	Status          domain.PositionStatus
	CreatedAt       time.Time
	SettledAt       *time.Time
}

type CreatePositionInput struct {
	ID              domain.PositionID
	PlayerID        domain.PlayerID
	MarketID        domain.MarketID
	Side            domain.PositionSide
	StakeAmount     string
	PotentialPayout string
}

type UpdatePositionSettlementInput struct {
	PositionID domain.PositionID
	Status     domain.PositionStatus
	SettledAt  time.Time
}

type PositionRepository interface {
	Create(ctx context.Context, input CreatePositionInput) (Position, error)
	ListByPlayerID(ctx context.Context, playerID domain.PlayerID, limit int) ([]Position, error)
	ListByMarketID(ctx context.Context, marketID domain.MarketID) ([]Position, error)
	UpdateSettlement(ctx context.Context, input UpdatePositionSettlementInput) (Position, error)
}
