# backend/internal/httpapi

## Purpose

I use this folder for HTTP routing, request validation, response shaping, and transport-only concerns.

## Architectural Decisions And Tradeoffs

- I keep handlers thin so business rules stay in domain packages.
- I shape responses here because the frontend should consume product-shaped payloads, not raw Pacifica responses.
- The tradeoff is an extra translation layer, but it gives me a cleaner contract and safer future changes.

## Logic Tracking

- To find HTTP transport logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- To find market domain logic visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find player and balance logic visit [../player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).

## Component And Connection Map

- The client-facing API transport can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- The realtime client stream connection can be found in [../realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
