# Frontend/src/routes/api

## Purpose

I use this folder for the SvelteKit catch-all backend proxy route.

## Architectural Decisions And Tradeoffs

- I proxy browser `/api/*` traffic through SvelteKit so cookies and backend URL handling stay centralized.
- This adds one hop, but it keeps the frontend boundary much cleaner.

## Logic Tracking

- To find the catch-all backend proxy handler visit [catch-all proxy +server.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/[...path]/+server.ts).

## Component And Connection Map

- The browser-to-Go-backend transport can be found in [catch-all proxy +server.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/[...path]/+server.ts).
- The backend URL builder can be found in [../../lib/server/backend-proxy.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/backend-proxy.ts).
