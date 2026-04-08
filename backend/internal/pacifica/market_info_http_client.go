package pacifica

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type HTTPMarketInfoClient struct {
	baseURL    string
	httpClient *http.Client
	cacheTTL   time.Duration

	mu          sync.RWMutex
	cachedAt    time.Time
	cachedItems []MarketInfo
}

type marketInfoResponse struct {
	Success bool                `json:"success"`
	Data    []marketInfoAPIItem `json:"data"`
	Error   any                 `json:"error"`
	Code    any                 `json:"code"`
}

type marketInfoAPIItem struct {
	Symbol          string `json:"symbol"`
	TickSize        string `json:"tick_size"`
	MinTick         string `json:"min_tick"`
	MaxTick         string `json:"max_tick"`
	LotSize         string `json:"lot_size"`
	MaxLeverage     int    `json:"max_leverage"`
	IsolatedOnly    bool   `json:"isolated_only"`
	MinOrderSize    string `json:"min_order_size"`
	MaxOrderSize    string `json:"max_order_size"`
	FundingRate     string `json:"funding_rate"`
	NextFundingRate string `json:"next_funding_rate"`
	CreatedAtMS     int64  `json:"created_at"`
}

type pricesResponse struct {
	Success bool           `json:"success"`
	Data    []priceAPIItem `json:"data"`
	Error   any            `json:"error"`
	Code    any            `json:"code"`
}

type priceAPIItem struct {
	Symbol          string `json:"symbol"`
	MarkPrice       string `json:"mark"`
	MidPrice        string `json:"mid"`
	OraclePrice     string `json:"oracle"`
	FundingRate     string `json:"funding"`
	NextFundingRate string `json:"next_funding"`
	OpenInterest    string `json:"open_interest"`
	Volume24H       string `json:"volume_24h"`
	YesterdayPrice  string `json:"yesterday_price"`
	TimestampMS     int64  `json:"timestamp"`
}

type candlesResponse struct {
	Success bool            `json:"success"`
	Data    []candleAPIItem `json:"data"`
	Error   any             `json:"error"`
	Code    any             `json:"code"`
}

type candleAPIItem struct {
	OpenTimeMS  int64  `json:"t"`
	OpenPrice   string `json:"o"`
	HighPrice   string `json:"h"`
	LowPrice    string `json:"l"`
	ClosePrice  string `json:"c"`
	Volume      string `json:"v"`
	CloseTimeMS int64  `json:"T"`
	TradeCount  int64  `json:"n"`
}

type fundingHistoryResponse struct {
	Success bool             `json:"success"`
	Data    []fundingAPIItem `json:"data"`
	Error   any              `json:"error"`
	Code    any              `json:"code"`
}

type fundingAPIItem struct {
	OraclePrice      string `json:"oracle_price"`
	BidImpactPrice   string `json:"bid_impact_price"`
	AskImpactPrice   string `json:"ask_impact_price"`
	FundingRate      string `json:"funding_rate"`
	NextFundingRate  string `json:"next_funding_rate"`
	SettlementTimeMS int64  `json:"created_at"`
}

func NewHTTPMarketInfoClient(baseURL string, httpClient *http.Client, cacheTTL time.Duration) *HTTPMarketInfoClient {
	return &HTTPMarketInfoClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
		cacheTTL:   cacheTTL,
	}
}

func NewHTTPRESTClient(baseURL string, httpClient *http.Client, cacheTTL time.Duration) *HTTPMarketInfoClient {
	return NewHTTPMarketInfoClient(baseURL, httpClient, cacheTTL)
}

func (c *HTTPMarketInfoClient) ListMarketInfo(ctx context.Context) ([]MarketInfo, error) {
	if cached, ok := c.getCached(); ok {
		return cached, nil
	}

	var payload marketInfoResponse
	if err := c.fetchJSON(ctx, "/api/v1/info", nil, &payload); err != nil {
		return nil, fmt.Errorf("fetch market info: %w", err)
	}
	if !payload.Success {
		return nil, fmt.Errorf("fetch market info: pacifica returned success=false")
	}

	items := make([]MarketInfo, 0, len(payload.Data))
	for _, item := range payload.Data {
		items = append(items, MarketInfo{
			Symbol:          item.Symbol,
			TickSize:        item.TickSize,
			MinTick:         item.MinTick,
			MaxTick:         item.MaxTick,
			LotSize:         item.LotSize,
			MaxLeverage:     item.MaxLeverage,
			IsolatedOnly:    item.IsolatedOnly,
			MinOrderSize:    item.MinOrderSize,
			MaxOrderSize:    item.MaxOrderSize,
			FundingRate:     item.FundingRate,
			NextFundingRate: item.NextFundingRate,
			CreatedAt:       unixMillisUTC(item.CreatedAtMS),
			RawCreatedAtMS:  item.CreatedAtMS,
		})
	}

	c.setCache(items)
	return items, nil
}

func (c *HTTPMarketInfoClient) ListPrices(ctx context.Context, filter PriceFilter) ([]PriceSnapshot, error) {
	var payload pricesResponse
	if err := c.fetchJSON(ctx, "/api/v1/info/prices", nil, &payload); err != nil {
		return nil, fmt.Errorf("fetch prices: %w", err)
	}
	if !payload.Success {
		return nil, fmt.Errorf("fetch prices: pacifica returned success=false")
	}

	symbolFilter := makeSymbolFilter(filter.Symbols)
	items := make([]PriceSnapshot, 0, len(payload.Data))
	for _, item := range payload.Data {
		if len(symbolFilter) > 0 && !symbolFilter[item.Symbol] {
			continue
		}

		items = append(items, PriceSnapshot{
			Symbol:          item.Symbol,
			MarkPrice:       item.MarkPrice,
			MidPrice:        item.MidPrice,
			OraclePrice:     item.OraclePrice,
			FundingRate:     item.FundingRate,
			NextFundingRate: item.NextFundingRate,
			OpenInterest:    item.OpenInterest,
			Volume24H:       item.Volume24H,
			YesterdayPrice:  item.YesterdayPrice,
			Timestamp:       unixMillisUTC(item.TimestampMS),
			RawTimestampMS:  item.TimestampMS,
		})
	}

	return items, nil
}

