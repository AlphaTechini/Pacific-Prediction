# backend/migrations

## Purpose

I use this folder for database schema migrations that back the live backend implementation.

## Architectural Decisions And Tradeoffs

- I keep migrations separate from repositories because schema history should be auditable on its own.
- This folder will hold table creation and later schema changes for players, balances, markets, positions, and settlements.
- Read-optimized indexes for leaderboard queries should also live here so ranking performance stays explicit and reversible.
- The tradeoff is another artifact stream to maintain, but it gives me safe schema evolution instead of ad hoc SQL drift.

## Logic Tracking

- To find database schema planning visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
- To find identity and balance schema creation visit [0001_identity_and_balances.up.sql](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/0001_identity_and_balances.up.sql).
- To find market schema creation visit [0002_markets.up.sql](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/0002_markets.up.sql).
- To find position schema creation visit [0003_positions.up.sql](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/0003_positions.up.sql).
- To find settlement audit schema creation visit [0004_market_settlements.up.sql](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/0004_market_settlements.up.sql).
- To find leaderboard index tuning visit [0005_leaderboard_indexes.up.sql](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/0005_leaderboard_indexes.up.sql).
- To find the storage layer that uses this schema visit [../internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- To find the architecture-level data model visit [../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).

## Component And Connection Map

- The database schema evolution path can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
- The position table definition can be found in [0003_positions.up.sql](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/0003_positions.up.sql).
- The settlement audit table definition can be found in [0004_market_settlements.up.sql](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/0004_market_settlements.up.sql).
- The leaderboard index path can be found in [0005_leaderboard_indexes.up.sql](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/0005_leaderboard_indexes.up.sql).
- The PostgreSQL repository connection can be found in [../internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
