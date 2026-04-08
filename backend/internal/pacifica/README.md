# backend/internal/pacifica

## Purpose

I use this folder to isolate every Pacifica REST and WebSocket integration behind one backend-owned boundary.

## Architectural Decisions And Tradeoffs

- I do not want Pacifica-specific request shapes leaking into domain packages.
- This package should own subscription management, heartbeat handling, reconnect behavior, and settlement lookups.
- Read-only REST fetches for metadata, prices, mark-price candles, and funding history should all live here so settlement and validation code do not depend on raw vendor payloads.
- Market metadata should be fetched here and cached aggressively so symbol validation can stay dynamic without making market creation depend on a fresh Pacifica call every time.
- The tradeoff is a translation boundary, but it protects the rest of the app from vendor-specific transport details.

## Logic Tracking

- To find Pacifica integration logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find unified Pacifica REST client contracts visit [rest_client.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/rest_client.go).
- To find Pacifica WebSocket client contracts visit [websocket_client.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/websocket_client.go).
- To find Pacifica subscription manager contracts visit [subscription_manager.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/subscription_manager.go).
- To find market metadata client contracts visit [market_info.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/market_info.go).
- To find normalized Pacifica price models visit [prices.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/prices.go).
- To find normalized Pacifica mark-price candle models visit [candles.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/candles.go).
- To find normalized Pacifica funding-history models visit [funding.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/funding.go).
- To find normalized Pacifica live event models visit [live_events.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/live_events.go).
- To find the current HTTP REST client implementation visit [market_info_http_client.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/market_info_http_client.go).
- To find reconnect and heartbeat orchestration visit [subscription_manager_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/subscription_manager_impl.go).
- To find the market module that consumes this metadata for validation visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find settlement dependency mapping visit [../settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find backend-wide architecture decisions visit [../../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).

## Component And Connection Map

- The Pacifica REST and WebSocket connection can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- The Pacifica REST contract boundary can be found in [rest_client.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/rest_client.go).
- The Pacifica WebSocket contract boundary can be found in [websocket_client.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/websocket_client.go).
- The Pacifica subscription manager boundary can be found in [subscription_manager.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/subscription_manager.go).
- The current read-only HTTP REST fetch path can be found in [market_info_http_client.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/market_info_http_client.go).
- The normalized settlement data models can be found in [prices.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/prices.go), [candles.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/candles.go), and [funding.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/funding.go).
- The normalized live update models can be found in [live_events.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/live_events.go).
- The reconnect and heartbeat implementation can be found in [subscription_manager_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/subscription_manager_impl.go).
- The client-facing realtime handoff can be found in [../realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
