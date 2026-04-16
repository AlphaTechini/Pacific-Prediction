# Frontend

## Purpose

I use this frontend to present the real Pacifica Pulse product flow: guest entry, live market browsing, market creation, market detail, portfolio, and leaderboard.

## Architectural Decisions And Tradeoffs

- I keep the frontend talking to our backend only, never directly to Pacifica.
- I use a SvelteKit proxy route under `/api/*` so cookies and backend URL management stay centralized.
- I keep most business logic in the backend and let the frontend focus on rendering and user input.
- I use a mix of route-level and client-side loading today because it keeps the current app simple without blocking progress.
- I cache the guest player id and display name in localStorage so revisiting the app does not silently create another guest identity.
- I serve the PNG favicon from `static` and reuse it as the top-navigation brand mark instead of keeping the old Svelte SVG.
- I accept that the landing page still contains older concept copy while the product pages are already much closer to backend truth.

## Current Frontend Capabilities

- Landing page shell at `/`
- Guest-session-aware dashboard at `/dashboard`
- Real create-market flow at `/markets/create`
- Real market detail view at `/markets/[id]`
- Resolved market view at `/markets/[id]/resolved`
- Portfolio view at `/portfolio`
- Server-loaded leaderboard at `/leaderboard`
- Backend proxy transport at `/api/[...path]`
- Static PNG brand/favicon asset at `/favicon.png`

## Logic Tracking

- To find the repo-level structure map visit [../structure.md](file:///C:/Hackathons/Pacific%20Prediction/structure.md).
- To find platform architecture decisions visit [../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
- To find frontend notes and remaining alignment work visit [../frontend-notes.md](file:///C:/Hackathons/Pacific%20Prediction/frontend-notes.md) and [../frontend-integration.md](file:///C:/Hackathons/Pacific%20Prediction/frontend-integration.md).
- To find frontend source ownership visit [src/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/README.md).
- To find static public assets visit [static/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/static/README.md).
- To find mockup and stitch references visit [stitch_screens/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/stitch_screens/README.md).

## Component And Connection Map

- The frontend application boundary can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/README.md).
- The frontend source tree can be found in [src/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/README.md).
- The route ownership map can be found in [src/routes/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/README.md).
- The shared frontend data and UI layer can be found in [src/lib/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/README.md).
- The backend proxy connection can be found in [src/lib/server/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/README.md).
- The static brand asset boundary can be found in [static/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/static/README.md).
