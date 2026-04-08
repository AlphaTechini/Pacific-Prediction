package pacifica

import "time"

type FundingRateHistoryQuery struct {
	Symbol               string
	SettlementTimeAfter  time.Time
	SettlementTimeBefore time.Time
	Limit                int
}

type FundingRateHistoryEntry struct {
	Symbol              string
	FundingRate         string
	NextFundingRate     string
	OraclePrice         string
	BidImpactPrice      string
	AskImpactPrice      string
	SettlementTime      time.Time
	RawSettlementTimeMS int64
}
