# Frontend/src/routes/api/[...path]

## Purpose

I use this folder for the concrete catch-all proxy handler that forwards frontend API requests to the Go backend.

## Architectural Decisions And Tradeoffs

- I keep the actual proxy implementation isolated here so the route boundary is explicit.
- This keeps browser-facing transport code out of feature pages, even though it adds one more folder in the route tree.

## Logic Tracking

- To find the catch-all proxy implementation visit [+server.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/[...path]/+server.ts).
- To find backend URL construction visit [../../../lib/server/backend-proxy.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/backend-proxy.ts).

## Component And Connection Map

- The frontend proxy handler can be found in [+server.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/[...path]/+server.ts).
- The backend URL builder can be found in [../../../lib/server/backend-proxy.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/backend-proxy.ts).
