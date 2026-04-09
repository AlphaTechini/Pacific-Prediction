# backend/internal/realtime

## Purpose

I use this folder for pushing live backend state to the client-facing stream.

## Architectural Decisions And Tradeoffs

- I prefer SSE first because the frontend mainly needs outbound updates, not bidirectional socket control.
- This package should expose backend-owned stream contracts for market lifecycle and settlement updates instead of reusing raw Pacifica live models directly.
- I am using a typed stream event envelope plus explicit market and settlement snapshots so the future SSE layer can serialize stable product events without relying on `map[string]any`.
- The first transport path is a public SSE subscription endpoint with heartbeat keepalives, because market and settlement updates are public dashboard data in v1.
- The tradeoff is less long-term flexibility than a full app WebSocket layer, but it is the right level of complexity for v1.
- The tradeoff of typed snapshots is a little extra mapping work, but it keeps Pacifica-specific churn out of the frontend stream contract.
- The transport starts with a no-op subscription service and heartbeat loop so `T7.2` can establish the stream boundary cleanly before `T7.3` adds the in-process event hub.

## Logic Tracking

- To find realtime event type definitions visit [event_type.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/event_type.go).
- To find the shared realtime stream event envelope visit [event.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/event.go).
- To find market stream snapshot contracts visit [market_snapshot.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/market_snapshot.go).
- To find settlement stream snapshot contracts visit [settlement_snapshot.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/settlement_snapshot.go).
- To find realtime subscription contracts visit [subscription.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/subscription.go).
- To find realtime service contracts visit [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/service.go).
- To find realtime service behavior visit [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/service_impl.go).
- To find realtime controller contracts visit [controller.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/controller.go).
- To find realtime controller behavior visit [controller_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/controller_impl.go).
- To find Pacifica live data sourcing visit [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find HTTP stream exposure visit [../httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).

## Component And Connection Map

- The client-facing event envelope can be found in [event.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/event.go).
- The client subscription lifecycle can be found in [subscription.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/subscription.go), [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/service.go), and [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/service_impl.go).
- The market lifecycle stream payload can be found in [market_snapshot.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/market_snapshot.go).
- The settlement stream payload can be found in [settlement_snapshot.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/settlement_snapshot.go).
- The upstream live data connection can be found in [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
