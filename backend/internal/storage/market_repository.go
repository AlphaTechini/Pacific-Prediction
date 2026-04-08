package storage

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type Market struct {
	ID                domain.MarketID
	Title             string
	Symbol            string
	MarketType        domain.MarketType
	ConditionOperator domain.ConditionOperator
	ThresholdValue    string
	SourceType        domain.SourceType
	SourceInterval    string
	ReferenceValue    string
	ExpiryTime        time.Time
	Status            domain.MarketStatus
	Result            domain.MarketResult
	SettlementValue   string
	ResolvedAt        *time.Time
	ResolutionReason  string
	CreatedByPlayerID domain.PlayerID
	CreatedAt         time.Time
}

type CreateMarketInput struct {
	ID                domain.MarketID
	Title             string
	Symbol            string
	MarketType        domain.MarketType
	ConditionOperator domain.ConditionOperator
	ThresholdValue    string
	SourceType        domain.SourceType
	SourceInterval    string
	ReferenceValue    string
	ExpiryTime        time.Time
	CreatedByPlayerID domain.PlayerID
}

type UpdateMarketSettlementInput struct {
	MarketID         domain.MarketID
	Status           domain.MarketStatus
	Result           domain.MarketResult
	SettlementValue  string
	ResolvedAt       time.Time
	ResolutionReason string
}

type MarketRepository interface {
	Create(ctx context.Context, input CreateMarketInput) (Market, error)
	GetByID(ctx context.Context, marketID domain.MarketID) (Market, error)
	ListByStatus(ctx context.Context, status domain.MarketStatus, limit int) ([]Market, error)
	ListExpiringBefore(ctx context.Context, before time.Time, limit int) ([]Market, error)
	UpdateSettlement(ctx context.Context, input UpdateMarketSettlementInput) (Market, error)
}
