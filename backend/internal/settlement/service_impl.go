package settlement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"prediction/internal/domain"
	"prediction/internal/pacifica"
	"prediction/internal/storage"

	"github.com/jackc/pgx/v5"
)

type service struct {
	marketRepository            storage.MarketRepository
	priceResolver               PriceResolver
	candleResolver              CandleResolver
	fundingResolver             FundingResolver
	txManager                   settlementTxManager
	marketRepositoryFactory     func(storage.Queryer) storage.MarketRepository
	settlementRepositoryFactory func(storage.Queryer) storage.SettlementRepository
	priceRetryInterval          time.Duration
	sleep                       func(context.Context, time.Duration) error
}

type ServiceDeps struct {
	MarketRepository            storage.MarketRepository
	PacificaClient              pacifica.RESTClient
	PriceResolver               PriceResolver
	CandleResolver              CandleResolver
	FundingResolver             FundingResolver
	TxManager                   settlementTxManager
	MarketRepositoryFactory     func(storage.Queryer) storage.MarketRepository
	SettlementRepositoryFactory func(storage.Queryer) storage.SettlementRepository
	PriceRetryInterval          time.Duration
}

type settlementTxManager interface {
	WithinTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error
}

func NewService(deps ServiceDeps) Service {
	marketRepositoryFactory := deps.MarketRepositoryFactory
	if marketRepositoryFactory == nil {
		marketRepositoryFactory = func(queryer storage.Queryer) storage.MarketRepository {
			return storage.NewMarketPostgresRepository(queryer)
		}
	}

	settlementRepositoryFactory := deps.SettlementRepositoryFactory
	if settlementRepositoryFactory == nil {
		settlementRepositoryFactory = func(queryer storage.Queryer) storage.SettlementRepository {
			return storage.NewSettlementPostgresRepository(queryer)
		}
	}

	priceResolver := deps.PriceResolver
	if priceResolver == nil && deps.PacificaClient != nil {
		priceResolver = NewPriceResolver(deps.PacificaClient)
	}

	candleResolver := deps.CandleResolver
	if candleResolver == nil && deps.PacificaClient != nil {
		candleResolver = NewCandleResolver(deps.PacificaClient)
	}

	fundingResolver := deps.FundingResolver
	if fundingResolver == nil && deps.PacificaClient != nil {
		fundingResolver = NewFundingResolver(deps.PacificaClient)
	}

	return &service{
		marketRepository:            deps.MarketRepository,
		priceResolver:               priceResolver,
		candleResolver:              candleResolver,
		fundingResolver:             fundingResolver,
		txManager:                   deps.TxManager,
		marketRepositoryFactory:     marketRepositoryFactory,
		settlementRepositoryFactory: settlementRepositoryFactory,
		priceRetryInterval:          deps.PriceRetryInterval,
		sleep: func(ctx context.Context, duration time.Duration) error {
			timer := time.NewTimer(duration)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				return nil
			}
		},
	}
}

func (s *service) SettleMarket(ctx context.Context, marketID domain.MarketID) (Attempt, error) {
	item, err := s.marketRepository.GetByID(ctx, marketID)
	if err != nil {
		return Attempt{}, fmt.Errorf("get market for settlement scan: %w", err)
	}

	if item.Status != domain.MarketStatusActive {
		return attemptFromMarket(item), nil
	}

	if item.ExpiryTime.After(domain.NowUTC()) {
		return Attempt{}, domain.NewValidationError("market_id", "market has not expired yet", marketID)
	}

	switch item.MarketType {
	case domain.MarketTypePriceThreshold:
		return s.settlePriceMarket(ctx, item)
	case domain.MarketTypeCandleDirection:
		return s.settleCandleMarket(ctx, item)
	case domain.MarketTypeFundingThreshold:
		return s.settleFundingMarket(ctx, item)
	default:
		return attemptFromMarket(item), nil
	}
}

func (s *service) SettleDueMarkets(ctx context.Context, filter DueMarketFilter) ([]Attempt, error) {
	normalized := normalizeDueMarketFilter(filter)

	items, err := s.marketRepository.ListExpiringBefore(ctx, normalized.Before, normalized.Limit)
	if err != nil {
		return nil, fmt.Errorf("list due markets for settlement scan: %w", err)
	}

	return s.settleDueMarketItems(ctx, items)
}

func (s *service) PlanPriceFetchBatches(ctx context.Context, filter PriceFetchPlanFilter) ([]PriceFetchBatch, error) {
	normalized := normalizePriceFetchPlanFilter(filter)

	items, err := s.marketRepository.ListExpiringBefore(ctx, normalized.Before, normalized.Limit)
	if err != nil {
		return nil, fmt.Errorf("list upcoming price markets for fetch planning: %w", err)
	}

	return buildPriceFetchBatches(items, normalized.After), nil
}

