package settlement

import (
	"context"
	"fmt"
	"log"
	"time"

	"prediction/internal/domain"
)

type worker struct {
	logger        *log.Logger
	service       Service
	scanInterval  time.Duration
	scanBatchSize int
}

type WorkerDeps struct {
	Logger        *log.Logger
	Service       Service
	ScanInterval  time.Duration
	ScanBatchSize int
}

func NewWorker(deps WorkerDeps) Worker {
	logger := deps.Logger
	if logger == nil {
		logger = log.Default()
	}

	return &worker{
		logger:        logger,
		service:       deps.Service,
		scanInterval:  deps.ScanInterval,
		scanBatchSize: deps.ScanBatchSize,
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
	attempts, err := w.service.SettleDueMarkets(ctx, DueMarketFilter{
		Before: domain.NowUTC(),
		Limit:  w.scanBatchSize,
	})
	if err != nil {
		return fmt.Errorf("scan due markets: %w", err)
	}

	if len(attempts) == 0 {
		return nil
	}

	w.logger.Printf("settlement expiry scan discovered %d due market(s)", len(attempts))
	return nil
}
