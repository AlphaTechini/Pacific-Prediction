# backend/internal

## Purpose

I use `internal` to hold the actual application logic and prevent external packages from depending on backend-only code.

## Architectural Decisions And Tradeoffs

- I split by responsibility so auth, balances, Pacifica integration, market rules, storage, and settlement can evolve independently.
- I prefer domain-oriented package boundaries here instead of one giant `service` package.
- The tradeoff is more folders, but it keeps the backend from collapsing into a monolith as features are added.

## Logic Tracking

- To find auth and session logic visit [auth/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/auth/README.md).
- To find balance logic visit [balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- To find configuration logic visit [config/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/README.md).
- To find shared domain vocabulary visit [domain/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md).
- To find HTTP transport logic visit [httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- To find market lifecycle logic visit [market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find Pacifica integration logic visit [pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find player identity logic visit [player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- To find realtime delivery logic visit [realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- To find settlement logic visit [settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find storage logic visit [storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).

## Component And Connection Map

- The auth system boundary can be found in [auth/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/auth/README.md).
- The balance system boundary can be found in [balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- The shared backend domain vocabulary can be found in [domain/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md).
- The Pacifica system connection can be found in [pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- The PostgreSQL connection boundary can be found in [storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
