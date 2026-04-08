package domain

type MarketType string

const (
	MarketTypePriceThreshold   MarketType = "price_threshold"
	MarketTypeCandleDirection  MarketType = "candle_direction"
	MarketTypeFundingThreshold MarketType = "funding_threshold"
)

type SourceType string

const (
	SourceTypeMarkPrice       SourceType = "mark_price"
	SourceTypeMarkPriceCandle SourceType = "mark_price_candle"
	SourceTypeFundingRate     SourceType = "funding_rate"
)

const SourceIntervalFundingEpoch = "funding_epoch"

type ConditionOperator string

const (
	ConditionOperatorGT           ConditionOperator = "gt"
	ConditionOperatorGTE          ConditionOperator = "gte"
	ConditionOperatorLT           ConditionOperator = "lt"
	ConditionOperatorLTE          ConditionOperator = "lte"
	ConditionOperatorBullishClose ConditionOperator = "bullish_close"
	ConditionOperatorBearishClose ConditionOperator = "bearish_close"
	ConditionOperatorPositive     ConditionOperator = "positive"
	ConditionOperatorNegative     ConditionOperator = "negative"
)

type MarketStatus string

const (
	MarketStatusActive    MarketStatus = "active"
	MarketStatusResolving MarketStatus = "resolving"
	MarketStatusResolved  MarketStatus = "resolved"
	MarketStatusCancelled MarketStatus = "cancelled"
)

type MarketResult string

const (
	MarketResultYes MarketResult = "yes"
	MarketResultNo  MarketResult = "no"
)
