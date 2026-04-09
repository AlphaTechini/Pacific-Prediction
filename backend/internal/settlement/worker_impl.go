package settlement

import (
	"context"
	"fmt"
	"log"
	"time"

	"prediction/internal/domain"
)

type worker struct {
	logger             *log.Logger
	service            Service
	scanInterval       time.Duration
	scanBatchSize      int
	priceLookahead     time.Duration
	priceRetryInterval time.Duration
}

type WorkerDeps struct {
	Logger             *log.Logger
	Service            Service
	ScanInterval       time.Duration
	ScanBatchSize      int
	PriceLookahead     time.Duration
	PriceRetryInterval time.Duration
}

func NewWorker(deps WorkerDeps) Worker {
	logger := deps.Logger
	if logger == nil {
		logger = log.Default()
	}

	return &worker{
		logger:             logger,
		service:            deps.Service,
		scanInterval:       deps.ScanInterval,
		scanBatchSize:      deps.ScanBatchSize,
		priceLookahead:     deps.PriceLookahead,
		priceRetryInterval: deps.PriceRetryInterval,
	}
}

func (w *worker) Run(ctx context.Context) error {
	if err := w.scanOnce(ctx); err != nil {
		w.logger.Printf("settlement expiry scan failed: %v", err)
	}

	ticker := time.NewTicker(w.scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := w.scanOnce(ctx); err != nil {
				w.logger.Printf("settlement expiry scan failed: %v", err)
			}
		}
	}
}

func (w *worker) scanOnce(ctx context.Context) error {
	now := domain.NowUTC()

	attempts, err := w.service.SettleDueMarkets(ctx, DueMarketFilter{
		Before: now,
		Limit:  w.scanBatchSize,
	})
	if err != nil {
		return fmt.Errorf("scan due markets: %w", err)
	}

	if len(attempts) == 0 {
		if err := w.logUpcomingPriceFetchBatches(ctx, now); err != nil {
			return err
		}

		return nil
	}

	w.logger.Printf("settlement expiry scan discovered %d due market(s)", len(attempts))

	return w.logUpcomingPriceFetchBatches(ctx, now)
}

func (w *worker) logUpcomingPriceFetchBatches(ctx context.Context, now time.Time) error {
	batches, err := w.service.PlanPriceFetchBatches(ctx, PriceFetchPlanFilter{
		After:  now,
		Before: now.Add(w.priceLookahead),
		Limit:  w.scanBatchSize,
	})
	if err != nil {
		return fmt.Errorf("plan upcoming price fetch batches: %w", err)
	}

	if len(batches) == 0 {
		return nil
	}

	for _, batch := range batches {
		untilExpiry := batch.ExpiryTime.Sub(now)
		if untilExpiry < 0 {
			untilExpiry = 0
		}

		w.logger.Printf(
			"settlement price fetch plan expiry=%s symbols=%d markets=%d retry_interval=%s until_expiry=%s",
			batch.ExpiryTime.Format(time.RFC3339),
			len(batch.Symbols),
			len(batch.Targets),
			w.priceRetryInterval,
			untilExpiry,
		)
	}

	return nil
}
