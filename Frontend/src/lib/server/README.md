# Frontend/src/lib/server

## Purpose

I use this folder for server-only helpers that let SvelteKit reach the Go backend cleanly.

## Architectural Decisions And Tradeoffs

- I keep backend URL construction here so route handlers and loads do not duplicate environment logic.
- I prefer one backend URL helper over ad hoc string concatenation across routes.

## Logic Tracking

- To find backend proxy URL construction visit [backend-proxy.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/backend-proxy.ts).

## Component And Connection Map

- The SvelteKit-to-backend connection can be found in [backend-proxy.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/backend-proxy.ts).
