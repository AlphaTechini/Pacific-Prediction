CREATE TABLE IF NOT EXISTS markets (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    symbol TEXT NOT NULL,
    market_type TEXT NOT NULL,
    condition_operator TEXT NOT NULL,
    threshold_value NUMERIC(20, 8),
    source_type TEXT NOT NULL,
    source_interval TEXT,
    reference_value NUMERIC(20, 8),
    expiry_time TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    result TEXT,
    settlement_value NUMERIC(20, 8),
    resolved_at TIMESTAMPTZ,
    resolution_reason TEXT,
    created_by_player_id TEXT NOT NULL REFERENCES players(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_markets_market_type CHECK (market_type IN ('price_threshold', 'candle_direction', 'funding_threshold')),
    CONSTRAINT chk_markets_condition_operator CHECK (condition_operator IN ('gt', 'gte', 'lt', 'lte', 'bullish_close', 'bearish_close', 'positive', 'negative')),
    CONSTRAINT chk_markets_source_type CHECK (source_type IN ('mark_price', 'mark_price_candle', 'funding_rate')),
    CONSTRAINT chk_markets_status CHECK (status IN ('active', 'resolving', 'resolved', 'cancelled')),
    CONSTRAINT chk_markets_result CHECK (result IS NULL OR result IN ('yes', 'no')),
    CONSTRAINT chk_markets_price_threshold_shape CHECK (
        market_type <> 'price_threshold'
        OR (
            source_type = 'mark_price'
            AND source_interval IS NULL
            AND threshold_value IS NOT NULL
            AND condition_operator IN ('gt', 'gte', 'lt', 'lte')
        )
    ),
    CONSTRAINT chk_markets_candle_direction_shape CHECK (
        market_type <> 'candle_direction'
        OR (
            source_type = 'mark_price_candle'
            AND source_interval IS NOT NULL
            AND threshold_value IS NULL
            AND condition_operator IN ('bullish_close', 'bearish_close')
        )
    ),
    CONSTRAINT chk_markets_funding_threshold_shape CHECK (
        market_type <> 'funding_threshold'
        OR (
            source_type = 'funding_rate'
            AND source_interval = 'funding_epoch'
            AND (
                (condition_operator IN ('gt', 'gte', 'lt', 'lte') AND threshold_value IS NOT NULL)
                OR (condition_operator IN ('positive', 'negative') AND threshold_value IS NULL)
            )
        )
    ),
    CONSTRAINT chk_markets_resolution_state CHECK (
        (status IN ('active', 'resolving') AND result IS NULL)
        OR (status IN ('resolved', 'cancelled'))
    ),
    CONSTRAINT chk_markets_resolved_columns CHECK (
        status NOT IN ('resolved', 'cancelled')
        OR resolved_at IS NOT NULL
    )
);

CREATE INDEX IF NOT EXISTS idx_markets_status_expiry_time
    ON markets (status, expiry_time);

CREATE INDEX IF NOT EXISTS idx_markets_created_by_player_id_created_at
    ON markets (created_by_player_id, created_at DESC);
