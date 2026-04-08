package storage

import "context"

type Repositories interface {
	Players() PlayerRepository
	Sessions() SessionRepository
	Balances() BalanceRepository
	Markets() MarketRepository
	Positions() PositionRepository
	Settlements() SettlementRepository
}

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
}
