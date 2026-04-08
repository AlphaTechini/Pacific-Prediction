package pacifica

import "time"

type MarkPriceCandleQuery struct {
	Symbol    string
	Interval  string
	StartTime time.Time
	EndTime   time.Time
	Limit     int
}

type MarkPriceCandle struct {
	Symbol         string
	Interval       string
	OpenPrice      string
	HighPrice      string
	LowPrice       string
	ClosePrice     string
	Volume         string
	TradeCount     int64
	OpenTime       time.Time
	CloseTime      time.Time
	RawOpenTimeMS  int64
	RawCloseTimeMS int64
}
