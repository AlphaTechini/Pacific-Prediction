# backend/internal/realtime

## Purpose

I use this folder for pushing live backend state to the client-facing stream.

## Architectural Decisions And Tradeoffs

- I prefer SSE first because the frontend mainly needs outbound updates, not bidirectional socket control.
- This package should translate backend events into stream-safe payloads for active markets, countdowns, and resolutions.
- The tradeoff is less long-term flexibility than a full app WebSocket layer, but it is the right level of complexity for v1.

## Logic Tracking

- To find realtime delivery logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- To find Pacifica live data sourcing visit [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find HTTP stream exposure visit [../httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).

## Component And Connection Map

- The client-facing event stream can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- The upstream live data connection can be found in [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
