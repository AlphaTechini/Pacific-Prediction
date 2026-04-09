# Frontend/src/routes

## Purpose

I use this folder for SvelteKit route ownership, page-level UI, and the app-side backend proxy route.

## Architectural Decisions And Tradeoffs

- I keep route folders aligned with product surfaces so each page owns its own behavior.
- I keep the backend proxy route here because it is part of the app's transport boundary.
- I accept a mix of client-side and route-level loading for now because it keeps the current implementation practical.

## Logic Tracking

- To find the landing page visit [+page.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/+page.svelte).
- To find the shared app layout visit [+layout.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/+layout.svelte) and [layout.css](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/layout.css).
- To find the backend proxy route visit [api/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/README.md).
- To find the dashboard route visit [dashboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/dashboard/README.md).
- To find the leaderboard route visit [leaderboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/leaderboard/README.md).
- To find the market routes visit [markets/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/README.md).
- To find the portfolio route visit [portfolio/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/portfolio/README.md).

## Component And Connection Map

- The browser-to-backend proxy connection can be found in [api/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/README.md).
- The app route ownership map can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/README.md).
