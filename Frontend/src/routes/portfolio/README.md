# Frontend/src/routes/portfolio

## Purpose

I use this folder for the player portfolio page.

## Architectural Decisions And Tradeoffs

- I keep portfolio focused on balance and position history instead of speculative analytics.
- The page derives its display from backend state and does not try to own balance math locally.
- The page ensures the cached guest session before loading player-specific portfolio data.

## Logic Tracking

- To find the portfolio page visit [+page.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/portfolio/+page.svelte).
- To find portfolio data helpers visit [../../lib/portfolio-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/portfolio-data.ts).
- To find guest-session flow helpers and the localStorage guest cache visit [../../lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).

## Component And Connection Map

- The portfolio route-to-backend connection can be found in [../../lib/portfolio-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/portfolio-data.ts).
- The portfolio guest-session dependency can be found in [../../lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).
