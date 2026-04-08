package storage

import (
	"context"
	"errors"
	"fmt"

	"prediction/internal/domain"

	"github.com/jackc/pgx/v5"
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

func (r *BalancePostgresRepository) LockStake(ctx context.Context, input LockStakeInput) (Balance, error) {
	const query = `
UPDATE player_balances
SET
    available_balance = available_balance - $2::numeric,
    locked_balance = locked_balance + $2::numeric,
    updated_at = NOW()
WHERE player_id = $1
  AND available_balance >= $2::numeric
RETURNING player_id, available_balance, locked_balance, updated_at;
`

	var balance Balance
	err := r.queryer.QueryRow(ctx, query, string(input.PlayerID), input.Amount).Scan(
		&balance.PlayerID,
		&balance.AvailableBalance,
		&balance.LockedBalance,
		&balance.UpdatedAt,
	)
	if err == nil {
		balance.UpdatedAt = domain.NormalizeTime(balance.UpdatedAt)
		return balance, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		exists, existsErr := r.balanceExists(ctx, input.PlayerID)
		if existsErr != nil {
			return Balance{}, existsErr
		}
		if !exists {
			return Balance{}, domain.ErrNotFound
		}

		return Balance{}, domain.NewValidationError("stake_amount", "insufficient available balance", input.Amount)
	}

	return Balance{}, fmt.Errorf("lock stake: %w", err)
}

func (r *BalancePostgresRepository) balanceExists(ctx context.Context, playerID domain.PlayerID) (bool, error) {
	var exists bool
	if err := r.queryer.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM player_balances WHERE player_id = $1)", string(playerID)).Scan(&exists); err != nil {
		return false, fmt.Errorf("check balance existence: %w", err)
	}

	return exists, nil
}