func normalizeDueMarketFilter(filter DueMarketFilter) DueMarketFilter {
	if filter.Before.IsZero() {
		filter.Before = domain.NowUTC()
	} else {
		filter.Before = domain.NormalizeTime(filter.Before)
	}

	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	return filter
}

func attemptFromMarket(item storage.Market) Attempt {
	return Attempt{
		MarketID:   item.ID,
		MarketType: item.MarketType,
		Settled:    false,
	}
}

func (s *service) settleDueMarketItems(ctx context.Context, items []storage.Market) ([]Attempt, error) {
	marketsByID := make(map[domain.MarketID]storage.Market, len(items))
	attemptByID := make(map[domain.MarketID]Attempt, len(items))
	for _, item := range items {
		marketsByID[item.ID] = item
		attemptByID[item.ID] = attemptFromMarket(item)
	}

	for _, batch := range buildDuePriceFetchBatches(items) {
		batchAttempts, err := s.settlePriceBatch(ctx, batch, marketsByID)
		if err != nil {
			return nil, err
		}

		for _, attempt := range batchAttempts {
			attemptByID[attempt.MarketID] = attempt
		}
	}

	for _, item := range items {
		switch item.MarketType {
		case domain.MarketTypeCandleDirection:
			attempt, err := s.settleCandleMarket(ctx, item)
			if err != nil {
				return nil, err
			}

			attemptByID[item.ID] = attempt
		case domain.MarketTypeFundingThreshold:
			attempt, err := s.settleFundingMarket(ctx, item)
			if err != nil {
				return nil, err
			}

			attemptByID[item.ID] = attempt
		}
	}

	attempts := make([]Attempt, 0, len(items))
	for _, item := range items {
		attempts = append(attempts, attemptByID[item.ID])
	}

	return attempts, nil
}

func (s *service) settlePriceMarket(ctx context.Context, item storage.Market) (Attempt, error) {
	if item.MarketType != domain.MarketTypePriceThreshold {
		return attemptFromMarket(item), nil
	}

	resolution, err := s.resolveDuePriceMarket(ctx, item)
	if err != nil {
		if errors.Is(err, errSettlementSourceNotReady) {
			return attemptFromMarket(item), nil
		}

		return Attempt{}, err
	}

	return s.persistPriceResolution(ctx, item, resolution)
}

func (s *service) settleCandleMarket(ctx context.Context, item storage.Market) (Attempt, error) {
	if item.MarketType != domain.MarketTypeCandleDirection {
		return attemptFromMarket(item), nil
	}

	resolution, err := s.resolveDueCandleMarket(ctx, item)
	if err != nil {
		if errors.Is(err, errSettlementSourceNotReady) {
			return attemptFromMarket(item), nil
		}

		return Attempt{}, err
	}

	return s.persistCandleResolution(ctx, item, resolution)
}

func (s *service) settleFundingMarket(ctx context.Context, item storage.Market) (Attempt, error) {
	if item.MarketType != domain.MarketTypeFundingThreshold {
		return attemptFromMarket(item), nil
	}

	resolution, err := s.resolveDueFundingMarket(ctx, item)
	if err != nil {
		if errors.Is(err, errSettlementSourceNotReady) {
			return attemptFromMarket(item), nil
		}

		return Attempt{}, err
	}

	return s.persistFundingResolution(ctx, item, resolution)
}

