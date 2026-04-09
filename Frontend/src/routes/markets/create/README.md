# Frontend/src/routes/markets/create

## Purpose

I use this folder for the create-market page.

## Architectural Decisions And Tradeoffs

- I keep market creation as one guided flow that includes creator side and creator stake.
- The page still initializes after mount because it first ensures a guest session and then loads backend context.

## Logic Tracking

- To find the create-market page visit [+page.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/create/+page.svelte).
- To find create-market data helpers visit [../../../lib/create-market-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/create-market-data.ts).
- To find guest-session flow helpers visit [../../../lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).

## Component And Connection Map

- The create-market route-to-backend context connection can be found in [../../../lib/create-market-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/create-market-data.ts).
