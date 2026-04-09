package settlement

import (
	"context"
	"errors"
	"testing"
	"time"

	"prediction/internal/domain"
	"prediction/internal/pacifica"
	"prediction/internal/storage"

	"github.com/jackc/pgx/v5"
)

func TestPriceResolverResolveBatchUsesSingleFetchAndProducesResolutions(t *testing.T) {
	now := time.Date(2026, 4, 9, 0, 10, 0, 0, time.UTC)
	client := &fakePacificaClient{
		priceResponses: [][]pacifica.PriceSnapshot{
			{
				{
					Symbol:      "BTC",
					MarkPrice:   "71000.10",
					OraclePrice: "71010.00",
					Timestamp:   now.Add(3 * time.Second),
				},
				{
					Symbol:      "ETH",
					MarkPrice:   "2200.00",
					OraclePrice: "2201.00",
					Timestamp:   now.Add(3 * time.Second),
				},
			},
		},
	}

	resolver := NewPriceResolver(client)
	resolutions, err := resolver.ResolveBatch(context.Background(), []PriceMarket{
		{
			ID:                domain.MarketID("market_btc"),
			Symbol:            "BTC",
			ConditionOperator: domain.ConditionOperatorGT,
			ThresholdValue:    "70000",
			ExpiryTime:        now,
		},
		{
			ID:                domain.MarketID("market_eth"),
			Symbol:            "ETH",
			ConditionOperator: domain.ConditionOperatorLT,
			ThresholdValue:    "2300",
			ExpiryTime:        now,
		},
	})
	if err != nil {
		t.Fatalf("ResolveBatch returned error: %v", err)
	}

	if client.listPricesCalls != 1 {
		t.Fatalf("expected one batched ListPrices call, got %d", client.listPricesCalls)
	}
	if len(resolutions) != 2 {
		t.Fatalf("expected 2 resolutions, got %d", len(resolutions))
	}
	if resolutions[0].Result != domain.MarketResultYes {
		t.Fatalf("expected BTC result yes, got %s", resolutions[0].Result)
	}
	if resolutions[1].Result != domain.MarketResultYes {
		t.Fatalf("expected ETH result yes, got %s", resolutions[1].Result)
	}
}

func TestPriceResolverResolveReturnsNotReadyWhenPacificaTimestampPredatesExpiry(t *testing.T) {
	now := time.Date(2026, 4, 9, 0, 11, 0, 0, time.UTC)
	client := &fakePacificaClient{
		priceResponses: [][]pacifica.PriceSnapshot{
			{
				{
					Symbol:    "BTC",
					MarkPrice: "71000.10",
					Timestamp: now.Add(-1 * time.Second),
				},
			},
		},
	}

	resolver := NewPriceResolver(client)
	_, err := resolver.Resolve(context.Background(), PriceMarket{
		ID:                domain.MarketID("market_btc"),
		Symbol:            "BTC",
		ConditionOperator: domain.ConditionOperatorGT,
		ThresholdValue:    "70000",
		ExpiryTime:        now,
	})
	if !errors.Is(err, errSettlementSourceNotReady) {
		t.Fatalf("expected errSettlementSourceNotReady, got %v", err)
	}
}

func TestSettleDueMarketsSettlesPriceMarketWhenPacificaTimestampQualifies(t *testing.T) {
	now := time.Date(2026, 4, 9, 0, 12, 0, 0, time.UTC)
	repo := &fakeMarketRepository{
		markets: []storage.Market{
			{
				ID:                domain.MarketID("market_btc"),
				Symbol:            "BTC",
				MarketType:        domain.MarketTypePriceThreshold,
				ConditionOperator: domain.ConditionOperatorGT,
				ThresholdValue:    "70000",
				SourceType:        domain.SourceTypeMarkPrice,
				Status:            domain.MarketStatusActive,
				ExpiryTime:        now,
			},
		},
	}
	writeRepo := &fakeMarketWriteRepository{}
	settlementRepo := &fakeSettlementRepository{}
	positionRepo := &fakePositionRepository{}
	balanceRepo := &fakeBalanceRepository{}
	service := NewService(ServiceDeps{
		MarketRepository: repo,
		PriceResolver: &fakePriceResolver{
			batchResolutions: []PriceResolution{
				{
					MarketID:            domain.MarketID("market_btc"),
					PacificaSource:      pacificaPricesSource,
					SourceTimestamp:     now.Add(2 * time.Second),
					SettlementMarkPrice: "71000.10",
					Result:              domain.MarketResultYes,
					RawPayload:          []byte(`{"symbol":"BTC","mark":"71000.10"}`),
				},
			},
		},
		TxManager: &fakeSettlementTxManager{},
		MarketRepositoryFactory: func(storage.Queryer) storage.MarketRepository {
			return writeRepo
		},
		PositionRepositoryFactory: func(storage.Queryer) storage.PositionRepository {
			return positionRepo
		},
		BalanceRepositoryFactory: func(storage.Queryer) storage.BalanceRepository {
			return balanceRepo
		},
		SettlementRepositoryFactory: func(storage.Queryer) storage.SettlementRepository {
			return settlementRepo
		},
	})

	attempts, err := service.SettleDueMarkets(context.Background(), DueMarketFilter{
		Before: now,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("SettleDueMarkets returned error: %v", err)
	}

	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}
	if !attempts[0].Settled {
		t.Fatalf("expected settled attempt, got %+v", attempts[0])
	}
	if len(writeRepo.updates) != 1 {
		t.Fatalf("expected 1 market update, got %d", len(writeRepo.updates))
	}
	if len(settlementRepo.creates) != 1 {
		t.Fatalf("expected 1 settlement audit write, got %d", len(settlementRepo.creates))
	}
	if writeRepo.updates[0].Status != domain.MarketStatusResolved {
		t.Fatalf("expected resolved market status, got %s", writeRepo.updates[0].Status)
	}
}

