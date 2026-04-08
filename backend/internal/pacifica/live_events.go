package pacifica

import "time"

type LiveEvent interface {
	Channel() SubscriptionChannel
	SymbolKey() string
	EventTime() time.Time
}

type PriceStreamEvent struct {
	Symbol         string
	MarkPrice      string
	OraclePrice    string
	FundingRate    string
	OpenInterest   string
	Volume24H      string
	Timestamp      time.Time
	RawTimestampMS int64
}

func (e PriceStreamEvent) Channel() SubscriptionChannel {
	return SubscriptionChannelPrices
}

func (e PriceStreamEvent) SymbolKey() string {
	return e.Symbol
}

func (e PriceStreamEvent) EventTime() time.Time {
	return e.Timestamp
}

type CandleStreamEvent struct {
	Symbol         string
	Interval       string
	OpenPrice      string
	HighPrice      string
	LowPrice       string
	ClosePrice     string
	Volume         string
	OpenTime       time.Time
	CloseTime      time.Time
	RawOpenTimeMS  int64
	RawCloseTimeMS int64
}

func (e CandleStreamEvent) Channel() SubscriptionChannel {
	return SubscriptionChannelCandle
}

func (e CandleStreamEvent) SymbolKey() string {
	return e.Symbol
}

func (e CandleStreamEvent) EventTime() time.Time {
	return e.CloseTime
}

type MarkPriceCandleStreamEvent struct {
	Symbol         string
	Interval       string
	OpenPrice      string
	HighPrice      string
	LowPrice       string
	ClosePrice     string
	Volume         string
	OpenTime       time.Time
	CloseTime      time.Time
	RawOpenTimeMS  int64
	RawCloseTimeMS int64
}

func (e MarkPriceCandleStreamEvent) Channel() SubscriptionChannel {
	return SubscriptionChannelMarkPriceCandle
}

func (e MarkPriceCandleStreamEvent) SymbolKey() string {
	return e.Symbol
}

func (e MarkPriceCandleStreamEvent) EventTime() time.Time {
	return e.CloseTime
}

type TradeStreamEvent struct {
	Symbol         string
	Side           string
	Cause          string
	Price          string
	Amount         string
	LastOrderID    string
	Timestamp      time.Time
	RawTimestampMS int64
}

func (e TradeStreamEvent) Channel() SubscriptionChannel {
	return SubscriptionChannelTrades
}

func (e TradeStreamEvent) SymbolKey() string {
	return e.Symbol
}

func (e TradeStreamEvent) EventTime() time.Time {
	return e.Timestamp
}
