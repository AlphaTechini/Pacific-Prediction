package settlement

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"prediction/internal/domain"
	"prediction/internal/pacifica"
)

const pacificaPricesSource = "/api/v1/info/prices"

type priceResolver struct {
	pacificaClient pacifica.RESTClient
}

func NewPriceResolver(pacificaClient pacifica.RESTClient) PriceResolver {
	return &priceResolver{
		pacificaClient: pacificaClient,
	}
}

func (r *priceResolver) Resolve(ctx context.Context, market PriceMarket) (PriceResolution, error) {
	items, err := r.fetchPrices(ctx, []string{market.Symbol})
	if err != nil {
		return PriceResolution{}, fmt.Errorf("fetch price for market settlement: %w", err)
	}

	return resolvePriceSnapshot(market, items)
}

func (r *priceResolver) ResolveBatch(ctx context.Context, markets []PriceMarket) ([]PriceResolution, error) {
	if len(markets) == 0 {
		return nil, nil
	}

	items, err := r.fetchPrices(ctx, uniqueMarketSymbols(markets))
	if err != nil {
		return nil, fmt.Errorf("fetch batch prices for market settlement: %w", err)
	}

	resolutions := make([]PriceResolution, 0, len(markets))
	for _, market := range markets {
		resolution, err := resolvePriceSnapshot(market, items)
		if err != nil {
			return nil, err
		}

		resolutions = append(resolutions, resolution)
	}

	return resolutions, nil
}

func (r *priceResolver) fetchPrices(ctx context.Context, symbols []string) ([]pacifica.PriceSnapshot, error) {
	return r.pacificaClient.ListPrices(ctx, pacifica.PriceFilter{
		Symbols: symbols,
	})
}

func resolvePriceSnapshot(market PriceMarket, items []pacifica.PriceSnapshot) (PriceResolution, error) {
	if strings.TrimSpace(market.Symbol) == "" {
		return PriceResolution{}, domain.NewValidationError("symbol", "symbol is required for price settlement", market.Symbol)
	}

	snapshot, err := findPriceSnapshot(items, market.Symbol)
	if err != nil {
		return PriceResolution{}, err
	}

	if snapshot.Timestamp.Before(domain.NormalizeTime(market.ExpiryTime)) {
		return PriceResolution{}, fmt.Errorf(
			"%w: price snapshot predates market expiry: snapshot=%s expiry=%s",
			errSettlementSourceNotReady,
			snapshot.Timestamp.Format(time.RFC3339Nano),
			domain.NormalizeTime(market.ExpiryTime).Format(time.RFC3339Nano),
		)
	}

	result, err := compareThreshold(snapshot.MarkPrice, market.ThresholdValue, market.ConditionOperator)
	if err != nil {
		return PriceResolution{}, err
	}

	rawPayload, err := json.Marshal(snapshot)
	if err != nil {
		return PriceResolution{}, fmt.Errorf("marshal price snapshot payload: %w", err)
	}

	return PriceResolution{
		MarketID:            market.ID,
		PacificaSource:      pacificaPricesSource,
		SourceTimestamp:     snapshot.Timestamp,
		SettlementMarkPrice: snapshot.MarkPrice,
		Result:              result,
		RawPayload:          rawPayload,
	}, nil
}

func findPriceSnapshot(items []pacifica.PriceSnapshot, symbol string) (pacifica.PriceSnapshot, error) {
	target := strings.TrimSpace(symbol)
	for _, item := range items {
		if strings.EqualFold(item.Symbol, target) {
			return item, nil
		}
	}

	return pacifica.PriceSnapshot{}, domain.NewValidationError("symbol", "price snapshot is missing for settlement symbol", symbol)
}

func compareThreshold(actualValue, thresholdValue string, operator domain.ConditionOperator) (domain.MarketResult, error) {
	actual, err := parseSettlementDecimal(actualValue)
	if err != nil {
		return "", domain.NewValidationError("mark_price", "price snapshot mark value must be a valid decimal string", actualValue)
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
		return "", domain.NewValidationError("condition_operator", "price settlement operator is not supported", operator)
	}

	return domain.MarketResultNo, nil
}

func parseSettlementDecimal(value string) (*big.Rat, error) {
	parsed, ok := new(big.Rat).SetString(strings.TrimSpace(value))
	if !ok {
		return nil, fmt.Errorf("parse settlement decimal %q", value)
	}

	return parsed, nil
}

func uniqueMarketSymbols(markets []PriceMarket) []string {
	items := make([]string, 0, len(markets))
	seen := make(map[string]struct{}, len(markets))
	for _, market := range markets {
		symbol := strings.TrimSpace(market.Symbol)
		if _, ok := seen[symbol]; ok {
			continue
		}

		seen[symbol] = struct{}{}
		items = append(items, symbol)
	}

	return items
}
