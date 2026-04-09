package market

import "prediction/internal/domain"

type ValidationModel struct {
	MarketType        domain.MarketType
	SourceType        domain.SourceType
	AllowedOperators  []domain.ConditionOperator
	RequiresThreshold bool
	RequiresInterval  bool
	AllowedIntervals  []string
}

type PriceThresholdValidation struct {
	MarketType        domain.MarketType
	SourceType        domain.SourceType
	AllowedOperators  []domain.ConditionOperator
	RequiresThreshold bool
}

type CandleDirectionValidation struct {
	MarketType       domain.MarketType
	SourceType       domain.SourceType
	AllowedOperators []domain.ConditionOperator
	RequiresInterval bool
	AllowedIntervals []string
}

type FundingThresholdValidation struct {
	MarketType        domain.MarketType
	SourceType        domain.SourceType
	AllowedOperators  []domain.ConditionOperator
	RequiresThreshold bool
	RequiresInterval  bool
	SourceInterval    string
}

func SupportedValidationModels() []ValidationModel {
	priceValidation := SupportedPriceThresholdValidation()
	candleValidation := SupportedCandleDirectionValidation()
	fundingValidation := SupportedFundingThresholdValidation()

	return []ValidationModel{
		{
			MarketType:        priceValidation.MarketType,
			SourceType:        priceValidation.SourceType,
			AllowedOperators:  priceValidation.AllowedOperators,
			RequiresThreshold: priceValidation.RequiresThreshold,
		},
		{
			MarketType:       candleValidation.MarketType,
			SourceType:       candleValidation.SourceType,
			AllowedOperators: candleValidation.AllowedOperators,
			RequiresInterval: candleValidation.RequiresInterval,
			AllowedIntervals: candleValidation.AllowedIntervals,
		},
		{
			MarketType:        fundingValidation.MarketType,
			SourceType:        fundingValidation.SourceType,
			AllowedOperators:  fundingValidation.AllowedOperators,
			RequiresThreshold: fundingValidation.RequiresThreshold,
			RequiresInterval:  fundingValidation.RequiresInterval,
			AllowedIntervals:  []string{fundingValidation.SourceInterval},
		},
	}
}

func SupportedPriceThresholdValidation() PriceThresholdValidation {
	return PriceThresholdValidation{
		MarketType: domain.MarketTypePriceThreshold,
		SourceType: domain.SourceTypeMarkPrice,
		AllowedOperators: []domain.ConditionOperator{
			domain.ConditionOperatorGT,
			domain.ConditionOperatorGTE,
			domain.ConditionOperatorLT,
			domain.ConditionOperatorLTE,
		},
		RequiresThreshold: true,
	}
}

func SupportedCandleDirectionValidation() CandleDirectionValidation {
	return CandleDirectionValidation{
		MarketType: domain.MarketTypeCandleDirection,
		SourceType: domain.SourceTypeMarkPriceCandle,
		AllowedOperators: []domain.ConditionOperator{
			domain.ConditionOperatorBullishClose,
			domain.ConditionOperatorBearishClose,
		},
		RequiresInterval: true,
		AllowedIntervals: []string{
			"1m",
			"3m",
			"5m",
			"15m",
			"30m",
			"1h",
			"2h",
			"4h",
			"8h",
			"12h",
			"1d",
		},
	}
}

func SupportedFundingThresholdValidation() FundingThresholdValidation {
	return FundingThresholdValidation{
		MarketType: domain.MarketTypeFundingThreshold,
		SourceType: domain.SourceTypeFundingRate,
		AllowedOperators: []domain.ConditionOperator{
			domain.ConditionOperatorGT,
			domain.ConditionOperatorGTE,
			domain.ConditionOperatorLT,
			domain.ConditionOperatorLTE,
			domain.ConditionOperatorPositive,
			domain.ConditionOperatorNegative,
		},
		RequiresThreshold: true,
		RequiresInterval:  true,
		SourceInterval:    domain.SourceIntervalFundingEpoch,
	}
}
