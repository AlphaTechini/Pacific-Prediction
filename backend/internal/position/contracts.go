package position

import (
	"time"

	"prediction/internal/domain"
)

type Record struct {
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

type CreateInput struct {
	MarketID    domain.MarketID
	Side        domain.PositionSide
	StakeAmount string
}

type ListFilter struct {
	Limit int
}
