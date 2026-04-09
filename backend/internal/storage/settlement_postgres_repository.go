package storage

import (
	"context"
	"errors"
	"fmt"

	"prediction/internal/domain"

	"github.com/jackc/pgx/v5"
)

type SettlementPostgresRepository struct {
	queryer Queryer
}

func NewSettlementPostgresRepository(queryer Queryer) *SettlementPostgresRepository {
	return &SettlementPostgresRepository{queryer: queryer}
}

func (r *SettlementPostgresRepository) Create(ctx context.Context, input CreateSettlementInput) (Settlement, error) {
	const query = `
INSERT INTO market_settlements (
    id,
    market_id,
    pacifica_source,
    source_timestamp,
    raw_payload,
    settlement_value,
    result
)
VALUES ($1, $2, $3, $4, $5::jsonb, NULLIF($6, '')::numeric, $7)
RETURNING
    id,
    market_id,
    pacifica_source,
    source_timestamp,
    raw_payload::text,
    COALESCE(settlement_value::text, ''),
    result,
    created_at;
`

	return r.scanSettlementRow(
		r.queryer.QueryRow(
			ctx,
			query,
			string(input.ID),
			string(input.MarketID),
			input.PacificaSource,
			input.SourceTimestamp,
			string(input.RawPayload),
			input.SettlementValue,
			string(input.Result),
		),
		"create settlement",
	)
}

func (r *SettlementPostgresRepository) GetByMarketID(ctx context.Context, marketID domain.MarketID) (Settlement, error) {
	const query = `
SELECT
    id,
    market_id,
    pacifica_source,
    source_timestamp,
    raw_payload::text,
    COALESCE(settlement_value::text, ''),
    result,
    created_at
FROM market_settlements
WHERE market_id = $1;
`

	return r.scanSettlementRow(r.queryer.QueryRow(ctx, query, string(marketID)), "get settlement by market id")
}

func (r *SettlementPostgresRepository) scanSettlementRow(row interface{ Scan(...any) error }, operation string) (Settlement, error) {
	settlement, err := scanSettlement(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Settlement{}, domain.ErrNotFound
		}

		return Settlement{}, fmt.Errorf("%s: %w", operation, err)
	}

	return settlement, nil
}

func scanSettlement(row interface{ Scan(...any) error }) (Settlement, error) {
	var settlement Settlement
	var rawPayload string

	if err := row.Scan(
		&settlement.ID,
		&settlement.MarketID,
		&settlement.PacificaSource,
		&settlement.SourceTimestamp,
		&rawPayload,
		&settlement.SettlementValue,
		&settlement.Result,
		&settlement.CreatedAt,
	); err != nil {
		return Settlement{}, err
	}

	settlement.SourceTimestamp = domain.NormalizeTime(settlement.SourceTimestamp)
	settlement.CreatedAt = domain.NormalizeTime(settlement.CreatedAt)
	settlement.RawPayload = []byte(rawPayload)

	return settlement, nil
}
