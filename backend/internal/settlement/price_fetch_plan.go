package settlement

import (
	"sort"
	"time"

	"prediction/internal/domain"
	"prediction/internal/storage"
)

func normalizePriceFetchPlanFilter(filter PriceFetchPlanFilter) PriceFetchPlanFilter {
	if filter.After.IsZero() {
		filter.After = domain.NowUTC()
	} else {
		filter.After = domain.NormalizeTime(filter.After)
	}

	if filter.Before.IsZero() || filter.Before.Before(filter.After) {
		filter.Before = filter.After.Add(2 * time.Minute)
	} else {
		filter.Before = domain.NormalizeTime(filter.Before)
	}

	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	return filter
}

func buildPriceFetchBatches(items []storage.Market, after time.Time) []PriceFetchBatch {
	batchesByExpiry := map[time.Time][]PriceFetchTarget{}
	for _, item := range items {
		if !isPlannablePriceMarket(item, after) {
			continue
		}

		expiryTime := normalizeExpirySecond(item.ExpiryTime)
		batchesByExpiry[expiryTime] = append(batchesByExpiry[expiryTime], PriceFetchTarget{
			MarketID:   item.ID,
			Symbol:     item.Symbol,
			ExpiryTime: expiryTime,
		})
	}

	expiries := make([]time.Time, 0, len(batchesByExpiry))
	for expiry := range batchesByExpiry {
		expiries = append(expiries, expiry)
	}
	sort.Slice(expiries, func(i, j int) bool {
		return expiries[i].Before(expiries[j])
	})

	batches := make([]PriceFetchBatch, 0, len(expiries))
	for _, expiry := range expiries {
		targets := batchesByExpiry[expiry]
		sort.Slice(targets, func(i, j int) bool {
			if targets[i].Symbol == targets[j].Symbol {
				return string(targets[i].MarketID) < string(targets[j].MarketID)
			}

			return targets[i].Symbol < targets[j].Symbol
		})

		batches = append(batches, PriceFetchBatch{
			ExpiryTime: expiry,
			Symbols:    uniqueBatchSymbols(targets),
			Targets:    targets,
		})
	}

	return batches
}

func isPlannablePriceMarket(item storage.Market, after time.Time) bool {
	if item.Status != domain.MarketStatusActive {
		return false
	}
	if item.MarketType != domain.MarketTypePriceThreshold {
		return false
	}
	if item.SourceType != domain.SourceTypeMarkPrice {
		return false
	}

	return !item.ExpiryTime.Before(after)
}

func normalizeExpirySecond(value time.Time) time.Time {
	return domain.NormalizeTime(value).Truncate(time.Second)
}

func uniqueBatchSymbols(targets []PriceFetchTarget) []string {
	symbols := make([]string, 0, len(targets))
	seen := make(map[string]struct{}, len(targets))
	for _, target := range targets {
		if _, ok := seen[target.Symbol]; ok {
			continue
		}

		seen[target.Symbol] = struct{}{}
		symbols = append(symbols, target.Symbol)
	}

	sort.Strings(symbols)
	return symbols
}
