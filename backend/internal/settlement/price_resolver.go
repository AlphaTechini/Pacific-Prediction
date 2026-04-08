package settlement

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type PriceMarket struct {
	ID                domain.MarketID
	Symbol            string
	ConditionOperator domain.ConditionOperator
	ThresholdValue    string
	ExpiryTime        time.Time
}

type PriceResolution struct {
	MarketID            domain.MarketID
	PacificaSource      string
	SourceTimestamp     time.Time
	SettlementMarkPrice string
	Result              domain.MarketResult
	RawPayload          []byte
}

type PriceResolver interface {
	Resolve(ctx context.Context, market PriceMarket) (PriceResolution, error)
}
