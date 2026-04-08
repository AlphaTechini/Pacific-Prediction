# Backend

## Purpose

I am structuring the backend as a single Go service that owns prediction markets, guest users, virtual balances, settlement, and Pacifica read-only integrations.

## Architectural Decisions And Tradeoffs

- I chose Go because the difficult parts of v1 are long-lived data ingestion, deterministic settlement, and transaction safety.
- I am keeping Pacifica integration read-only in v1 so the service stays focused on game logic instead of wallet or order execution concerns.
- I am using PostgreSQL without Redis because correctness matters more than speculative optimization at this stage.

## Logic Tracking

- To find backend system boundaries visit [architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
- To find the backend module map visit [structure.md](file:///C:/Hackathons/Pacific%20Prediction/structure.md).
- To find backend entrypoint planning visit [cmd/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- To find domain and integration ownership visit [internal/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/README.md).

## Component And Connection Map

- The backend service boundary can be found in [architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
- The backend entrypoint planning can be found in [cmd/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/README.md).
- The backend domain packages can be found in [internal/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/README.md).
- The database migration planning can be found in [migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).
