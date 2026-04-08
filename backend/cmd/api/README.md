# backend/cmd/api

## Purpose

I reserve this folder for the main HTTP API process entrypoint.

## Architectural Decisions And Tradeoffs

- This package should wire config, database connections, Pacifica clients, HTTP routes, and background workers.
- I do not want route handlers, SQL, or settlement rules in this folder.
- The tradeoff is more package wiring, but startup remains explicit and testable.

## Logic Tracking

- To find API process startup planning visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- To find HTTP handler ownership visit [../../internal/httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- To find settlement worker ownership visit [../../internal/settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).

## Component And Connection Map

- The main API process can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- The HTTP transport can be found in [../../internal/httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- The Pacifica connection layer can be found in [../../internal/pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
