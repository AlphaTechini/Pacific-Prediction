package storage

import (
	"context"
	"fmt"

	"prediction/internal/domain"
)

type BalancePostgresRepository struct {
	queryer Queryer
}

func NewBalancePostgresRepository(queryer Queryer) *BalancePostgresRepository {
	return &BalancePostgresRepository{queryer: queryer}
}

func (r *BalancePostgresRepository) Create(ctx context.Context, input CreateBalanceInput) (Balance, error) {
	const query = `
INSERT INTO player_balances (player_id, available_balance, locked_balance)
VALUES ($1, $2, $3)
RETURNING player_id, available_balance, locked_balance, updated_at;
`

	var balance Balance
	if err := r.queryer.QueryRow(ctx, query, string(input.PlayerID), input.AvailableBalance, input.LockedBalance).Scan(
		&balance.PlayerID,
		&balance.AvailableBalance,
		&balance.LockedBalance,
		&balance.UpdatedAt,
	); err != nil {
		return Balance{}, fmt.Errorf("create balance: %w", err)
	}

	balance.UpdatedAt = domain.NormalizeTime(balance.UpdatedAt)

	return balance, nil
}

func (r *BalancePostgresRepository) GetByPlayerID(ctx context.Context, playerID domain.PlayerID) (Balance, error) {
	const query = `
SELECT player_id, available_balance, locked_balance, updated_at
FROM player_balances
WHERE player_id = $1;
`

	var balance Balance
	if err := r.queryer.QueryRow(ctx, query, string(playerID)).Scan(
		&balance.PlayerID,
		&balance.AvailableBalance,
		&balance.LockedBalance,
		&balance.UpdatedAt,
	); err != nil {
		return Balance{}, fmt.Errorf("get balance by player id: %w", err)
	}

	balance.UpdatedAt = domain.NormalizeTime(balance.UpdatedAt)

	return balance, nil
}
