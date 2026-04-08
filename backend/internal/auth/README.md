# backend/internal/auth

## Purpose

I use this folder for guest session creation, session validation, session lookup, and future authentication upgrades.

## Architectural Decisions And Tradeoffs

- I keep auth separate from player identity so guest sessions can evolve into wallet auth later without polluting the player module.
- This package should own session issuance, validation, revocation, and auth-context creation.
- The tradeoff is one more module, but it gives me a clean upgrade path for wallet auth and external identity providers.

## Logic Tracking

- To find auth and session logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/auth/README.md).
- To find player identity logic visit [../player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- To find the main route-wiring rule visit [../../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).

## Component And Connection Map

- The guest session system can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/auth/README.md).
- The persistent player identity it authenticates can be found in [../player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
