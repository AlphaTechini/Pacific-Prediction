# Frontend/src/routes/portfolio

## Purpose

I use this folder for the player portfolio page.

## Architectural Decisions And Tradeoffs

- I keep portfolio focused on balance and position history instead of speculative analytics.
- The page derives its display from backend state and does not try to own balance math locally.

## Logic Tracking

- To find the portfolio page visit [+page.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/portfolio/+page.svelte).
- To find portfolio data helpers visit [../../lib/portfolio-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/portfolio-data.ts).

## Component And Connection Map

- The portfolio route-to-backend connection can be found in [../../lib/portfolio-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/portfolio-data.ts).
