# Frontend/src/routes/markets/create

## Purpose

I use this folder for the create-market page.

## Architectural Decisions And Tradeoffs

- I keep market creation as one guided flow that includes creator side and creator stake.
- The page still initializes after mount because it first ensures a guest session through the cached guest helper and then loads backend context.
- I generate the visible market question from the selected symbol, rule, and timing so the saved title stays aligned with settlement logic instead of relying on freeform copy.
- I mirror the backend precision rules in the UI by preferring Pacifica `tick_size` for price-threshold guidance and only falling back to a positive `min_tick` when needed.

## Logic Tracking

- To find the create-market page visit [+page.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/create/+page.svelte).
- To find create-market data helpers visit [../../../lib/create-market-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/create-market-data.ts).
- To find guest-session flow helpers and the localStorage guest cache visit [../../../lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).

## Component And Connection Map

- The create-market route-to-backend context connection can be found in [../../../lib/create-market-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/create-market-data.ts).
