package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"prediction/internal/domain"

	"github.com/jackc/pgx/v5"
)

type MarketPostgresRepository struct {
	queryer Queryer
}

func NewMarketPostgresRepository(queryer Queryer) *MarketPostgresRepository {
	return &MarketPostgresRepository{queryer: queryer}
}

func (r *MarketPostgresRepository) Create(ctx context.Context, input CreateMarketInput) (Market, error) {
	const query = `
INSERT INTO markets (
    id,
    title,
    symbol,
    market_type,
    condition_operator,
    threshold_value,
    source_type,
    source_interval,
    reference_value,
    expiry_time,
    created_by_player_id
)
VALUES ($1, $2, $3, $4, $5, NULLIF($6, '')::numeric, $7, NULLIF($8, ''), NULLIF($9, '')::numeric, $10, $11)
RETURNING
    id,
    title,
    symbol,
    market_type,
    condition_operator,
    COALESCE(threshold_value::text, ''),
    source_type,
    COALESCE(source_interval, ''),
    COALESCE(reference_value::text, ''),
    expiry_time,
    status,
    COALESCE(result, ''),
    COALESCE(settlement_value::text, ''),
    resolved_at,
    COALESCE(resolution_reason, ''),
    created_by_player_id,
    created_at;
`

	return r.scanMarketRow(
		r.queryer.QueryRow(
			ctx,
			query,
			string(input.ID),
			input.Title,
			input.Symbol,
			string(input.MarketType),
			string(input.ConditionOperator),
			input.ThresholdValue,
			string(input.SourceType),
			input.SourceInterval,
			input.ReferenceValue,
			input.ExpiryTime,
			string(input.CreatedByPlayerID),
		),
		"create market",
	)
}

func (r *MarketPostgresRepository) GetByID(ctx context.Context, marketID domain.MarketID) (Market, error) {
	const query = `
SELECT
    id,
    title,
    symbol,
    market_type,
    condition_operator,
    COALESCE(threshold_value::text, ''),
    source_type,
    COALESCE(source_interval, ''),
    COALESCE(reference_value::text, ''),
    expiry_time,
    status,
    COALESCE(result, ''),
    COALESCE(settlement_value::text, ''),
    resolved_at,
    COALESCE(resolution_reason, ''),
    created_by_player_id,
    created_at
FROM markets
WHERE id = $1;
`

	return r.scanMarketRow(r.queryer.QueryRow(ctx, query, string(marketID)), "get market by id")
}

func (r *MarketPostgresRepository) ListByStatus(ctx context.Context, status domain.MarketStatus, limit int) ([]Market, error) {
	const query = `
SELECT
    id,
    title,
    symbol,
    market_type,
    condition_operator,
    COALESCE(threshold_value::text, ''),
    source_type,
    COALESCE(source_interval, ''),
    COALESCE(reference_value::text, ''),
    expiry_time,
    status,
    COALESCE(result, ''),
    COALESCE(settlement_value::text, ''),
    resolved_at,
    COALESCE(resolution_reason, ''),
    created_by_player_id,
    created_at
FROM markets
WHERE status = $1
ORDER BY expiry_time ASC
LIMIT $2;
`

	rows, err := r.queryer.Query(ctx, query, string(status), limit)
	if err != nil {
		return nil, fmt.Errorf("list markets by status: %w", err)
	}
	defer rows.Close()

	var markets []Market
	for rows.Next() {
		market, err := scanMarket(rows)
		if err != nil {
			return nil, fmt.Errorf("list markets by status: %w", err)
		}
		markets = append(markets, market)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list markets by status: %w", err)
	}

	return markets, nil
}

func (r *MarketPostgresRepository) ListExpiringBefore(ctx context.Context, before time.Time, limit int) ([]Market, error) {
	const query = `
SELECT
    id,
    title,
    symbol,
    market_type,
    condition_operator,
    COALESCE(threshold_value::text, ''),
    source_type,
    COALESCE(source_interval, ''),
    COALESCE(reference_value::text, ''),
    expiry_time,
    status,
    COALESCE(result, ''),
    COALESCE(settlement_value::text, ''),
    resolved_at,
    COALESCE(resolution_reason, ''),
    created_by_player_id,
    created_at
FROM markets
WHERE status = $1
  AND expiry_time <= $2
ORDER BY expiry_time ASC
LIMIT $3;
`

	rows, err := r.queryer.Query(ctx, query, string(domain.MarketStatusActive), before.UTC(), limit)
	if err != nil {
		return nil, fmt.Errorf("list expiring markets: %w", err)
	}
	defer rows.Close()

	var markets []Market
	for rows.Next() {
		market, err := scanMarket(rows)
		if err != nil {
			return nil, fmt.Errorf("list expiring markets: %w", err)
		}
		markets = append(markets, market)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list expiring markets: %w", err)
	}

	return markets, nil
}

func (r *MarketPostgresRepository) UpdateSettlement(ctx context.Context, input UpdateMarketSettlementInput) (Market, error) {
	const query = `
UPDATE markets
SET
    status = $2,
    result = $3,
    settlement_value = NULLIF($4, '')::numeric,
    resolved_at = $5,
    resolution_reason = $6
WHERE id = $1
RETURNING
    id,
    title,
    symbol,
    market_type,
    condition_operator,
    COALESCE(threshold_value::text, ''),
    source_type,
    COALESCE(source_interval, ''),
    COALESCE(reference_value::text, ''),
    expiry_time,
    status,
    COALESCE(result, ''),
    COALESCE(settlement_value::text, ''),
    resolved_at,
    COALESCE(resolution_reason, ''),
    created_by_player_id,
    created_at;
`

	return r.scanMarketRow(
		r.queryer.QueryRow(
			ctx,
			query,
			string(input.MarketID),
			string(input.Status),
			string(input.Result),
			input.SettlementValue,
			input.ResolvedAt,
			input.ResolutionReason,
		),
		"update market settlement",
	)
}

func (r *MarketPostgresRepository) scanMarketRow(row interface{ Scan(...any) error }, operation string) (Market, error) {
	market, err := scanMarket(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Market{}, domain.ErrNotFound
		}

		return Market{}, fmt.Errorf("%s: %w", operation, err)
	}

	return market, nil
}

func scanMarket(row interface{ Scan(...any) error }) (Market, error) {
	var market Market
	var resolvedAt *time.Time

	if err := row.Scan(
		&market.ID,
		&market.Title,
		&market.Symbol,
		&market.MarketType,
		&market.ConditionOperator,
		&market.ThresholdValue,
		&market.SourceType,
		&market.SourceInterval,
		&market.ReferenceValue,
		&market.ExpiryTime,
		&market.Status,
		&market.Result,
		&market.SettlementValue,
		&resolvedAt,
		&market.ResolutionReason,
		&market.CreatedByPlayerID,
		&market.CreatedAt,
	); err != nil {
		return Market{}, err
	}

	market.ExpiryTime = domain.NormalizeTime(market.ExpiryTime)
	market.CreatedAt = domain.NormalizeTime(market.CreatedAt)
	if resolvedAt != nil {
		value := domain.NormalizeTime(*resolvedAt)
		market.ResolvedAt = &value
	}

	return market, nil
}
