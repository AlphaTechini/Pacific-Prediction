package pacifica

import (
	"context"
	"time"
)

type MarketInfoClient interface {
	ListMarketInfo(ctx context.Context) ([]MarketInfo, error)
}

type MarketInfo struct {
	Symbol          string
	TickSize        string
	MinTick         string
	MaxTick         string
	LotSize         string
	MaxLeverage     int
	IsolatedOnly    bool
	MinOrderSize    string
	MaxOrderSize    string
	FundingRate     string
	NextFundingRate string
	CreatedAt       time.Time
	RawCreatedAtMS  int64
}
