# backend/internal/storage

## Purpose

I use this folder for PostgreSQL repositories, transaction boundaries, and persistence abstractions.

## Architectural Decisions And Tradeoffs

- I want SQL and transaction control isolated from HTTP handlers and domain coordination code.
- This package should own persistence for players, markets, positions, balances, and settlement records.
- This package should also own the aggregate SQL used for the leaderboard read model.
- Balance persistence should support stake locking at entry time and settlement-time clearing or payout application without leaking SQL into orchestration code.
- I split repository contracts by module concern and expose them through a transaction-scoped repository provider so domain services can coordinate writes safely without sharing raw SQL handles.
- The tradeoff is explicit mapping code and a little more wiring, but that keeps the database layer predictable and replaceable.

## Logic Tracking

- To find persistence and repository logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- To find player persistence contracts visit [player_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/player_repository.go).
- To find player PostgreSQL persistence visit [player_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/player_postgres_repository.go).
- To find session persistence contracts visit [session_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/session_repository.go).
- To find session PostgreSQL persistence visit [session_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/session_postgres_repository.go).
- To find balance persistence contracts visit [balance_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/balance_repository.go).
- To find balance PostgreSQL persistence visit [balance_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/balance_postgres_repository.go).
- To find market persistence contracts visit [market_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/market_repository.go).
- To find market PostgreSQL persistence visit [market_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/market_postgres_repository.go).
- To find position persistence contracts visit [position_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/position_repository.go).
- To find position PostgreSQL persistence visit [position_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/position_postgres_repository.go).
- To find settlement persistence contracts visit [settlement_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/settlement_repository.go).
- To find settlement PostgreSQL persistence visit [settlement_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/settlement_postgres_repository.go).
- To find leaderboard persistence contracts visit [leaderboard_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/leaderboard_repository.go).
- To find leaderboard PostgreSQL aggregation reads visit [leaderboard_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/leaderboard_postgres_repository.go).
- To find database connection setup visit [db.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/db.go).
- To find migration runner logic visit [migrator.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/migrator.go).
- To find transaction helpers visit [tx.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/tx.go) and [transaction_contracts.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/transaction_contracts.go).
- To find the shared domain vocabulary used by repository contracts visit [../domain/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md).
- To find schema evolution planning visit [../../migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
- To find the domain packages that depend on storage visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md), [../player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md), and [../position/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/README.md).

## Component And Connection Map

- The PostgreSQL persistence boundary can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- The transaction-safe repository contract pattern can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- The position persistence connection can be found in [position_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/position_repository.go) and [position_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/position_postgres_repository.go).
- The settlement audit persistence connection can be found in [settlement_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/settlement_repository.go) and [settlement_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/settlement_postgres_repository.go).
- The leaderboard aggregation connection can be found in [leaderboard_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/leaderboard_repository.go) and [leaderboard_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/leaderboard_postgres_repository.go).
- The database schema planning can be found in [../../migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
