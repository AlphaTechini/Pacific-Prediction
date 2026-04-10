package market

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"prediction/internal/domain"
	"prediction/internal/pacifica"
	"prediction/internal/realtime"
	"prediction/internal/storage"

	"github.com/jackc/pgx/v5"
)

type service struct {
	marketRepository          storage.MarketRepository
	createContextProvider     CreateContextProvider
	validator                 Validator
	publisher                 realtime.Publisher
	txManager                 *storage.TxManager
	marketRepositoryFactory   func(storage.Queryer) storage.MarketRepository
	balanceRepositoryFactory  func(storage.Queryer) storage.BalanceRepository
	positionRepositoryFactory func(storage.Queryer) storage.PositionRepository
}

type ServiceDeps struct {
	MarketRepository          storage.MarketRepository
	CreateContextProvider     CreateContextProvider
	Validator                 Validator
	Publisher                 realtime.Publisher
	TxManager                 *storage.TxManager
	MarketRepositoryFactory   func(storage.Queryer) storage.MarketRepository
	BalanceRepositoryFactory  func(storage.Queryer) storage.BalanceRepository
	PositionRepositoryFactory func(storage.Queryer) storage.PositionRepository
}

type CreateContextProvider interface {
	ListMarketInfo(ctx context.Context) ([]pacifica.MarketInfo, error)
	ListPrices(ctx context.Context, filter pacifica.PriceFilter) ([]pacifica.PriceSnapshot, error)
}

func NewService(marketRepository storage.MarketRepository, createContextProvider CreateContextProvider, validator Validator) Service {
	return NewServiceWithDeps(ServiceDeps{
		MarketRepository:      marketRepository,
		CreateContextProvider: createContextProvider,
		Validator:             validator,
	})
}

func NewServiceWithDeps(deps ServiceDeps) Service {
	marketRepositoryFactory := deps.MarketRepositoryFactory
	if marketRepositoryFactory == nil {
		marketRepositoryFactory = func(queryer storage.Queryer) storage.MarketRepository {
			return storage.NewMarketPostgresRepository(queryer)
		}
	}

	balanceRepositoryFactory := deps.BalanceRepositoryFactory
	if balanceRepositoryFactory == nil {
		balanceRepositoryFactory = func(queryer storage.Queryer) storage.BalanceRepository {
			return storage.NewBalancePostgresRepository(queryer)
		}
	}

	positionRepositoryFactory := deps.PositionRepositoryFactory
	if positionRepositoryFactory == nil {
		positionRepositoryFactory = func(queryer storage.Queryer) storage.PositionRepository {
			return storage.NewPositionPostgresRepository(queryer)
		}
	}

	return &service{
		marketRepository:          deps.MarketRepository,
		createContextProvider:     deps.CreateContextProvider,
		validator:                 deps.Validator,
		publisher:                 deps.Publisher,
		txManager:                 deps.TxManager,
		marketRepositoryFactory:   marketRepositoryFactory,
		balanceRepositoryFactory:  balanceRepositoryFactory,
		positionRepositoryFactory: positionRepositoryFactory,
	}
}

