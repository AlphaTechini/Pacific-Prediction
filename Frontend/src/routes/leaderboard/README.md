# Frontend/src/routes/leaderboard

## Purpose

I use this folder for the leaderboard page and its route-level data load.

## Architectural Decisions And Tradeoffs

- I load the leaderboard at the route level because one backend snapshot is cheaper and simpler than several client-side requests.
- I keep ranking logic on the backend so this route stays mostly presentational.

## Logic Tracking

- To find the leaderboard page load visit [+page.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/leaderboard/+page.ts).
- To find the leaderboard page UI visit [+page.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/leaderboard/+page.svelte).
- To find leaderboard data helpers visit [../../lib/leaderboard-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/leaderboard-data.ts).

## Component And Connection Map

- The leaderboard route-to-backend snapshot connection can be found in [../../lib/leaderboard-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/leaderboard-data.ts).
