package market

import (
	"context"
	"fmt"
	"strings"

	"prediction/internal/domain"
	"prediction/internal/pacifica"
	"prediction/internal/storage"
)

type service struct {
	marketRepository      storage.MarketRepository
	createContextProvider CreateContextProvider
	validator             Validator
}

type CreateContextProvider interface {
	ListMarketInfo(ctx context.Context) ([]pacifica.MarketInfo, error)
	ListPrices(ctx context.Context, filter pacifica.PriceFilter) ([]pacifica.PriceSnapshot, error)
}

func NewService(marketRepository storage.MarketRepository, createContextProvider CreateContextProvider, validator Validator) Service {
	return &service{
		marketRepository:      marketRepository,
		createContextProvider: createContextProvider,
		validator:             validator,
	}
}

func (s *service) Create(ctx context.Context, input CreateInput) (Record, error) {
	normalized := normalizeCreateInput(input)
	if err := s.validator.ValidateCreateInput(ctx, normalized); err != nil {
		return Record{}, err
	}

	marketID, err := NewMarketID()
	if err != nil {
		return Record{}, err
	}

	created, err := s.marketRepository.Create(ctx, storage.CreateMarketInput{
		ID:                marketID,
		Title:             normalized.Title,
		Symbol:            normalized.Symbol,
		MarketType:        normalized.MarketType,
		ConditionOperator: normalized.ConditionOperator,
		ThresholdValue:    normalized.ThresholdValue,
		SourceType:        normalized.SourceType,
		SourceInterval:    normalized.SourceInterval,
		ReferenceValue:    normalized.ReferenceValue,
		ExpiryTime:        normalized.ExpiryTime,
		CreatedByPlayerID: normalized.CreatedByPlayerID,
	})
	if err != nil {
		return Record{}, fmt.Errorf("create market: %w", err)
	}

	return toRecord(created), nil
}

func (s *service) List(ctx context.Context, filter ListFilter) ([]Record, error) {
	if filter.Status == "" {
		filter.Status = domain.MarketStatusActive
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	items, err := s.marketRepository.ListByStatus(ctx, filter.Status, filter.Limit)
	if err != nil {
		return nil, fmt.Errorf("list markets: %w", err)
	}

	records := make([]Record, 0, len(items))
	for _, item := range items {
		records = append(records, toRecord(item))
	}

	return records, nil
}

func (s *service) ListCatalog(ctx context.Context, limitPerStatus int) (Catalog, error) {
	if limitPerStatus <= 0 {
		limitPerStatus = 50
	}

	active, err := s.List(ctx, ListFilter{
		Status: domain.MarketStatusActive,
		Limit:  limitPerStatus,
	})
	if err != nil {
		return Catalog{}, fmt.Errorf("list active markets: %w", err)
	}

	resolved, err := s.List(ctx, ListFilter{
		Status: domain.MarketStatusResolved,
		Limit:  limitPerStatus,
	})
	if err != nil {
		return Catalog{}, fmt.Errorf("list resolved markets: %w", err)
	}

	return Catalog{
		Active:   active,
		Resolved: resolved,
	}, nil
}

func (s *service) GetByID(ctx context.Context, marketID domain.MarketID) (Record, error) {
	item, err := s.marketRepository.GetByID(ctx, marketID)
	if err != nil {
		return Record{}, fmt.Errorf("get market by id: %w", err)
	}

	return toRecord(item), nil
}

func (s *service) GetCreateContext(ctx context.Context) (CreateContext, error) {
	marketInfoItems, err := s.createContextProvider.ListMarketInfo(ctx)
	if err != nil {
		return CreateContext{}, fmt.Errorf("list market info for create context: %w", err)
	}

	priceItems, err := s.createContextProvider.ListPrices(ctx, pacifica.PriceFilter{})
	if err != nil {
		return CreateContext{}, fmt.Errorf("list prices for create context: %w", err)
	}

	priceBySymbol := make(map[string]pacifica.PriceSnapshot, len(priceItems))
	for _, item := range priceItems {
		priceBySymbol[item.Symbol] = item
	}

	symbols := make([]CreateContextSymbol, 0, len(marketInfoItems))
	for _, item := range marketInfoItems {
		price := priceBySymbol[item.Symbol]
		symbols = append(symbols, CreateContextSymbol{
			Symbol:          item.Symbol,
			TickSize:        item.TickSize,
			MinTick:         item.MinTick,
			MaxTick:         item.MaxTick,
			LotSize:         item.LotSize,
			MinOrderSize:    item.MinOrderSize,
			MaxOrderSize:    item.MaxOrderSize,
			MaxLeverage:     item.MaxLeverage,
			IsolatedOnly:    item.IsolatedOnly,
			MarkPrice:       price.MarkPrice,
			OraclePrice:     price.OraclePrice,
			FundingRate:     price.FundingRate,
			NextFundingRate: price.NextFundingRate,
			OpenInterest:    price.OpenInterest,
			Volume24H:       price.Volume24H,
			UpdatedAt:       price.Timestamp,
		})
	}

	return CreateContext{
		Symbols:          symbols,
		ValidationModels: SupportedValidationModels(),
	}, nil
}

func (s *service) ValidateCreateInput(ctx context.Context, input CreateInput) error {
	return s.validator.ValidateCreateInput(ctx, normalizeCreateInput(input))
}

func (s *service) SupportedValidationModels() []ValidationModel {
	return SupportedValidationModels()
}

func normalizeCreateInput(input CreateInput) CreateInput {
	input.Title = strings.TrimSpace(input.Title)
	input.Symbol = strings.ToUpper(strings.TrimSpace(input.Symbol))
	input.ThresholdValue = strings.TrimSpace(input.ThresholdValue)
	input.SourceInterval = strings.TrimSpace(input.SourceInterval)
	input.ReferenceValue = strings.TrimSpace(input.ReferenceValue)
	input.ExpiryTime = domain.NormalizeTime(input.ExpiryTime)
	return input
}

func toRecord(item storage.Market) Record {
	return Record{
		ID:                item.ID,
		Title:             item.Title,
		Symbol:            item.Symbol,
		MarketType:        item.MarketType,
		ConditionOperator: item.ConditionOperator,
		ThresholdValue:    item.ThresholdValue,
		SourceType:        item.SourceType,
		SourceInterval:    item.SourceInterval,
		ReferenceValue:    item.ReferenceValue,
		ExpiryTime:        item.ExpiryTime,
		Status:            item.Status,
		Result:            item.Result,
		SettlementValue:   item.SettlementValue,
		ResolvedAt:        item.ResolvedAt,
		ResolutionReason:  item.ResolutionReason,
		CreatedByPlayerID: item.CreatedByPlayerID,
		CreatedAt:         item.CreatedAt,
	}
}
