package market

import "time"

type CreateContext struct {
	Symbols                           []CreateContextSymbol
	ValidationModels                  []ValidationModel
	PriceThresholdCreationBandPercent string
}

type CreateContextSymbol struct {
	Symbol          string
	TickSize        string
	MinTick         string
	MaxTick         string
	LotSize         string
	MinOrderSize    string
	MaxOrderSize    string
	MaxLeverage     int
	IsolatedOnly    bool
	MarkPrice       string
	OraclePrice     string
	FundingRate     string
	NextFundingRate string
	OpenInterest    string
	Volume24H       string
	UpdatedAt       time.Time
}
