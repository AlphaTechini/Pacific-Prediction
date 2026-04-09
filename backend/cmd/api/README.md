# backend/cmd/api

## Purpose

I reserve this folder for the main HTTP API process entrypoint.

## Architectural Decisions And Tradeoffs

- This package should wire config, database connections, Pacifica clients, HTTP routes, and background workers.
- The public leaderboard snapshot route should also be wired here because this folder owns API composition.
- I do not want route handlers, SQL, or settlement rules in this folder.
- I keep `main.go` as the composition root and use the shared `httpapi.Application` and `httpapi.Router` types to wire dependencies and register method-aware routes without letting modules self-register.
- The tradeoff is more package wiring, but startup remains explicit and testable.

## Logic Tracking

- To find API process startup planning visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- To find API startup wiring visit [main.go](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/main.go).
- To find HTTP handler ownership visit [../../internal/httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- To find leaderboard route ownership visit [../../internal/leaderboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/README.md).
- To find realtime stream wiring visit [../../internal/realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- To find settlement worker ownership visit [../../internal/settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).

## Component And Connection Map

- The main API process can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- The runtime worker wiring can be found in [main.go](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/main.go).
- The HTTP transport can be found in [../../internal/httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- The leaderboard ranking package can be found in [../../internal/leaderboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/README.md).
- The realtime stream package can be found in [../../internal/realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- The Pacifica connection layer can be found in [../../internal/pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
