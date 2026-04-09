package settlement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"prediction/internal/domain"
	"prediction/internal/pacifica"
)

const pacificaMarkPriceCandlesSource = "/api/v1/kline/mark"

type candleResolver struct {
	pacificaClient pacifica.RESTClient
}

func NewCandleResolver(pacificaClient pacifica.RESTClient) CandleResolver {
	return &candleResolver{
		pacificaClient: pacificaClient,
	}
}

func (r *candleResolver) Resolve(ctx context.Context, market CandleMarket) (CandleResolution, error) {
	if strings.TrimSpace(market.Symbol) == "" {
		return CandleResolution{}, domain.NewValidationError("symbol", "symbol is required for candle settlement", market.Symbol)
	}

	candleOpenTime, candleCloseTime, err := domain.CandleWindowForExpiry(market.ExpiryTime, market.SourceInterval)
	if err != nil {
		return CandleResolution{}, domain.NewValidationError("expiry_time", "candle market expiry is not aligned to its interval", market.ExpiryTime)
	}

	items, err := r.pacificaClient.ListMarkPriceCandles(ctx, pacifica.MarkPriceCandleQuery{
		Symbol:    market.Symbol,
		Interval:  market.SourceInterval,
		StartTime: candleOpenTime,
		EndTime:   candleCloseTime,
		Limit:     3,
	})
	if err != nil {
		wrapped := fmt.Errorf("fetch mark price candle for market settlement: %w", err)
		if errors.Is(err, pacifica.ErrTemporaryFailure) {
			return CandleResolution{}, markTemporarySettlementError(wrapped)
		}

		return CandleResolution{}, wrapped
	}

	candle, err := findSettlementCandle(items, market.Symbol, market.SourceInterval, candleOpenTime, candleCloseTime)
	if err != nil {
		return CandleResolution{}, err
	}

	result, err := compareCandleDirection(candle.OpenPrice, candle.ClosePrice, market.ConditionOperator)
	if err != nil {
		return CandleResolution{}, err
	}

	rawPayload, err := json.Marshal(candle)
	if err != nil {
		return CandleResolution{}, fmt.Errorf("marshal mark price candle payload: %w", err)
	}

	return CandleResolution{
		MarketID:       market.ID,
		PacificaSource: pacificaMarkPriceCandlesSource,
		OpenTime:       candle.OpenTime,
		CloseTime:      candle.CloseTime,
		OpenPrice:      candle.OpenPrice,
		ClosePrice:     candle.ClosePrice,
		Result:         result,
		RawPayload:     rawPayload,
	}, nil
}

func findSettlementCandle(
	items []pacifica.MarkPriceCandle,
	symbol string,
	interval string,
	expectedOpenTime time.Time,
	expectedCloseTime time.Time,
) (pacifica.MarkPriceCandle, error) {
	for _, item := range items {
		if !strings.EqualFold(item.Symbol, symbol) {
			continue
		}
		if strings.TrimSpace(item.Interval) != strings.TrimSpace(interval) {
			continue
		}
		if !domain.NormalizeTime(item.OpenTime).Equal(expectedOpenTime) {
			continue
		}
		if !domain.NormalizeTime(item.CloseTime).Equal(expectedCloseTime) {
			continue
		}

		return item, nil
	}

	return pacifica.MarkPriceCandle{}, fmt.Errorf(
		"%w: mark-price candle is not yet available for symbol=%s interval=%s open=%s close=%s",
		errSettlementSourceNotReady,
		symbol,
		interval,
		expectedOpenTime.Format(time.RFC3339Nano),
		expectedCloseTime.Format(time.RFC3339Nano),
	)
}

func compareCandleDirection(openPrice string, closePrice string, operator domain.ConditionOperator) (domain.MarketResult, error) {
	openValue, err := parseSettlementDecimal(openPrice)
	if err != nil {
		return "", domain.NewValidationError("open_price", "candle open price must be a valid decimal string", openPrice)
	}

	closeValue, err := parseSettlementDecimal(closePrice)
	if err != nil {
		return "", domain.NewValidationError("close_price", "candle close price must be a valid decimal string", closePrice)
	}

	comparison := closeValue.Cmp(openValue)

	switch operator {
	case domain.ConditionOperatorBullishClose:
		if comparison > 0 {
			return domain.MarketResultYes, nil
		}
	case domain.ConditionOperatorBearishClose:
		if comparison < 0 {
			return domain.MarketResultYes, nil
		}
	default:
		return "", domain.NewValidationError("condition_operator", "candle settlement operator is not supported", operator)
	}

	return domain.MarketResultNo, nil
}
