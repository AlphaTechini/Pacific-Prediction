# Backend

## Purpose

I use the backend as one Go service that owns guest identity, virtual balances, prediction markets, participant positions, leaderboard reads, deterministic settlement, and Pacifica read-only integrations.

## Architectural Decisions And Tradeoffs

- I chose Go because the difficult parts of v1 are long-lived settlement work, transaction safety, and predictable orchestration.
- I keep Pacifica integration read-only in v1 so the service stays focused on game logic instead of wallet or order execution concerns.
- I use PostgreSQL without Redis because current scale does not justify another runtime dependency.
- I keep market creation product-shaped so one request can create the market and the creator's first staked position together.
- I keep the leaderboard as a derived Postgres read model instead of a separate cache or write-time table.
- I allow local `.env` loading in the config layer so `go run` works from normal backend working directories without manual environment exporting.
- I keep the Postgres pool explicitly configurable so the service can hold warm connections longer and avoid unnecessary reconnect churn against Supabase session pooling.

## Current Backend Capabilities

- Guest session creation and secure-cookie session lookup
- Cookie-backed authorization remains server-authoritative; the frontend guest cache is not trusted for backend access
- Current player profile and balance reads
- Product-shaped market creation with creator side and creator stake
- Additional YES or NO position placement on existing markets
- Active and resolved market listing plus market detail reads
- Market-create context reads for the frontend form
- Public leaderboard snapshot reads
- Public SSE stream reads
- Price, candle, and funding settlement from Pacifica-backed sources
- Transactional payout application that updates positions, balances, and settlement audit records together

## Current Public API Surface

- `POST /api/v1/players/guest`
- `GET /api/v1/players/me`
- `GET /api/v1/players/me/balance`
- `GET /api/v1/players/me/positions`
- `GET /api/v1/leaderboard`
- `GET /api/v1/stream`
- `POST /api/v1/markets`
- `GET /api/v1/markets`
- `GET /api/v1/markets/context`
- `GET /api/v1/markets/{market_id}`
- `POST /api/v1/markets/{market_id}/positions`

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
- The leaderboard package can be found in [internal/leaderboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/README.md).
- The database migration planning can be found in [migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
