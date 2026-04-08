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

type Attempt struct {
	MarketID      domain.MarketID
	MarketType    domain.MarketType
	Settled       bool
	SettlementID  *domain.SettlementID
	SettledAt     *time.Time
	SettlementRef string
}

type Service interface {
	SettleMarket(ctx context.Context, marketID domain.MarketID) (Attempt, error)
	SettleDueMarkets(ctx context.Context, filter DueMarketFilter) ([]Attempt, error)
}
