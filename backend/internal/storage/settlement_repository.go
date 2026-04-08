package storage

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type Settlement struct {
	ID              domain.SettlementID
	MarketID        domain.MarketID
	PacificaSource  string
	SourceTimestamp time.Time
	RawPayload      []byte
	SettlementValue string
	Result          domain.MarketResult
	CreatedAt       time.Time
}

type CreateSettlementInput struct {
	ID              domain.SettlementID
	MarketID        domain.MarketID
	PacificaSource  string
	SourceTimestamp time.Time
	RawPayload      []byte
	SettlementValue string
	Result          domain.MarketResult
}

type SettlementRepository interface {
	Create(ctx context.Context, input CreateSettlementInput) (Settlement, error)
	GetByMarketID(ctx context.Context, marketID domain.MarketID) (Settlement, error)
}
