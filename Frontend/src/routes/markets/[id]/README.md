# Frontend/src/routes/markets/[id]

## Purpose

I use this folder for the live market-detail route keyed by market id.

## Architectural Decisions And Tradeoffs

- I keep this page focused on one market plus the current player's ability to take a position.
- The page reads live market data from the backend instead of caching its own source-of-truth state locally.
- The page ensures the cached guest session before loading player-specific balance and position state.

## Logic Tracking

- To find the market-detail page visit [+page.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/[id]/+page.svelte).
- To find market-detail data helpers visit [../../../lib/market-detail-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/market-detail-data.ts).
- To find guest-session flow helpers and the localStorage guest cache visit [../../../lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).

## Component And Connection Map

- The market-detail route-to-backend connection can be found in [../../../lib/market-detail-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/market-detail-data.ts).
- The market-detail guest-session dependency can be found in [../../../lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).
