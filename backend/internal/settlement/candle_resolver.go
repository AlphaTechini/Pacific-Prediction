package settlement

import (
	"context"
	"time"

	"prediction/internal/domain"
)

type CandleMarket struct {
	ID                domain.MarketID
	Symbol            string
	ConditionOperator domain.ConditionOperator
	SourceInterval    string
	ExpiryTime        time.Time
}

type CandleResolution struct {
	MarketID       domain.MarketID
	PacificaSource string
	OpenTime       time.Time
	CloseTime      time.Time
	OpenPrice      string
	ClosePrice     string
	Result         domain.MarketResult
	RawPayload     []byte
}

type CandleResolver interface {
	Resolve(ctx context.Context, market CandleMarket) (CandleResolution, error)
}
