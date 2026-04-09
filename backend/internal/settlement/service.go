package settlement

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type DueMarketFilter struct {
	Before time.Time
	Limit  int
}

type PriceFetchPlanFilter struct {
	After  time.Time
	Before time.Time
	Limit  int
}

type Attempt struct {
	MarketID      domain.MarketID
	MarketType    domain.MarketType
	Settled       bool
	SettlementID  *domain.SettlementID
	SettledAt     *time.Time
	SettlementRef string
}

type PriceFetchTarget struct {
	MarketID   domain.MarketID
	Symbol     string
	ExpiryTime time.Time
}

type PriceFetchBatch struct {
	ExpiryTime time.Time
	Symbols    []string
	Targets    []PriceFetchTarget
}

type Service interface {
	SettleMarket(ctx context.Context, marketID domain.MarketID) (Attempt, error)
	SettleDueMarkets(ctx context.Context, filter DueMarketFilter) ([]Attempt, error)
	PlanPriceFetchBatches(ctx context.Context, filter PriceFetchPlanFilter) ([]PriceFetchBatch, error)
}
