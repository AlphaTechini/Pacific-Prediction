# backend/internal/realtime

## Purpose

I use this folder for pushing live backend state to the client-facing stream.

## Architectural Decisions And Tradeoffs

- I prefer SSE first because the frontend mainly needs outbound updates, not bidirectional socket control.
- This package should expose backend-owned stream contracts for market lifecycle and settlement updates instead of reusing raw Pacifica live models directly.
- I am using a typed stream event envelope plus explicit market and settlement snapshots so the future SSE layer can serialize stable product events without relying on `map[string]any`.
- The tradeoff is less long-term flexibility than a full app WebSocket layer, but it is the right level of complexity for v1.
- The tradeoff of typed snapshots is a little extra mapping work, but it keeps Pacifica-specific churn out of the frontend stream contract.

## Logic Tracking

- To find realtime event type definitions visit [event_type.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/event_type.go).
- To find the shared realtime stream event envelope visit [event.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/event.go).
- To find market stream snapshot contracts visit [market_snapshot.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/market_snapshot.go).
- To find settlement stream snapshot contracts visit [settlement_snapshot.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/settlement_snapshot.go).
- To find Pacifica live data sourcing visit [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find HTTP stream exposure visit [../httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).

## Component And Connection Map

- The client-facing event envelope can be found in [event.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/event.go).
- The market lifecycle stream payload can be found in [market_snapshot.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/market_snapshot.go).
- The settlement stream payload can be found in [settlement_snapshot.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/settlement_snapshot.go).
- The upstream live data connection can be found in [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
