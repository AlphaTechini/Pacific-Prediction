package pacifica

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
	CreatedAt       string
}

type MarketInfoClient interface {
	ListMarketInfo() ([]MarketInfo, error)
}
