# backend/internal/pacifica

## Purpose

I use this folder to isolate every Pacifica REST and WebSocket integration behind one backend-owned boundary.

## Architectural Decisions And Tradeoffs

- I do not want Pacifica-specific request shapes leaking into domain packages.
- This package should own subscription management, heartbeat handling, reconnect behavior, and settlement lookups.
- Market metadata should be fetched here and cached aggressively so symbol validation can stay dynamic without making market creation depend on a fresh Pacifica call every time.
- The tradeoff is a translation boundary, but it protects the rest of the app from vendor-specific transport details.

## Logic Tracking

- To find Pacifica integration logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find market metadata client contracts visit [market_info.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/market_info.go).
- To find the cached REST market metadata client visit [market_info_http_client.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/market_info_http_client.go).
- To find the market module that consumes this metadata for validation visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find settlement dependency mapping visit [../settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find backend-wide architecture decisions visit [../../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).

## Component And Connection Map

- The Pacifica REST and WebSocket connection can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- The cached market metadata fetch path can be found in [market_info_http_client.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/market_info_http_client.go).
- The client-facing realtime handoff can be found in [../realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
