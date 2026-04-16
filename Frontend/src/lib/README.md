# Frontend/src/lib

## Purpose

I use this folder for shared frontend data helpers, API types, reusable components, and backend-proxy utilities.

## Architectural Decisions And Tradeoffs

- I keep API-shaping helpers here so route files stay focused on page behavior.
- I keep shared UI components here so the app shell stays visually consistent.
- I keep the backend proxy utilities isolated under `server` because they are not browser-only code.
- I let `guest-session.ts` cache the guest player id and display name in browser localStorage so returning to the app does not silently create another guest identity.
- I keep that cache as a UI continuity hint only; backend authorization still depends on the secure session cookie.

## Logic Tracking

- To find shared API types visit [api-types.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/api-types.ts).
- To find backend fetch helpers visit [backend-api.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/backend-api.ts).
- To find guest-session flow helpers visit [guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).
- To find dashboard data helpers visit [dashboard-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/dashboard-data.ts).
- To find create-market data helpers visit [create-market-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/create-market-data.ts).
- To find market-detail data helpers visit [market-detail-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/market-detail-data.ts).
- To find portfolio data helpers visit [portfolio-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/portfolio-data.ts).
- To find leaderboard data helpers visit [leaderboard-data.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/leaderboard-data.ts).
- To find shared UI components visit [components/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/README.md).
- To find backend-proxy utilities visit [server/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/README.md).
- To find assets and mockups visit [assets/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/assets/README.md).

## Component And Connection Map

- The shared frontend type system can be found in [api-types.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/api-types.ts).
- The backend session connection and localStorage guest cache can be found in [guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts).
- The backend proxy connection can be found in [server/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/README.md).
