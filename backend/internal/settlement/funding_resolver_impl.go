package settlement

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"prediction/internal/domain"
	"prediction/internal/pacifica"
)

const pacificaFundingHistorySource = "/api/v1/funding_rate/history"

type fundingResolver struct {
	pacificaClient pacifica.RESTClient
}

func NewFundingResolver(pacificaClient pacifica.RESTClient) FundingResolver {
	return &fundingResolver{
		pacificaClient: pacificaClient,
	}
}

func (r *fundingResolver) Resolve(ctx context.Context, market FundingMarket) (FundingResolution, error) {
	if strings.TrimSpace(market.Symbol) == "" {
		return FundingResolution{}, domain.NewValidationError("symbol", "symbol is required for funding settlement", market.Symbol)
	}

	items, err := r.pacificaClient.ListFundingRateHistory(ctx, pacifica.FundingRateHistoryQuery{
		Symbol:              market.Symbol,
		SettlementTimeAfter: market.ExpiryTime,
		Limit:               5,
	})
	if err != nil {
		return FundingResolution{}, fmt.Errorf("fetch funding history for market settlement: %w", err)
	}

	entry, err := findFundingSettlementEntry(items, market.Symbol, market.ExpiryTime)
	if err != nil {
		return FundingResolution{}, err
	}

	result, err := compareFundingRate(entry.FundingRate, market.ThresholdValue, market.ConditionOperator)
	if err != nil {
		return FundingResolution{}, err
	}

	rawPayload, err := json.Marshal(entry)
	if err != nil {
		return FundingResolution{}, fmt.Errorf("marshal funding history payload: %w", err)
	}

	return FundingResolution{
		MarketID:       market.ID,
		PacificaSource: pacificaFundingHistorySource,
		SettlementTime: entry.SettlementTime,
		FundingRate:    entry.FundingRate,
		Result:         result,
		RawPayload:     rawPayload,
	}, nil
}

func findFundingSettlementEntry(items []pacifica.FundingRateHistoryEntry, symbol string, expiry time.Time) (pacifica.FundingRateHistoryEntry, error) {
	normalizedExpiry := domain.NormalizeTime(expiry)
	for _, item := range items {
		if !strings.EqualFold(item.Symbol, symbol) {
			continue
		}
		if domain.NormalizeTime(item.SettlementTime).Before(normalizedExpiry) {
			continue
		}

		return item, nil
	}

	return pacifica.FundingRateHistoryEntry{}, fmt.Errorf(
		"%w: funding record is not yet available for symbol=%s expiry=%s",
		errSettlementSourceNotReady,
		symbol,
		normalizedExpiry.Format(time.RFC3339Nano),
	)
}

func compareFundingRate(actualValue string, thresholdValue string, operator domain.ConditionOperator) (domain.MarketResult, error) {
	actual, err := parseSettlementDecimal(actualValue)
	if err != nil {
		return "", domain.NewValidationError("funding_rate", "funding rate must be a valid decimal string", actualValue)
	}

	switch operator {
	case domain.ConditionOperatorPositive:
		if actual.Sign() > 0 {
			return domain.MarketResultYes, nil
		}
		return domain.MarketResultNo, nil
	case domain.ConditionOperatorNegative:
		if actual.Sign() < 0 {
			return domain.MarketResultYes, nil
		}
		return domain.MarketResultNo, nil
	}

	threshold, err := parseSettlementDecimal(thresholdValue)
	if err != nil {
		return "", domain.NewValidationError("threshold_value", "market threshold must be a valid decimal string", thresholdValue)
	}

	comparison := actual.Cmp(threshold)

	switch operator {
	case domain.ConditionOperatorGT:
		if comparison > 0 {
			return domain.MarketResultYes, nil
		}
	case domain.ConditionOperatorGTE:
		if comparison >= 0 {
			return domain.MarketResultYes, nil
		}
	case domain.ConditionOperatorLT:
		if comparison < 0 {
			return domain.MarketResultYes, nil
		}
	case domain.ConditionOperatorLTE:
		if comparison <= 0 {
			return domain.MarketResultYes, nil
		}
	default:
		return "", domain.NewValidationError("condition_operator", "funding settlement operator is not supported", operator)
	}

	return domain.MarketResultNo, nil
}
