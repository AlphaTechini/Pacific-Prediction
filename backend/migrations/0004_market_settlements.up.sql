CREATE TABLE IF NOT EXISTS market_settlements (
    id TEXT PRIMARY KEY,
    market_id TEXT NOT NULL UNIQUE REFERENCES markets(id) ON DELETE CASCADE,
    pacifica_source TEXT NOT NULL,
    source_timestamp TIMESTAMPTZ NOT NULL,
    raw_payload JSONB NOT NULL,
    settlement_value NUMERIC(20, 8),
    result TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_market_settlements_result CHECK (result IN ('yes', 'no'))
);

CREATE INDEX IF NOT EXISTS idx_market_settlements_market_id
    ON market_settlements (market_id);