func (s *service) Create(ctx context.Context, input CreateInput) (Record, error) {
	now := domain.NowUTC()
	normalized := normalizeCreateInput(input, now)
	enriched, err := s.enrichCreateInput(ctx, normalized)
	if err != nil {
		return Record{}, err
	}

	if err := s.validator.ValidateCreateInput(ctx, enriched); err != nil {
		return Record{}, err
	}

	marketID, err := NewMarketID()
	if err != nil {
		return Record{}, err
	}

	creatorPositionID, err := newCreatorPositionID()
	if err != nil {
		return Record{}, err
	}

	potentialPayout, err := domain.CalculateEvenOddsPayout(normalized.CreatorStakeAmount)
	if err != nil {
		return Record{}, domain.NewValidationError("creator_stake_amount", "creator stake amount must be a valid decimal value", normalized.CreatorStakeAmount)
	}

	createMarketInput := storage.CreateMarketInput{
		ID:                marketID,
		Title:             enriched.Title,
		Symbol:            enriched.Symbol,
		MarketType:        enriched.MarketType,
		ConditionOperator: enriched.ConditionOperator,
		ThresholdValue:    enriched.ThresholdValue,
		SourceType:        enriched.SourceType,
		SourceInterval:    enriched.SourceInterval,
		ReferenceValue:    enriched.ReferenceValue,
		ExpiryTime:        enriched.ExpiryTime,
		CreatedByPlayerID: enriched.CreatedByPlayerID,
	}

	var created storage.Market
	if s.txManager == nil {
		created, err = s.marketRepository.Create(ctx, createMarketInput)
		if err != nil {
			return Record{}, fmt.Errorf("create market: %w", err)
		}

		return toRecord(created), nil
	}

	if err := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		marketRepository := s.marketRepositoryFactory(tx)
		balanceRepository := s.balanceRepositoryFactory(tx)
		positionRepository := s.positionRepositoryFactory(tx)

		var createErr error
		created, createErr = marketRepository.Create(ctx, createMarketInput)
		if createErr != nil {
			return fmt.Errorf("create market: %w", createErr)
		}

		if _, createErr = balanceRepository.LockStake(ctx, storage.LockStakeInput{
			PlayerID: enriched.CreatedByPlayerID,
			Amount:   enriched.CreatorStakeAmount,
		}); createErr != nil {
			return fmt.Errorf("lock creator stake: %w", createErr)
		}

		if _, createErr = positionRepository.Create(ctx, storage.CreatePositionInput{
			ID:              creatorPositionID,
			PlayerID:        enriched.CreatedByPlayerID,
			MarketID:        marketID,
			Side:            enriched.CreatorSide,
			StakeAmount:     enriched.CreatorStakeAmount,
			PotentialPayout: potentialPayout,
		}); createErr != nil {
			return fmt.Errorf("create creator opening position: %w", createErr)
		}

		return nil
	}); err != nil {
		return Record{}, err
	}

	s.publishMarketCreated(ctx, created)

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
		Symbols:                           symbols,
		ValidationModels:                  SupportedValidationModels(),
		PriceThresholdCreationBandPercent: s.validatorConfig().PriceThresholdCreationBandPercent,
	}, nil
}

func (s *service) ValidateCreateInput(ctx context.Context, input CreateInput) error {
	normalized := normalizeCreateInput(input, domain.NowUTC())
	enriched, err := s.enrichCreateInput(ctx, normalized)
	if err != nil {
		return err
	}

	return s.validator.ValidateCreateInput(ctx, enriched)
}

func (s *service) SupportedValidationModels() []ValidationModel {
	return SupportedValidationModels()
}

func (s *service) validatorConfig() ValidationConfig {
	validationService, ok := s.validator.(*validationService)
	if !ok {
		return ValidationConfig{}
	}

	return validationService.config
}

func normalizeCreateInput(input CreateInput, now time.Time) CreateInput {
	input.Title = strings.TrimSpace(input.Title)
	input.Symbol = strings.ToUpper(strings.TrimSpace(input.Symbol))
	input.SymbolTickSize = strings.TrimSpace(input.SymbolTickSize)
	input.CreatorStakeAmount = strings.TrimSpace(input.CreatorStakeAmount)
	input.ThresholdValue = strings.TrimSpace(input.ThresholdValue)
	input.SourceInterval = strings.TrimSpace(input.SourceInterval)
	input.ReferenceValue = strings.TrimSpace(input.ReferenceValue)
	input.ExpiryTime = deriveCreateExpiryTime(input, now)
	return input
}

func deriveCreateExpiryTime(input CreateInput, now time.Time) time.Time {
	switch input.MarketType {
	case domain.MarketTypeCandleDirection:
		expiryTime, err := domain.NextCandleExpiryFromCreation(now, input.SourceInterval)
		if err == nil {
			return expiryTime
		}
	case domain.MarketTypeFundingThreshold:
		return domain.NextFundingEpochFromCreation(now)
	}

	return domain.NormalizeTime(input.ExpiryTime)
}

func newCreatorPositionID() (domain.PositionID, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate creator position id: %w", err)
	}

	return domain.PositionID("position_" + hex.EncodeToString(buf)), nil
}

func (s *service) publishMarketCreated(ctx context.Context, item storage.Market) {
	if s.publisher == nil {
		return
	}

	_ = s.publisher.Publish(context.WithoutCancel(ctx), realtime.NewMarketCreatedEvent(item))
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
