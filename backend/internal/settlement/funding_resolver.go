package settlement

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type FundingMarket struct {
	ID                domain.MarketID
	Symbol            string
	ConditionOperator domain.ConditionOperator
	ThresholdValue    string
	ExpiryTime        time.Time
}

type FundingResolution struct {
	MarketID       domain.MarketID
	PacificaSource string
	SettlementTime time.Time
	FundingRate    string
	Result         domain.MarketResult
	RawPayload     []byte
}

type FundingResolver interface {
	Resolve(ctx context.Context, market FundingMarket) (FundingResolution, error)
}
