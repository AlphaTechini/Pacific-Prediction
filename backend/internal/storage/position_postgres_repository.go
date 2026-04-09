package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"prediction/internal/domain"

	"github.com/jackc/pgx/v5"
)

type PositionPostgresRepository struct {
	queryer Queryer
}

func NewPositionPostgresRepository(queryer Queryer) *PositionPostgresRepository {
	return &PositionPostgresRepository{queryer: queryer}
}

func (r *PositionPostgresRepository) Create(ctx context.Context, input CreatePositionInput) (Position, error) {
	const query = `
INSERT INTO positions (
    id,
    player_id,
    market_id,
    side,
    stake_amount,
    potential_payout
)
VALUES ($1, $2, $3, $4, $5::numeric, $6::numeric)
RETURNING
    id,
    player_id,
    market_id,
    side,
    stake_amount::text,
    potential_payout::text,
    status,
    created_at,
    settled_at;
`

	return r.scanPositionRow(
		r.queryer.QueryRow(
			ctx,
			query,
			string(input.ID),
			string(input.PlayerID),
			string(input.MarketID),
			string(input.Side),
			input.StakeAmount,
			input.PotentialPayout,
		),
		"create position",
	)
}

func (r *PositionPostgresRepository) ListByPlayerID(ctx context.Context, playerID domain.PlayerID, limit int) ([]Position, error) {
	const query = `
SELECT
    id,
    player_id,
    market_id,
    side,
    stake_amount::text,
    potential_payout::text,
    status,
    created_at,
    settled_at
FROM positions
WHERE player_id = $1
ORDER BY created_at DESC
LIMIT $2;
`

	rows, err := r.queryer.Query(ctx, query, string(playerID), limit)
	if err != nil {
		return nil, fmt.Errorf("list positions by player id: %w", err)
	}
	defer rows.Close()

	return scanPositions(rows, "list positions by player id")
}

func (r *PositionPostgresRepository) ListByMarketID(ctx context.Context, marketID domain.MarketID) ([]Position, error) {
	const query = `
SELECT
    id,
    player_id,
    market_id,
    side,
    stake_amount::text,
    potential_payout::text,
    status,
    created_at,
    settled_at
FROM positions
WHERE market_id = $1
ORDER BY created_at ASC;
`

	rows, err := r.queryer.Query(ctx, query, string(marketID))
	if err != nil {
		return nil, fmt.Errorf("list positions by market id: %w", err)
	}
	defer rows.Close()

	return scanPositions(rows, "list positions by market id")
}

func (r *PositionPostgresRepository) UpdateSettlement(ctx context.Context, input UpdatePositionSettlementInput) (Position, error) {
	const query = `
UPDATE positions
SET
    status = $2,
    settled_at = $3
WHERE id = $1
RETURNING
    id,
    player_id,
    market_id,
    side,
    stake_amount::text,
    potential_payout::text,
    status,
    created_at,
    settled_at;
`

	return r.scanPositionRow(
		r.queryer.QueryRow(
			ctx,
			query,
			string(input.PositionID),
			string(input.Status),
			input.SettledAt,
		),
		"update position settlement",
	)
}

func (r *PositionPostgresRepository) scanPositionRow(row interface{ Scan(...any) error }, operation string) (Position, error) {
	position, err := scanPosition(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Position{}, domain.ErrNotFound
		}

		return Position{}, fmt.Errorf("%s: %w", operation, err)
	}

	return position, nil
}

func scanPositions(rows pgx.Rows, operation string) ([]Position, error) {
	var positions []Position
	for rows.Next() {
		position, err := scanPosition(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}

		positions = append(positions, position)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return positions, nil
}

func scanPosition(row interface{ Scan(...any) error }) (Position, error) {
	var position Position
	var settledAt *time.Time

	if err := row.Scan(
		&position.ID,
		&position.PlayerID,
		&position.MarketID,
		&position.Side,
		&position.StakeAmount,
		&position.PotentialPayout,
		&position.Status,
		&position.CreatedAt,
		&settledAt,
	); err != nil {
		return Position{}, err
	}

	position.CreatedAt = domain.NormalizeTime(position.CreatedAt)
	if settledAt != nil {
		value := domain.NormalizeTime(*settledAt)
		position.SettledAt = &value
	}

	return position, nil
}
