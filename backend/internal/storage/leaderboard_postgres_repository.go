package storage

import (
	"context"
	"fmt"
)

type LeaderboardPostgresRepository struct {
	queryer Queryer
}

func NewLeaderboardPostgresRepository(queryer Queryer) *LeaderboardPostgresRepository {
	return &LeaderboardPostgresRepository{queryer: queryer}
}

func (r *LeaderboardPostgresRepository) GetOverview(ctx context.Context) (LeaderboardOverview, error) {
	const query = `
WITH player_win_rates AS (
    SELECT
        ROUND(
            100.0 * COUNT(*) FILTER (WHERE status = 'won')::numeric / NULLIF(COUNT(*), 0),
            2
        ) AS win_rate
    FROM positions
    WHERE status IN ('won', 'lost')
    GROUP BY player_id
)
SELECT
    (SELECT COUNT(*) FROM positions) AS total_predictions,
    (SELECT COUNT(*) FROM positions WHERE status IN ('won', 'lost')) AS resolved_predictions,
    (SELECT COUNT(DISTINCT player_id) FROM positions) AS active_predictors,
    (SELECT COUNT(DISTINCT created_by_player_id) FROM markets) AS active_creators,
    COALESCE((SELECT ROUND(AVG(win_rate), 2)::text FROM player_win_rates), '0.00') AS average_win_rate;
`

	var overview LeaderboardOverview
	if err := r.queryer.QueryRow(ctx, query).Scan(
		&overview.TotalPredictions,
		&overview.ResolvedPredictions,
		&overview.ActivePredictors,
		&overview.ActiveCreators,
		&overview.AverageWinRate,
	); err != nil {
		return LeaderboardOverview{}, fmt.Errorf("get leaderboard overview: %w", err)
	}

	return overview, nil
}

func (r *LeaderboardPostgresRepository) ListTopPredictors(ctx context.Context, limit int) ([]LeaderboardPredictor, error) {
	const query = `
WITH predictor_stats AS (
    SELECT
        p.player_id,
        COUNT(*) AS resolved_positions,
        COUNT(*) FILTER (WHERE p.status = 'won') AS won_positions,
        COUNT(*) FILTER (WHERE p.status = 'lost') AS lost_positions,
        COALESCE(SUM(p.stake_amount), 0)::text AS total_staked,
        COALESCE(
            SUM(
                CASE
                    WHEN p.status = 'won' THEN p.potential_payout - p.stake_amount
                    ELSE -p.stake_amount
                END
            ),
            0
        )::text AS net_profit,
        ROUND(
            100.0 * COUNT(*) FILTER (WHERE p.status = 'won')::numeric / NULLIF(COUNT(*), 0),
            2
        )::text AS win_rate
    FROM positions p
    WHERE p.status IN ('won', 'lost')
    GROUP BY p.player_id
)
SELECT
    ps.player_id,
    pl.display_name,
    ps.resolved_positions,
    ps.won_positions,
    ps.lost_positions,
    ps.win_rate,
    ps.net_profit,
    ps.total_staked
FROM predictor_stats ps
JOIN players pl ON pl.id = ps.player_id
ORDER BY
    ps.net_profit::numeric DESC,
    ps.win_rate::numeric DESC,
    ps.resolved_positions DESC,
    pl.display_name ASC
LIMIT $1;
`

	rows, err := r.queryer.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("list top predictors: %w", err)
	}
	defer rows.Close()

	var items []LeaderboardPredictor
	for rows.Next() {
		var item LeaderboardPredictor
		if err := rows.Scan(
			&item.PlayerID,
			&item.DisplayName,
			&item.ResolvedPositions,
			&item.WonPositions,
			&item.LostPositions,
			&item.WinRate,
			&item.NetProfit,
			&item.TotalStaked,
		); err != nil {
			return nil, fmt.Errorf("list top predictors: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list top predictors: %w", err)
	}

	return items, nil
}

