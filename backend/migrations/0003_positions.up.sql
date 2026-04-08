CREATE TABLE IF NOT EXISTS positions (
    id TEXT PRIMARY KEY,
    player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    market_id TEXT NOT NULL REFERENCES markets(id) ON DELETE CASCADE,
    side TEXT NOT NULL,
    stake_amount NUMERIC(20, 8) NOT NULL,
    potential_payout NUMERIC(20, 8) NOT NULL,
    status TEXT NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    settled_at TIMESTAMPTZ,
    CONSTRAINT chk_positions_side CHECK (side IN ('yes', 'no')),
    CONSTRAINT chk_positions_stake_amount_positive CHECK (stake_amount > 0),
    CONSTRAINT chk_positions_potential_payout_positive CHECK (potential_payout > 0),
    CONSTRAINT chk_positions_status CHECK (status IN ('open', 'won', 'lost', 'cancelled')),
    CONSTRAINT chk_positions_settlement_state CHECK (
        (status = 'open' AND settled_at IS NULL)
        OR (status IN ('won', 'lost', 'cancelled') AND settled_at IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_positions_player_id_created_at
    ON positions (player_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_positions_market_id_created_at
    ON positions (market_id, created_at DESC);
