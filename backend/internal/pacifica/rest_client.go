package pacifica

import "context"

type RESTClient interface {
	ListMarketInfo(ctx context.Context) ([]MarketInfo, error)
	ListPrices(ctx context.Context, filter PriceFilter) ([]PriceSnapshot, error)
	ListMarkPriceCandles(ctx context.Context, query MarkPriceCandleQuery) ([]MarkPriceCandle, error)
	ListFundingRateHistory(ctx context.Context, query FundingRateHistoryQuery) ([]FundingRateHistoryEntry, error)
}
