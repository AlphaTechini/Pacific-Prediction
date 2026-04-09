package settlement

import (
	"testing"
	"time"

	"prediction/internal/domain"
	"prediction/internal/storage"
)

func TestBuildPriceFetchBatchesGroupsUpcomingPriceMarketsByExpirySecond(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)
	expiryA := now.Add(45 * time.Second)
	expiryB := now.Add(90 * time.Second)

	items := []storage.Market{
		{
			ID:         domain.MarketID("m-1"),
			Symbol:     "BTC-PERP",
			MarketType: domain.MarketTypePriceThreshold,
			SourceType: domain.SourceTypeMarkPrice,
			Status:     domain.MarketStatusActive,
			ExpiryTime: expiryA.Add(200 * time.Millisecond),
		},
		{
			ID:         domain.MarketID("m-2"),
			Symbol:     "ETH-PERP",
			MarketType: domain.MarketTypePriceThreshold,
			SourceType: domain.SourceTypeMarkPrice,
			Status:     domain.MarketStatusActive,
			ExpiryTime: expiryA.Add(900 * time.Millisecond),
		},
		{
			ID:         domain.MarketID("m-3"),
			Symbol:     "BTC-PERP",
			MarketType: domain.MarketTypePriceThreshold,
			SourceType: domain.SourceTypeMarkPrice,
			Status:     domain.MarketStatusActive,
			ExpiryTime: expiryA.Add(300 * time.Millisecond),
		},
		{
			ID:         domain.MarketID("m-4"),
			Symbol:     "SOL-PERP",
			MarketType: domain.MarketTypePriceThreshold,
			SourceType: domain.SourceTypeMarkPrice,
			Status:     domain.MarketStatusActive,
			ExpiryTime: expiryB,
		},
		{
			ID:         domain.MarketID("m-5"),
			Symbol:     "XRP-PERP",
			MarketType: domain.MarketTypeCandleDirection,
			SourceType: domain.SourceTypeMarkPriceCandle,
			Status:     domain.MarketStatusActive,
			ExpiryTime: expiryA,
		},
		{
			ID:         domain.MarketID("m-6"),
			Symbol:     "ADA-PERP",
			MarketType: domain.MarketTypePriceThreshold,
			SourceType: domain.SourceTypeMarkPrice,
			Status:     domain.MarketStatusResolved,
			ExpiryTime: expiryA,
		},
		{
			ID:         domain.MarketID("m-7"),
			Symbol:     "DOGE-PERP",
			MarketType: domain.MarketTypePriceThreshold,
			SourceType: domain.SourceTypeMarkPrice,
			Status:     domain.MarketStatusActive,
			ExpiryTime: now.Add(-5 * time.Second),
		},
	}

	batches := buildPriceFetchBatches(items, now)
	if len(batches) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(batches))
	}

	first := batches[0]
	if !first.ExpiryTime.Equal(expiryA.Truncate(time.Second)) {
		t.Fatalf("expected first batch expiry %s, got %s", expiryA.Truncate(time.Second), first.ExpiryTime)
	}
	if got, want := first.Symbols, []string{"BTC-PERP", "ETH-PERP"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("expected first batch symbols %v, got %v", want, got)
	}
	if len(first.Targets) != 3 {
		t.Fatalf("expected 3 targets in first batch, got %d", len(first.Targets))
	}

	second := batches[1]
	if !second.ExpiryTime.Equal(expiryB.Truncate(time.Second)) {
		t.Fatalf("expected second batch expiry %s, got %s", expiryB.Truncate(time.Second), second.ExpiryTime)
	}
	if got, want := second.Symbols, []string{"SOL-PERP"}; len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("expected second batch symbols %v, got %v", want, got)
	}
	if len(second.Targets) != 1 {
		t.Fatalf("expected 1 target in second batch, got %d", len(second.Targets))
	}
}
