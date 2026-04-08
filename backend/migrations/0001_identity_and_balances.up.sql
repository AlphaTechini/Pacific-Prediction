CREATE TABLE IF NOT EXISTS players (
    id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS player_sessions (
    id TEXT PRIMARY KEY,
    player_id TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    session_token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_player_sessions_token_hash
    ON player_sessions (session_token_hash);

CREATE TABLE IF NOT EXISTS player_balances (
    player_id TEXT PRIMARY KEY REFERENCES players(id) ON DELETE CASCADE,
    available_balance NUMERIC(20, 8) NOT NULL DEFAULT 0,
    locked_balance NUMERIC(20, 8) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_available_balance_nonnegative CHECK (available_balance >= 0),
    CONSTRAINT chk_locked_balance_nonnegative CHECK (locked_balance >= 0)
);
