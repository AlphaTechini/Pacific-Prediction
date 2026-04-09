# backend/internal/httpapi

## Purpose

I use this folder for HTTP routing, request validation, response shaping, and transport-only concerns.

## Architectural Decisions And Tradeoffs

- I keep handlers thin so business rules stay in domain packages.
- I shape responses here because the frontend should consume product-shaped payloads, not raw Pacifica responses.
- I keep route ownership in `main.go`, but I use a shared `Application` container and `Router` helper here so dependency wiring and HTTP method checks stay consistent across modules.
- I let auth middleware resolve the session cookie into request context here so domain packages do not depend on HTTP cookie parsing.
- Market creation uses the same thin-handler pattern, with request parsing and RFC3339 coercion in HTTP and domain validation/persistence in the market module.
- Market creation should stay product-shaped in HTTP, which now means the route accepts creator side and stake so the frontend can submit one request instead of coordinating a second opening-position call.
- Market listing returns grouped active and resolved catalogs for the default dashboard read path.
- Market detail and position placement use Go's method-aware and path-aware `http.ServeMux` patterns so nested market routes stay explicit without adding a third-party router.
- The tradeoff is an extra translation layer, but it gives me a cleaner contract and safer future changes.

## Logic Tracking

- To find HTTP transport logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- To find market domain logic visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find position domain logic visit [../position/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/README.md).
- To find player and balance logic visit [../player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- To find the API process that owns route registration visit [../../cmd/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).

## Component And Connection Map

- The client-facing API transport can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- The realtime client stream connection can be found in [../realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- The composition root that wires controllers into this transport can be found in [../../cmd/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
