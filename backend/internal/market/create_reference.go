package market

import (
	"context"
	"fmt"
	"strings"

	"prediction/internal/domain"
	"prediction/internal/pacifica"
)

func (s *service) enrichCreateInput(ctx context.Context, input CreateInput) (CreateInput, error) {
	if input.MarketType != domain.MarketTypePriceThreshold {
		return input, nil
	}

	if strings.TrimSpace(input.Symbol) == "" {
		return input, nil
	}

	if s.createContextProvider == nil {
		return CreateInput{}, fmt.Errorf("create context provider is not configured")
	}

	marketInfoItems, err := s.createContextProvider.ListMarketInfo(ctx)
	if err != nil {
		return CreateInput{}, fmt.Errorf("list market info for price-threshold creation: %w", err)
	}

	marketInfo, ok := findMarketInfoBySymbol(marketInfoItems, input.Symbol)
	if !ok {
		return CreateInput{}, domain.NewValidationError("symbol", "symbol is not supported by Pacifica", input.Symbol)
	}

	priceItems, err := s.createContextProvider.ListPrices(ctx, pacifica.PriceFilter{
		Symbols: []string{input.Symbol},
	})
	if err != nil {
		return CreateInput{}, fmt.Errorf("list live prices for price-threshold creation: %w", err)
	}

	priceSnapshot, ok := findPriceSnapshotBySymbol(priceItems, input.Symbol)
	if !ok || strings.TrimSpace(priceSnapshot.MarkPrice) == "" {
		return CreateInput{}, domain.NewValidationError("symbol", "live mark price is not available for this symbol", input.Symbol)
	}

	input.SymbolPriceIncrement = strings.TrimSpace(marketInfo.MinTick)
	input.ReferenceValue = strings.TrimSpace(priceSnapshot.MarkPrice)

	return input, nil
}

func findMarketInfoBySymbol(items []pacifica.MarketInfo, symbol string) (pacifica.MarketInfo, bool) {
	for _, item := range items {
		if strings.EqualFold(item.Symbol, symbol) {
			return item, true
		}
	}

	return pacifica.MarketInfo{}, false
}

func findPriceSnapshotBySymbol(items []pacifica.PriceSnapshot, symbol string) (pacifica.PriceSnapshot, bool) {
	for _, item := range items {
		if strings.EqualFold(item.Symbol, symbol) {
			return item, true
		}
	}

	return pacifica.PriceSnapshot{}, false
}
