# backend/migrations

## Purpose

I use this folder for database schema migrations and seed strategy once backend implementation starts.

## Architectural Decisions And Tradeoffs

- I keep migrations separate from repositories because schema history should be auditable on its own.
- This folder will hold table creation and later schema changes for players, balances, markets, positions, and settlements.
- The tradeoff is another artifact stream to maintain, but it gives me safe schema evolution instead of ad hoc SQL drift.

## Logic Tracking

- To find database schema planning visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
- To find the storage layer that uses this schema visit [../internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- To find the architecture-level data model visit [../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).

## Component And Connection Map

- The database schema evolution path can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
- The PostgreSQL repository connection can be found in [../internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
