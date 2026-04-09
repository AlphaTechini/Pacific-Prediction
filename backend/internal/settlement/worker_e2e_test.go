package settlement

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"prediction/internal/domain"
	"prediction/internal/storage"
)

func TestWorkerScanOnceLogsDueAndUpcomingPriceFetchPlan(t *testing.T) {
	now := domain.NowUTC()
	sharedExpiry := now.Add(40 * time.Second).Truncate(time.Second)
	repo := &fakeMarketRepository{
		markets: []storage.Market{
			{
				ID:         domain.MarketID("due-price"),
				Symbol:     "BTC-PERP",
				MarketType: domain.MarketTypePriceThreshold,
				SourceType: domain.SourceTypeMarkPrice,
				Status:     domain.MarketStatusActive,
				ExpiryTime: now.Add(-1 * time.Second),
			},
			{
				ID:         domain.MarketID("soon-price-1"),
				Symbol:     "BTC-PERP",
				MarketType: domain.MarketTypePriceThreshold,
				SourceType: domain.SourceTypeMarkPrice,
				Status:     domain.MarketStatusActive,
				ExpiryTime: sharedExpiry,
			},
			{
				ID:         domain.MarketID("soon-price-2"),
				Symbol:     "ETH-PERP",
				MarketType: domain.MarketTypePriceThreshold,
				SourceType: domain.SourceTypeMarkPrice,
				Status:     domain.MarketStatusActive,
				ExpiryTime: sharedExpiry.Add(400 * time.Millisecond),
			},
			{
				ID:         domain.MarketID("later-price"),
				Symbol:     "SOL-PERP",
				MarketType: domain.MarketTypePriceThreshold,
				SourceType: domain.SourceTypeMarkPrice,
				Status:     domain.MarketStatusActive,
				ExpiryTime: now.Add(5 * time.Minute),
			},
		},
	}

	service := NewService(ServiceDeps{
		MarketRepository: repo,
		PriceResolver: &fakePriceResolver{
			batchErr: errSettlementSourceNotReady,
		},
	})

	var output bytes.Buffer
	testWorker := NewWorker(WorkerDeps{
		Logger:             newTestLogger(&output),
		Service:            service,
		ScanInterval:       15 * time.Second,
		ScanBatchSize:      50,
		PriceLookahead:     2 * time.Minute,
		PriceRetryInterval: 500 * time.Millisecond,
	})

	impl, ok := testWorker.(*worker)
	if !ok {
		t.Fatalf("expected concrete worker implementation")
	}

	if err := impl.scanOnce(context.Background()); err != nil {
		t.Fatalf("scanOnce returned error: %v", err)
	}

	logs := output.String()
	assertLogContains(t, logs, "settlement expiry scan discovered 1 due market(s)")
	assertLogContains(t, logs, "symbols=2")
	assertLogContains(t, logs, "markets=2")
	assertLogContains(t, logs, "retry_interval=500ms")
}

type fakeMarketRepository struct {
	markets []storage.Market
}

func (r *fakeMarketRepository) Create(context.Context, storage.CreateMarketInput) (storage.Market, error) {
	panic("unexpected Create call")
}

func (r *fakeMarketRepository) GetByID(_ context.Context, marketID domain.MarketID) (storage.Market, error) {
	for _, item := range r.markets {
		if item.ID == marketID {
			return item, nil
		}
	}

	return storage.Market{}, domain.ErrNotFound
}

func (r *fakeMarketRepository) ListByStatus(_ context.Context, status domain.MarketStatus, limit int) ([]storage.Market, error) {
	var items []storage.Market
	for _, item := range r.markets {
		if item.Status != status {
			continue
		}

		items = append(items, item)
		if limit > 0 && len(items) >= limit {
			break
		}
	}

	return items, nil
}

func (r *fakeMarketRepository) ListExpiringBefore(_ context.Context, before time.Time, limit int) ([]storage.Market, error) {
	var items []storage.Market
	for _, item := range r.markets {
		if item.Status != domain.MarketStatusActive {
			continue
		}
		if item.ExpiryTime.After(before) {
			continue
		}

		items = append(items, item)
		if limit > 0 && len(items) >= limit {
			break
		}
	}

	return items, nil
}

func (r *fakeMarketRepository) UpdateSettlement(context.Context, storage.UpdateMarketSettlementInput) (storage.Market, error) {
	panic("unexpected UpdateSettlement call")
}

func newTestLogger(output *bytes.Buffer) *log.Logger {
	return log.New(output, "", 0)
}

func assertLogContains(t *testing.T, logs, fragment string) {
	t.Helper()

	if !strings.Contains(logs, fragment) {
		t.Fatalf("expected logs to contain %q, got %q", fragment, logs)
	}
}