func (r *LeaderboardPostgresRepository) ListTopCreators(ctx context.Context, limit int) ([]LeaderboardCreator, error) {
	const query = `
WITH creator_market_stats AS (
    SELECT
        m.created_by_player_id AS player_id,
        COUNT(*) AS created_markets,
        COUNT(*) FILTER (WHERE m.status = 'resolved') AS resolved_markets
    FROM markets m
    GROUP BY m.created_by_player_id
),
creator_engagement AS (
    SELECT
        m.created_by_player_id AS player_id,
        COUNT(p.id) AS total_positions,
        COUNT(DISTINCT p.player_id) AS unique_participants,
        COALESCE(SUM(p.stake_amount), 0)::text AS total_staked_on_markets
    FROM markets m
    LEFT JOIN positions p ON p.market_id = m.id
    GROUP BY m.created_by_player_id
)
SELECT
    cms.player_id,
    pl.display_name,
    cms.created_markets,
    cms.resolved_markets,
    COALESCE(ce.total_positions, 0) AS total_positions,
    COALESCE(ce.unique_participants, 0) AS unique_participants,
    COALESCE(ce.total_staked_on_markets, '0') AS total_staked_on_markets
FROM creator_market_stats cms
JOIN players pl ON pl.id = cms.player_id
LEFT JOIN creator_engagement ce ON ce.player_id = cms.player_id
ORDER BY
    cms.created_markets DESC,
    COALESCE(ce.unique_participants, 0) DESC,
    COALESCE(ce.total_positions, 0) DESC,
    pl.display_name ASC
LIMIT $1;
`

	rows, err := r.queryer.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("list top creators: %w", err)
	}
	defer rows.Close()

	var items []LeaderboardCreator
	for rows.Next() {
		var item LeaderboardCreator
		if err := rows.Scan(
			&item.PlayerID,
			&item.DisplayName,
			&item.CreatedMarkets,
			&item.ResolvedMarkets,
			&item.TotalPositions,
			&item.UniqueParticipants,
			&item.TotalStakedOnMarkets,
		); err != nil {
			return nil, fmt.Errorf("list top creators: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list top creators: %w", err)
	}

	return items, nil
}

func (r *LeaderboardPostgresRepository) ListBestStreaks(ctx context.Context, limit int) ([]LeaderboardStreak, error) {
	const query = `
WITH resolved_positions AS (
    SELECT
        p.player_id,
        p.status,
        p.stake_amount,
        p.potential_payout,
        p.settled_at,
        p.created_at,
        p.id
    FROM positions p
    WHERE p.status IN ('won', 'lost')
),
resolved_stats AS (
    SELECT
        rp.player_id,
        COUNT(*) AS resolved_positions,
        ROUND(
            100.0 * COUNT(*) FILTER (WHERE rp.status = 'won')::numeric / NULLIF(COUNT(*), 0),
            2
        )::text AS win_rate,
        COALESCE(
            SUM(
                CASE
                    WHEN rp.status = 'won' THEN rp.potential_payout - rp.stake_amount
                    ELSE -rp.stake_amount
                END
            ),
            0
        )::text AS net_profit
    FROM resolved_positions rp
    GROUP BY rp.player_id
),
ordered_forward AS (
    SELECT
        rp.player_id,
        rp.status,
        rp.settled_at,
        ROW_NUMBER() OVER (
            PARTITION BY rp.player_id
            ORDER BY rp.settled_at ASC, rp.created_at ASC, rp.id ASC
        ) AS overall_row,
        ROW_NUMBER() OVER (
            PARTITION BY rp.player_id, rp.status
            ORDER BY rp.settled_at ASC, rp.created_at ASC, rp.id ASC
        ) AS status_row
    FROM resolved_positions rp
),
win_runs AS (
    SELECT
        player_id,
        COUNT(*) AS streak_length
    FROM ordered_forward
    WHERE status = 'won'
    GROUP BY player_id, overall_row - status_row
),
longest_streaks AS (
    SELECT
        player_id,
        MAX(streak_length) AS longest_win_streak
    FROM win_runs
    GROUP BY player_id
),
ordered_reverse AS (
    SELECT
        rp.player_id,
        rp.status,
        ROW_NUMBER() OVER (
            PARTITION BY rp.player_id
            ORDER BY rp.settled_at DESC, rp.created_at DESC, rp.id DESC
        ) AS overall_row,
        ROW_NUMBER() OVER (
            PARTITION BY rp.player_id, rp.status
            ORDER BY rp.settled_at DESC, rp.created_at DESC, rp.id DESC
        ) AS status_row
    FROM resolved_positions rp
),
current_streaks AS (
    SELECT
        player_id,
        COUNT(*) AS current_win_streak
    FROM ordered_reverse
    WHERE status = 'won'
      AND overall_row - status_row = 0
    GROUP BY player_id
)
SELECT
    rs.player_id,
    pl.display_name,
    COALESCE(cs.current_win_streak, 0) AS current_win_streak,
    COALESCE(ls.longest_win_streak, 0) AS longest_win_streak,
    rs.resolved_positions,
    rs.win_rate,
    rs.net_profit
FROM resolved_stats rs
JOIN players pl ON pl.id = rs.player_id
LEFT JOIN current_streaks cs ON cs.player_id = rs.player_id
LEFT JOIN longest_streaks ls ON ls.player_id = rs.player_id
ORDER BY
    COALESCE(cs.current_win_streak, 0) DESC,
    COALESCE(ls.longest_win_streak, 0) DESC,
    rs.win_rate::numeric DESC,
    rs.resolved_positions DESC,
    pl.display_name ASC
LIMIT $1;
`

	rows, err := r.queryer.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("list best streaks: %w", err)
	}
	defer rows.Close()

	var items []LeaderboardStreak
	for rows.Next() {
		var item LeaderboardStreak
		if err := rows.Scan(
			&item.PlayerID,
			&item.DisplayName,
			&item.CurrentWinStreak,
			&item.LongestWinStreak,
			&item.ResolvedPositions,
			&item.WinRate,
			&item.NetProfit,
		); err != nil {
			return nil, fmt.Errorf("list best streaks: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list best streaks: %w", err)
	}

	return items, nil
}

