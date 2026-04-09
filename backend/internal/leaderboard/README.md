# backend/internal/leaderboard

## Purpose

I use this folder for leaderboard read models and ranking orchestration.

## Architectural Decisions And Tradeoffs

- I keep leaderboard logic in its own package because these rankings are product features, not generic player or balance reads.
- I derive the leaderboard from markets and positions instead of creating a dedicated leaderboard table right now.
- I fan out the category reads in parallel from Go so one slow category does not serialize the whole snapshot.
- The tradeoff is a few aggregate queries per request, but the system stays simple and transparent while current scale is modest.

## Logic Tracking

- To find leaderboard snapshot contracts visit [contracts.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/contracts.go).
- To find leaderboard service contracts visit [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/service.go).
- To find leaderboard service orchestration visit [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/service_impl.go).
- To find leaderboard controller contracts visit [controller.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/controller.go).
- To find leaderboard controller implementation visit [controller_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/controller_impl.go).
- To find leaderboard SQL reads visit [../storage/leaderboard_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/leaderboard_repository.go) and [../storage/leaderboard_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/leaderboard_postgres_repository.go).
- To find the HTTP route that exposes this package visit [../httpapi/leaderboard_handlers.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/leaderboard_handlers.go).

## Component And Connection Map

- The leaderboard domain boundary can be found in [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/service.go) and [controller.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/controller.go).
- The PostgreSQL aggregation connection can be found in [../storage/leaderboard_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/leaderboard_postgres_repository.go).
- The public leaderboard HTTP connection can be found in [../httpapi/leaderboard_handlers.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/leaderboard_handlers.go).
