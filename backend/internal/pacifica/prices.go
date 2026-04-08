package pacifica

import "time"

type PriceFilter struct {
	Symbols []string
}

type PriceSnapshot struct {
	Symbol          string
	MarkPrice       string
	MidPrice        string
	OraclePrice     string
	FundingRate     string
	NextFundingRate string
	OpenInterest    string
	Volume24H       string
	YesterdayPrice  string
	Timestamp       time.Time
	RawTimestampMS  int64
}
