CREATE INDEX IF NOT EXISTS idx_positions_resolved_leaderboard
    ON positions (player_id, settled_at DESC, created_at DESC, id DESC)
    WHERE status IN ('won', 'lost');

CREATE INDEX IF NOT EXISTS idx_positions_market_id_player_id
    ON positions (market_id, player_id);