func TestSettleDueMarketsLeavesPriceMarketUnsettledWhenSourceIsStillEarly(t *testing.T) {
	now := time.Date(2026, 4, 9, 0, 13, 0, 0, time.UTC)
	repo := &fakeMarketRepository{
		markets: []storage.Market{
			{
				ID:                domain.MarketID("market_btc"),
				Symbol:            "BTC",
				MarketType:        domain.MarketTypePriceThreshold,
				ConditionOperator: domain.ConditionOperatorGT,
				ThresholdValue:    "70000",
				SourceType:        domain.SourceTypeMarkPrice,
				Status:            domain.MarketStatusActive,
				ExpiryTime:        now,
			},
		},
	}
	writeRepo := &fakeMarketWriteRepository{}
	settlementRepo := &fakeSettlementRepository{}
	positionRepo := &fakePositionRepository{}
	balanceRepo := &fakeBalanceRepository{}
	service := NewService(ServiceDeps{
		MarketRepository: repo,
		PriceResolver: &fakePriceResolver{
			batchErr: errSettlementSourceNotReady,
		},
		TxManager: &fakeSettlementTxManager{},
		MarketRepositoryFactory: func(storage.Queryer) storage.MarketRepository {
			return writeRepo
		},
		PositionRepositoryFactory: func(storage.Queryer) storage.PositionRepository {
			return positionRepo
		},
		BalanceRepositoryFactory: func(storage.Queryer) storage.BalanceRepository {
			return balanceRepo
		},
		SettlementRepositoryFactory: func(storage.Queryer) storage.SettlementRepository {
			return settlementRepo
		},
	})

	attempts, err := service.SettleDueMarkets(context.Background(), DueMarketFilter{
		Before: now,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("SettleDueMarkets returned error: %v", err)
	}

	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}
	if attempts[0].Settled {
		t.Fatalf("expected unsettled attempt, got %+v", attempts[0])
	}
	if len(writeRepo.updates) != 0 {
		t.Fatalf("expected no market updates, got %d", len(writeRepo.updates))
	}
	if len(settlementRepo.creates) != 0 {
		t.Fatalf("expected no settlement audit writes, got %d", len(settlementRepo.creates))
	}
}

type fakePacificaClient struct {
	priceResponses  [][]pacifica.PriceSnapshot
	listPricesCalls int
}

func (c *fakePacificaClient) ListMarketInfo(context.Context) ([]pacifica.MarketInfo, error) {
	panic("unexpected ListMarketInfo call")
}

func (c *fakePacificaClient) ListPrices(_ context.Context, _ pacifica.PriceFilter) ([]pacifica.PriceSnapshot, error) {
	index := c.listPricesCalls
	c.listPricesCalls++

	if len(c.priceResponses) == 0 {
		return nil, nil
	}
	if index >= len(c.priceResponses) {
		return c.priceResponses[len(c.priceResponses)-1], nil
	}

	return c.priceResponses[index], nil
}

func (c *fakePacificaClient) ListMarkPriceCandles(context.Context, pacifica.MarkPriceCandleQuery) ([]pacifica.MarkPriceCandle, error) {
	panic("unexpected ListMarkPriceCandles call")
}

func (c *fakePacificaClient) ListFundingRateHistory(context.Context, pacifica.FundingRateHistoryQuery) ([]pacifica.FundingRateHistoryEntry, error) {
	panic("unexpected ListFundingRateHistory call")
}

