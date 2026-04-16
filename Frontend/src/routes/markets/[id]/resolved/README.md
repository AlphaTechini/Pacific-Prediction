# Frontend/src/routes/markets/[id]/resolved

## Purpose

I use this folder for the resolved market view.

## Architectural Decisions And Tradeoffs

- I keep the resolved screen read-only and tied to the market id route hierarchy.
- This keeps the result surface simple and avoids mixing settlement display concerns into the main trading state.
- The page still ensures the cached guest session so the resolved read can show player-context data consistently.

## Logic Tracking

- To find the resolved market page visit [+page.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/[id]/resolved/+page.svelte).
- To find market-detail and result data helpers visit [../../../../lib/market-detail-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/market-detail-data.ts).
- To find guest-session flow helpers and the localStorage guest cache visit [../../../../lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).

## Component And Connection Map

- The resolved-market route-to-backend connection can be found in [../../../../lib/market-detail-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/market-detail-data.ts).
- The resolved-market guest-session dependency can be found in [../../../../lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).
