# backend/internal/pacifica

## Purpose

I use this folder to isolate every Pacifica REST and WebSocket integration behind one backend-owned boundary.

## Architectural Decisions And Tradeoffs

- I do not want Pacifica-specific request shapes leaking into domain packages.
- This package should own subscription management, heartbeat handling, reconnect behavior, and settlement lookups.
- The tradeoff is a translation boundary, but it protects the rest of the app from vendor-specific transport details.

## Logic Tracking

- To find Pacifica integration logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find settlement dependency mapping visit [../settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find backend-wide architecture decisions visit [../../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).

## Component And Connection Map

- The Pacifica REST and WebSocket connection can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- The client-facing realtime handoff can be found in [../realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