type fakePriceResolver struct {
	singleResolution PriceResolution
	singleErr        error
	batchResolutions []PriceResolution
	batchErr         error
}

func (r *fakePriceResolver) Resolve(context.Context, PriceMarket) (PriceResolution, error) {
	return r.singleResolution, r.singleErr
}

func (r *fakePriceResolver) ResolveBatch(context.Context, []PriceMarket) ([]PriceResolution, error) {
	return r.batchResolutions, r.batchErr
}

type fakeSettlementTxManager struct{}

func (m *fakeSettlementTxManager) WithinTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return fn(nil)
}

type fakeMarketWriteRepository struct {
	updates []storage.UpdateMarketSettlementInput
}

func (r *fakeMarketWriteRepository) Create(context.Context, storage.CreateMarketInput) (storage.Market, error) {
	panic("unexpected Create call")
}

func (r *fakeMarketWriteRepository) GetByID(context.Context, domain.MarketID) (storage.Market, error) {
	panic("unexpected GetByID call")
}

func (r *fakeMarketWriteRepository) ListByStatus(context.Context, domain.MarketStatus, int) ([]storage.Market, error) {
	panic("unexpected ListByStatus call")
}

func (r *fakeMarketWriteRepository) ListExpiringBefore(context.Context, time.Time, int) ([]storage.Market, error) {
	panic("unexpected ListExpiringBefore call")
}

func (r *fakeMarketWriteRepository) UpdateSettlement(_ context.Context, input storage.UpdateMarketSettlementInput) (storage.Market, error) {
	r.updates = append(r.updates, input)
	return storage.Market{
		ID:              input.MarketID,
		Status:          input.Status,
		Result:          input.Result,
		SettlementValue: input.SettlementValue,
		ResolvedAt:      &input.ResolvedAt,
	}, nil
}

type fakeSettlementRepository struct {
	creates []storage.CreateSettlementInput
}

func (r *fakeSettlementRepository) Create(_ context.Context, input storage.CreateSettlementInput) (storage.Settlement, error) {
	r.creates = append(r.creates, input)
	return storage.Settlement{
		ID:              input.ID,
		MarketID:        input.MarketID,
		PacificaSource:  input.PacificaSource,
		SourceTimestamp: input.SourceTimestamp,
		RawPayload:      input.RawPayload,
		SettlementValue: input.SettlementValue,
		Result:          input.Result,
		CreatedAt:       input.SourceTimestamp,
	}, nil
}

func (r *fakeSettlementRepository) GetByMarketID(context.Context, domain.MarketID) (storage.Settlement, error) {
	panic("unexpected GetByMarketID call")
}

type fakePositionRepository struct {
	items   []storage.Position
	updates []storage.UpdatePositionSettlementInput
}

func (r *fakePositionRepository) Create(context.Context, storage.CreatePositionInput) (storage.Position, error) {
	panic("unexpected Create call")
}

func (r *fakePositionRepository) ListByPlayerID(context.Context, domain.PlayerID, int) ([]storage.Position, error) {
	panic("unexpected ListByPlayerID call")
}

func (r *fakePositionRepository) ListByMarketID(context.Context, domain.MarketID) ([]storage.Position, error) {
	return append([]storage.Position(nil), r.items...), nil
}

func (r *fakePositionRepository) UpdateSettlement(_ context.Context, input storage.UpdatePositionSettlementInput) (storage.Position, error) {
	r.updates = append(r.updates, input)
	return storage.Position{
		ID:        input.PositionID,
		Status:    input.Status,
		SettledAt: &input.SettledAt,
	}, nil
}

type fakeBalanceRepository struct {
	wins   []storage.SettleWonPositionInput
	losses []storage.SettleLostPositionInput
}

func (r *fakeBalanceRepository) Create(context.Context, storage.CreateBalanceInput) (storage.Balance, error) {
	panic("unexpected Create call")
}

func (r *fakeBalanceRepository) GetByPlayerID(context.Context, domain.PlayerID) (storage.Balance, error) {
	panic("unexpected GetByPlayerID call")
}

func (r *fakeBalanceRepository) LockStake(context.Context, storage.LockStakeInput) (storage.Balance, error) {
	panic("unexpected LockStake call")
}

func (r *fakeBalanceRepository) SettleWonPosition(_ context.Context, input storage.SettleWonPositionInput) (storage.Balance, error) {
	r.wins = append(r.wins, input)
	return storage.Balance{PlayerID: input.PlayerID}, nil
}

func (r *fakeBalanceRepository) SettleLostPosition(_ context.Context, input storage.SettleLostPositionInput) (storage.Balance, error) {
	r.losses = append(r.losses, input)
	return storage.Balance{PlayerID: input.PlayerID}, nil
}