func (r *LeaderboardPostgresRepository) ListMostActive(ctx context.Context, limit int) ([]LeaderboardActivity, error) {
	const query = `
WITH position_stats AS (
    SELECT
        p.player_id,
        COUNT(*) AS total_positions,
        COUNT(*) FILTER (WHERE p.status = 'open') AS open_positions,
        COUNT(*) FILTER (WHERE p.status IN ('won', 'lost')) AS resolved_positions,
        COALESCE(SUM(p.stake_amount), 0)::text AS total_staked
    FROM positions p
    GROUP BY p.player_id
),
market_stats AS (
    SELECT
        m.created_by_player_id AS player_id,
        COUNT(*) AS created_markets
    FROM markets m
    GROUP BY m.created_by_player_id
)
SELECT
    pl.id,
    pl.display_name,
    COALESCE(ps.total_positions, 0) AS total_positions,
    COALESCE(ps.open_positions, 0) AS open_positions,
    COALESCE(ps.resolved_positions, 0) AS resolved_positions,
    COALESCE(ms.created_markets, 0) AS created_markets,
    COALESCE(ps.total_staked, '0') AS total_staked
FROM players pl
LEFT JOIN position_stats ps ON ps.player_id = pl.id
LEFT JOIN market_stats ms ON ms.player_id = pl.id
WHERE COALESCE(ps.total_positions, 0) > 0
   OR COALESCE(ms.created_markets, 0) > 0
ORDER BY
    COALESCE(ps.total_positions, 0) DESC,
    COALESCE(ms.created_markets, 0) DESC,
    COALESCE(ps.resolved_positions, 0) DESC,
    pl.display_name ASC
LIMIT $1;
`

	rows, err := r.queryer.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("list most active players: %w", err)
	}
	defer rows.Close()

	var items []LeaderboardActivity
	for rows.Next() {
		var item LeaderboardActivity
		if err := rows.Scan(
			&item.PlayerID,
			&item.DisplayName,
			&item.TotalPositions,
			&item.OpenPositions,
			&item.ResolvedPositions,
			&item.CreatedMarkets,
			&item.TotalStaked,
		); err != nil {
			return nil, fmt.Errorf("list most active players: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list most active players: %w", err)
	}

	return items, nil
}

var _ LeaderboardRepository = (*LeaderboardPostgresRepository)(nil)