func (s *service) settlePriceBatch(ctx context.Context, batch PriceFetchBatch, marketsByID map[domain.MarketID]storage.Market) ([]Attempt, error) {
	if s.priceResolver == nil {
		return nil, fmt.Errorf("price resolver is not configured")
	}

	priceMarkets := make([]PriceMarket, 0, len(batch.Targets))
	for _, target := range batch.Targets {
		item, ok := marketsByID[target.MarketID]
		if !ok {
			continue
		}

		priceMarkets = append(priceMarkets, PriceMarket{
			ID:                item.ID,
			Symbol:            item.Symbol,
			ConditionOperator: item.ConditionOperator,
			ThresholdValue:    item.ThresholdValue,
			ExpiryTime:        item.ExpiryTime,
		})
	}

	resolutions, err := s.resolveDuePriceBatch(ctx, priceMarkets)
	if err != nil {
		if errors.Is(err, errSettlementSourceNotReady) {
			attempts := make([]Attempt, 0, len(priceMarkets))
			for _, market := range priceMarkets {
				attempts = append(attempts, Attempt{
					MarketID:   market.ID,
					MarketType: domain.MarketTypePriceThreshold,
					Settled:    false,
				})
			}

			return attempts, nil
		}

		return nil, fmt.Errorf("resolve price batch for expiry %s: %w", batch.ExpiryTime.Format(time.RFC3339), err)
	}

	resolutionByID := make(map[domain.MarketID]PriceResolution, len(resolutions))
	for _, resolution := range resolutions {
		resolutionByID[resolution.MarketID] = resolution
	}

	attempts := make([]Attempt, 0, len(priceMarkets))
	for _, market := range priceMarkets {
		item, ok := marketsByID[market.ID]
		if !ok {
			continue
		}

		resolution, ok := resolutionByID[market.ID]
		if !ok {
			attempts = append(attempts, attemptFromMarket(item))
			continue
		}

		attempt, err := s.persistPriceResolution(ctx, item, resolution)
		if err != nil {
			return nil, err
		}

		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

func (s *service) resolveDuePriceMarket(ctx context.Context, item storage.Market) (PriceResolution, error) {
	if s.priceResolver == nil {
		return PriceResolution{}, fmt.Errorf("price resolver is not configured")
	}

	market := PriceMarket{
		ID:                item.ID,
		Symbol:            item.Symbol,
		ConditionOperator: item.ConditionOperator,
		ThresholdValue:    item.ThresholdValue,
		ExpiryTime:        item.ExpiryTime,
	}

	resolution, err := s.priceResolver.Resolve(ctx, market)
	if err == nil {
		return resolution, nil
	}

	if s.priceRetryInterval <= 0 || !errors.Is(err, errSettlementSourceNotReady) {
		return PriceResolution{}, err
	}

	if sleepErr := s.sleep(ctx, s.priceRetryInterval); sleepErr != nil {
		return PriceResolution{}, err
	}

	return s.priceResolver.Resolve(ctx, market)
}

func (s *service) resolveDueCandleMarket(ctx context.Context, item storage.Market) (CandleResolution, error) {
	if s.candleResolver == nil {
		return CandleResolution{}, fmt.Errorf("candle resolver is not configured")
	}

	return s.candleResolver.Resolve(ctx, CandleMarket{
		ID:                item.ID,
		Symbol:            item.Symbol,
		ConditionOperator: item.ConditionOperator,
		SourceInterval:    item.SourceInterval,
		ExpiryTime:        item.ExpiryTime,
	})
}

func (s *service) resolveDueFundingMarket(ctx context.Context, item storage.Market) (FundingResolution, error) {
	if s.fundingResolver == nil {
		return FundingResolution{}, fmt.Errorf("funding resolver is not configured")
	}

	return s.fundingResolver.Resolve(ctx, FundingMarket{
		ID:                item.ID,
		Symbol:            item.Symbol,
		ConditionOperator: item.ConditionOperator,
		ThresholdValue:    item.ThresholdValue,
		ExpiryTime:        item.ExpiryTime,
	})
}

func (s *service) resolveDuePriceBatch(ctx context.Context, markets []PriceMarket) ([]PriceResolution, error) {
	resolutions, err := s.priceResolver.ResolveBatch(ctx, markets)
	if err == nil {
		return resolutions, nil
	}

	if s.priceRetryInterval <= 0 || !errors.Is(err, errSettlementSourceNotReady) {
		return nil, err
	}

	if sleepErr := s.sleep(ctx, s.priceRetryInterval); sleepErr != nil {
		return nil, err
	}

	return s.priceResolver.ResolveBatch(ctx, markets)
}

func (s *service) persistCandleResolution(ctx context.Context, item storage.Market, resolution CandleResolution) (Attempt, error) {
	if s.txManager == nil {
		return Attempt{}, fmt.Errorf("settlement transaction manager is not configured")
	}

	settlementID, err := NewSettlementID()
	if err != nil {
		return Attempt{}, err
	}

	resolvedAt := resolution.CloseTime
	if err := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		marketRepository := s.marketRepositoryFactory(tx)
		settlementRepository := s.settlementRepositoryFactory(tx)

		if _, err := settlementRepository.Create(ctx, storage.CreateSettlementInput{
			ID:              settlementID,
			MarketID:        item.ID,
			PacificaSource:  resolution.PacificaSource,
			SourceTimestamp: resolution.CloseTime,
			RawPayload:      resolution.RawPayload,
			SettlementValue: resolution.ClosePrice,
			Result:          resolution.Result,
		}); err != nil {
			return fmt.Errorf("create candle settlement audit: %w", err)
		}

		if _, err := marketRepository.UpdateSettlement(ctx, storage.UpdateMarketSettlementInput{
			MarketID:         item.ID,
			Status:           domain.MarketStatusResolved,
			Result:           resolution.Result,
			SettlementValue:  resolution.ClosePrice,
			ResolvedAt:       resolvedAt,
			ResolutionReason: "candle_direction_mark_price_close",
		}); err != nil {
			return fmt.Errorf("update market settlement state: %w", err)
		}

		return nil
	}); err != nil {
		return Attempt{}, fmt.Errorf("persist candle settlement: %w", err)
	}

	return Attempt{
		MarketID:      item.ID,
		MarketType:    item.MarketType,
		Settled:       true,
		SettlementID:  &settlementID,
		SettledAt:     &resolvedAt,
		SettlementRef: resolution.PacificaSource,
	}, nil
}

