package settlement

import (
	"context"
	"fmt"

	"prediction/internal/domain"
	"prediction/internal/storage"
)

type service struct {
	marketRepository storage.MarketRepository
}

type ServiceDeps struct {
	MarketRepository storage.MarketRepository
}

func NewService(deps ServiceDeps) Service {
	return &service{
		marketRepository: deps.MarketRepository,
	}
}

func (s *service) SettleMarket(ctx context.Context, marketID domain.MarketID) (Attempt, error) {
	item, err := s.marketRepository.GetByID(ctx, marketID)
	if err != nil {
		return Attempt{}, fmt.Errorf("get market for settlement scan: %w", err)
	}

	return attemptFromMarket(item), nil
}

func (s *service) SettleDueMarkets(ctx context.Context, filter DueMarketFilter) ([]Attempt, error) {
	normalized := normalizeDueMarketFilter(filter)

	items, err := s.marketRepository.ListExpiringBefore(ctx, normalized.Before, normalized.Limit)
	if err != nil {
		return nil, fmt.Errorf("list due markets for settlement scan: %w", err)
	}

	attempts := make([]Attempt, 0, len(items))
	for _, item := range items {
		attempts = append(attempts, attemptFromMarket(item))
	}

	return attempts, nil
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
