# Frontend/src/routes/dashboard

## Purpose

I use this folder for the live dashboard page.

## Architectural Decisions And Tradeoffs

- I keep this page focused on active and resolved markets plus the current player's state.
- The page still fetches after mount instead of using a route load, and the guest-session helper now checks the cached guest player before provisioning anything new.

## Logic Tracking

- To find the dashboard page visit [+page.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/dashboard/+page.svelte).
- To find dashboard data loading visit [../../lib/dashboard-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/dashboard-data.ts).
- To find guest-session bootstrapping and the localStorage guest cache visit [../../lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).

## Component And Connection Map

- The dashboard-to-backend read path can be found in [../../lib/dashboard-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/dashboard-data.ts).
