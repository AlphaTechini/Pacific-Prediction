# Backend

## Purpose

I am structuring the backend as a single Go service that owns guest identity, virtual balances, prediction markets, participant positions, deterministic settlement, and Pacifica read-only integrations.

## Architectural Decisions And Tradeoffs

- I chose Go because the difficult parts of v1 are long-lived data ingestion, deterministic settlement, and transaction safety.
- I am keeping Pacifica integration read-only in v1 so the service stays focused on game logic instead of wallet or order execution concerns.
- I am using PostgreSQL without Redis because correctness matters more than speculative optimization at this stage.
- I am keeping market creation product-shaped so one request can create the market and the creator's first staked position together.

## Current Backend Capabilities

- Guest session creation and secure-cookie session lookup
- Current player profile and balance reads
- Product-shaped market creation with creator side and creator stake
- Additional YES or NO position placement on existing markets
- Active and resolved market listing plus market detail reads
- Price, candle, and funding settlement from Pacifica-backed sources
- Transactional payout application that updates positions, balances, and settlement audit records together

## Logic Tracking

- To find backend system boundaries visit [architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
- To find the backend module map visit [structure.md](file:///C:/Hackathons/Pacific%20Prediction/structure.md).
- To find backend entrypoint planning visit [cmd/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- To find domain and integration ownership visit [internal/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/README.md).
- To find frontend flow notes that affect backend contracts visit [../frontend-notes.md](file:///C:/Hackathons/Pacific%20Prediction/frontend-notes.md).

## Component And Connection Map

- The backend service boundary can be found in [architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
- The backend entrypoint planning can be found in [cmd/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/README.md).
- The backend domain packages can be found in [internal/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/README.md).
- The database migration planning can be found in [migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