func (c *HTTPMarketInfoClient) ListMarkPriceCandles(ctx context.Context, query MarkPriceCandleQuery) ([]MarkPriceCandle, error) {
	if strings.TrimSpace(query.Symbol) == "" {
		return nil, fmt.Errorf("fetch mark price candles: symbol is required")
	}
	if strings.TrimSpace(query.Interval) == "" {
		return nil, fmt.Errorf("fetch mark price candles: interval is required")
	}
	if query.StartTime.IsZero() {
		return nil, fmt.Errorf("fetch mark price candles: start time is required")
	}

	params := url.Values{}
	params.Set("symbol", query.Symbol)
	params.Set("interval", query.Interval)
	params.Set("start_time", fmt.Sprintf("%d", query.StartTime.UTC().UnixMilli()))
	if !query.EndTime.IsZero() {
		params.Set("end_time", fmt.Sprintf("%d", query.EndTime.UTC().UnixMilli()))
	}
	if query.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", query.Limit))
	}

	var payload candlesResponse
	if err := c.fetchJSON(ctx, "/api/v1/kline/mark", params, &payload); err != nil {
		return nil, fmt.Errorf("fetch mark price candles: %w", err)
	}
	if !payload.Success {
		return nil, fmt.Errorf("fetch mark price candles: pacifica returned success=false")
	}

	items := make([]MarkPriceCandle, 0, len(payload.Data))
	for _, item := range payload.Data {
		items = append(items, MarkPriceCandle{
			Symbol:         query.Symbol,
			Interval:       query.Interval,
			OpenPrice:      item.OpenPrice,
			HighPrice:      item.HighPrice,
			LowPrice:       item.LowPrice,
			ClosePrice:     item.ClosePrice,
			Volume:         item.Volume,
			TradeCount:     item.TradeCount,
			OpenTime:       unixMillisUTC(item.OpenTimeMS),
			CloseTime:      unixMillisUTC(item.CloseTimeMS),
			RawOpenTimeMS:  item.OpenTimeMS,
			RawCloseTimeMS: item.CloseTimeMS,
		})
	}

	return items, nil
}

func (c *HTTPMarketInfoClient) ListFundingRateHistory(ctx context.Context, query FundingRateHistoryQuery) ([]FundingRateHistoryEntry, error) {
	if strings.TrimSpace(query.Symbol) == "" {
		return nil, fmt.Errorf("fetch funding history: symbol is required")
	}

	params := url.Values{}
	params.Set("symbol", query.Symbol)
	if query.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", query.Limit))
	}

	var payload fundingHistoryResponse
	if err := c.fetchJSON(ctx, "/api/v1/funding_rate/history", params, &payload); err != nil {
		return nil, fmt.Errorf("fetch funding history: %w", err)
	}
	if !payload.Success {
		return nil, fmt.Errorf("fetch funding history: pacifica returned success=false")
	}

	items := make([]FundingRateHistoryEntry, 0, len(payload.Data))
	for _, item := range payload.Data {
		settlementTime := unixMillisUTC(item.SettlementTimeMS)
		if !query.SettlementTimeAfter.IsZero() && settlementTime.Before(query.SettlementTimeAfter.UTC()) {
			continue
		}
		if !query.SettlementTimeBefore.IsZero() && settlementTime.After(query.SettlementTimeBefore.UTC()) {
			continue
		}

		items = append(items, FundingRateHistoryEntry{
			Symbol:              query.Symbol,
			FundingRate:         item.FundingRate,
			NextFundingRate:     item.NextFundingRate,
			OraclePrice:         item.OraclePrice,
			BidImpactPrice:      item.BidImpactPrice,
			AskImpactPrice:      item.AskImpactPrice,
			SettlementTime:      settlementTime,
			RawSettlementTimeMS: item.SettlementTimeMS,
		})
	}

	return items, nil
}

func (c *HTTPMarketInfoClient) getCached() ([]MarketInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cacheTTL <= 0 || c.cachedAt.IsZero() || time.Since(c.cachedAt) > c.cacheTTL {
		return nil, false
	}

	return append([]MarketInfo(nil), c.cachedItems...), true
}

func (c *HTTPMarketInfoClient) setCache(items []MarketInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cachedAt = time.Now().UTC()
	c.cachedItems = append([]MarketInfo(nil), items...)
}

func (c *HTTPMarketInfoClient) fetchJSON(ctx context.Context, path string, params url.Values, target any) error {
	requestURL := c.baseURL + path
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func makeSymbolFilter(symbols []string) map[string]bool {
	if len(symbols) == 0 {
		return nil
	}

	items := make(map[string]bool, len(symbols))
	for _, symbol := range symbols {
		trimmed := strings.TrimSpace(symbol)
		if trimmed == "" {
			continue
		}

		items[trimmed] = true
	}

	return items
}

func unixMillisUTC(value int64) time.Time {
	if value <= 0 {
		return time.Time{}
	}

	return time.UnixMilli(value).UTC()
}
