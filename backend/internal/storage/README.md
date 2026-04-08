# backend/internal/storage

## Purpose

I use this folder for PostgreSQL repositories, transaction boundaries, and persistence abstractions.

## Architectural Decisions And Tradeoffs

- I want SQL and transaction control isolated from HTTP handlers and domain coordination code.
- This package should own persistence for players, markets, positions, balances, and settlement records.
- I split repository contracts by module concern and expose them through a transaction-scoped repository provider so domain services can coordinate writes safely without sharing raw SQL handles.
- The tradeoff is explicit mapping code and a little more wiring, but that keeps the database layer predictable and replaceable.

## Logic Tracking

- To find persistence and repository logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- To find the shared domain vocabulary used by repository contracts visit [../domain/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md).
- To find schema evolution planning visit [../../migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
- To find the domain packages that depend on storage visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md) and [../player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).

## Component And Connection Map

- The PostgreSQL persistence boundary can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- The transaction-safe repository contract pattern can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- The database schema planning can be found in [../../migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