func (s *service) persistFundingResolution(ctx context.Context, item storage.Market, resolution FundingResolution) (Attempt, error) {
	if s.txManager == nil {
		return Attempt{}, fmt.Errorf("settlement transaction manager is not configured")
	}

	settlementID, err := NewSettlementID()
	if err != nil {
		return Attempt{}, err
	}

	resolvedAt := resolution.SettlementTime
	if err := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		marketRepository := s.marketRepositoryFactory(tx)
		settlementRepository := s.settlementRepositoryFactory(tx)

		if _, err := settlementRepository.Create(ctx, storage.CreateSettlementInput{
			ID:              settlementID,
			MarketID:        item.ID,
			PacificaSource:  resolution.PacificaSource,
			SourceTimestamp: resolution.SettlementTime,
			RawPayload:      resolution.RawPayload,
			SettlementValue: resolution.FundingRate,
			Result:          resolution.Result,
		}); err != nil {
			return fmt.Errorf("create funding settlement audit: %w", err)
		}

		if _, err := marketRepository.UpdateSettlement(ctx, storage.UpdateMarketSettlementInput{
			MarketID:         item.ID,
			Status:           domain.MarketStatusResolved,
			Result:           resolution.Result,
			SettlementValue:  resolution.FundingRate,
			ResolvedAt:       resolvedAt,
			ResolutionReason: "funding_threshold_history",
		}); err != nil {
			return fmt.Errorf("update market settlement state: %w", err)
		}

		return nil
	}); err != nil {
		return Attempt{}, fmt.Errorf("persist funding settlement: %w", err)
	}

	return Attempt{
		MarketID:      item.ID,
		MarketType:    item.MarketType,
		Settled:       true,
		SettlementID:  &settlementID,
		SettledAt:     &resolvedAt,
		SettlementRef: resolution.PacificaSource,
	}, nil
}

func (s *service) persistPriceResolution(ctx context.Context, item storage.Market, resolution PriceResolution) (Attempt, error) {
	if s.txManager == nil {
		return Attempt{}, fmt.Errorf("settlement transaction manager is not configured")
	}

	settlementID, err := NewSettlementID()
	if err != nil {
		return Attempt{}, err
	}

	resolvedAt := resolution.SourceTimestamp
	if err := s.txManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		marketRepository := s.marketRepositoryFactory(tx)
		settlementRepository := s.settlementRepositoryFactory(tx)

		if _, err := settlementRepository.Create(ctx, storage.CreateSettlementInput{
			ID:              settlementID,
			MarketID:        item.ID,
			PacificaSource:  resolution.PacificaSource,
			SourceTimestamp: resolution.SourceTimestamp,
			RawPayload:      resolution.RawPayload,
			SettlementValue: resolution.SettlementMarkPrice,
			Result:          resolution.Result,
		}); err != nil {
			return fmt.Errorf("create price settlement audit: %w", err)
		}

		if _, err := marketRepository.UpdateSettlement(ctx, storage.UpdateMarketSettlementInput{
			MarketID:         item.ID,
			Status:           domain.MarketStatusResolved,
			Result:           resolution.Result,
			SettlementValue:  resolution.SettlementMarkPrice,
			ResolvedAt:       resolvedAt,
			ResolutionReason: "price_threshold_mark_price",
		}); err != nil {
			return fmt.Errorf("update market settlement state: %w", err)
		}

		return nil
	}); err != nil {
		return Attempt{}, fmt.Errorf("persist price settlement: %w", err)
	}

	return Attempt{
		MarketID:      item.ID,
		MarketType:    item.MarketType,
		Settled:       true,
		SettlementID:  &settlementID,
		SettledAt:     &resolvedAt,
		SettlementRef: resolution.PacificaSource,
	}, nil
}
