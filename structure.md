# Structure

## Purpose

I use this file to map the current repository layout, point to the folder ownership docs, and show where the implemented frontend and backend logic now live.

This map is no longer backend-heavy. The repo now has a real application loop across both stacks, so the structure summary reflects both sides.

## Current Folder Structure

```text
Pacific Prediction/
+-- .agents/
|   +-- GUIDE.md
|   +-- README.md
+-- backend/
|   +-- README.md
|   +-- cmd/
|   |   +-- README.md
|   |   +-- api/
|   |       +-- README.md
|   +-- internal/
|   |   +-- README.md
|   |   +-- auth/
|   |   |   +-- README.md
|   |   +-- balance/
|   |   |   +-- README.md
|   |   +-- config/
|   |   |   +-- README.md
|   |   +-- domain/
|   |   |   +-- README.md
|   |   +-- httpapi/
|   |   |   +-- README.md
|   |   +-- leaderboard/
|   |   |   +-- README.md
|   |   +-- market/
|   |   |   +-- README.md
|   |   +-- pacifica/
|   |   |   +-- README.md
|   |   +-- player/
|   |   |   +-- README.md
|   |   +-- position/
|   |   |   +-- README.md
|   |   +-- realtime/
|   |   |   +-- README.md
|   |   +-- settlement/
|   |   |   +-- README.md
|   |   +-- storage/
|   |       +-- README.md
|   +-- migrations/
|   |   +-- README.md
+-- Frontend/
|   +-- README.md
|   +-- static/
|   |   +-- README.md
|   +-- stitch_screens/
|   |   +-- README.md
|   +-- src/
|       +-- README.md
|       +-- lib/
|       |   +-- README.md
|       |   +-- assets/
|       |   |   +-- README.md
|       |   |   +-- mockups/
|       |   |       +-- README.md
|       |   +-- components/
|       |   |   +-- README.md
|       |   +-- server/
|       |       +-- README.md
|       +-- routes/
|           +-- README.md
|           +-- api/
|           |   +-- README.md
|           |   +-- [...path]/
|           |       +-- README.md
|           +-- dashboard/
|           |   +-- README.md
|           +-- leaderboard/
|           |   +-- README.md
|           +-- markets/
|           |   +-- README.md
|           |   +-- create/
|           |   |   +-- README.md
|           |   +-- [id]/
|           |       +-- README.md
|           |       +-- resolved/
|           |           +-- README.md
|           +-- portfolio/
|               +-- README.md
+-- skills/
|   +-- README.md
|   +-- pacifica-pulse-v1/
|       +-- README.md
|       +-- SKILL.md
+-- Readme.md
+-- Resources.md
+-- Research.md
+-- architecture.md
+-- frontend-integration.md
+-- frontend-notes.md
+-- schema.md
+-- structure.md
+-- task.md
```

## High-Level Mapping

- Product-level overview and repo entrypoint live in [Readme.md](file:///C:/Hackathons/Pacific%20Prediction/Readme.md).
- Architecture decisions and tradeoffs live in [architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
- Schema reasoning and extension boundaries live in [schema.md](file:///C:/Hackathons/Pacific%20Prediction/schema.md).
- Project memory and confirmed constraints live in [.agents/GUIDE.md](file:///C:/Hackathons/Pacific%20Prediction/.agents/GUIDE.md).
- Frontend implementation notes live in [frontend-notes.md](file:///C:/Hackathons/Pacific%20Prediction/frontend-notes.md) and [frontend-integration.md](file:///C:/Hackathons/Pacific%20Prediction/frontend-integration.md).
- Backend package ownership begins in [backend/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/README.md).
- Frontend package ownership begins in [Frontend/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/README.md).
- Public OpenClaw skill packaging begins in [skills/README.md](file:///C:/Hackathons/Pacific%20Prediction/skills/README.md).

## Frontend Logic Map

- To find frontend source ownership visit [Frontend/src/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/README.md).
- To find shared frontend data, UI, and types visit [Frontend/src/lib/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/README.md).
- To find reusable UI components visit [Frontend/src/lib/components/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/README.md).
- To find backend proxy helpers visit [Frontend/src/lib/server/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/README.md).
- To find route ownership visit [Frontend/src/routes/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/README.md).
- To find the app-side API proxy visit [Frontend/src/routes/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/README.md) and [Frontend/src/routes/api catch-all README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/[...path]/README.md).
- To find the live dashboard route visit [Frontend/src/routes/dashboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/dashboard/README.md).
- To find the leaderboard route visit [Frontend/src/routes/leaderboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/leaderboard/README.md).
- To find market route ownership visit [Frontend/src/routes/markets/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/README.md).
- To find the create-market route visit [Frontend/src/routes/markets/create/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/create/README.md).
- To find the market-detail route visit [Frontend/src/routes/markets detail README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/[id]/README.md).
- To find the resolved-market route visit [Frontend/src/routes/resolved market README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/[id]/resolved/README.md).
- To find the portfolio route visit [Frontend/src/routes/portfolio/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/portfolio/README.md).

## Backend Logic Map

- To find backend entrypoint wiring visit [backend/cmd/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- To find auth and session logic visit [backend/internal/auth/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/auth/README.md).
- To find balance rules and settlement-time balance accounting visit [backend/internal/balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- To find configuration loading visit [backend/internal/config/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/README.md).
- To find shared domain vocabulary and decimal helpers visit [backend/internal/domain/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md).
- To find HTTP request and response shaping visit [backend/internal/httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- To find leaderboard ranking logic visit [backend/internal/leaderboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/README.md).
- To find market creation and creator auto-stake orchestration visit [backend/internal/market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find Pacifica REST and WebSocket integration visit [backend/internal/pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find direct participant position placement visit [backend/internal/position/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/README.md).
- To find player profile logic visit [backend/internal/player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- To find realtime stream ownership visit [backend/internal/realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- To find deterministic settlement and payout application visit [backend/internal/settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find PostgreSQL repositories and transaction boundaries visit [backend/internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- To find migration files visit [backend/migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).

## Skills Logic Map

- To find the public OpenClaw skill index visit [skills/README.md](file:///C:/Hackathons/Pacific%20Prediction/skills/README.md).
- To find the publishable Pacifica Pulse skill visit [skills/pacifica-pulse-v1/README.md](file:///C:/Hackathons/Pacific%20Prediction/skills/pacifica-pulse-v1/README.md).
- To find the actual OpenClaw skill entrypoint visit [skills/pacifica-pulse-v1/SKILL.md](file:///C:/Hackathons/Pacific%20Prediction/skills/pacifica-pulse-v1/SKILL.md).

## Folder README Index

- [.agents/README.md](file:///C:/Hackathons/Pacific%20Prediction/.agents/README.md)
- [backend/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/README.md)
- [backend/cmd/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/README.md)
- [backend/cmd/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md)
- [backend/internal/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/README.md)
- [backend/internal/auth/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/auth/README.md)
- [backend/internal/balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md)
- [backend/internal/config/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/README.md)
- [backend/internal/domain/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md)
- [backend/internal/httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md)
- [backend/internal/leaderboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/leaderboard/README.md)
- [backend/internal/market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md)
- [backend/internal/pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md)
- [backend/internal/player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md)
- [backend/internal/position/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/README.md)
- [backend/internal/realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md)
- [backend/internal/settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md)
- [backend/internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md)
- [backend/migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md)
- [Frontend/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/README.md)
- [Frontend/static/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/static/README.md)
- [Frontend/stitch_screens/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/stitch_screens/README.md)
- [Frontend/src/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/README.md)
- [Frontend/src/lib/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/README.md)
- [Frontend/src/lib/assets/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/assets/README.md)
- [Frontend/src/lib/assets/mockups/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/assets/mockups/README.md)
- [Frontend/src/lib/components/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/README.md)
- [Frontend/src/lib/server/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/server/README.md)
- [Frontend/src/routes/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/README.md)
- [Frontend/src/routes/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/README.md)
- [Frontend/src/routes/api catch-all README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/[...path]/README.md)
- [Frontend/src/routes/dashboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/dashboard/README.md)
- [Frontend/src/routes/leaderboard/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/leaderboard/README.md)
- [Frontend/src/routes/markets/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/README.md)
- [Frontend/src/routes/markets/create/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/create/README.md)
- [Frontend/src/routes/markets detail README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/[id]/README.md)
- [Frontend/src/routes/resolved market README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/[id]/resolved/README.md)
- [Frontend/src/routes/portfolio/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/portfolio/README.md)
- [skills/README.md](file:///C:/Hackathons/Pacific%20Prediction/skills/README.md)
- [skills/pacifica-pulse-v1/README.md](file:///C:/Hackathons/Pacific%20Prediction/skills/pacifica-pulse-v1/README.md)

## Boundary Summary

I am keeping the platform split practical:

- the root docs define the product boundary and architecture
- `Frontend` owns app routes, shared UI, and the browser-to-backend proxy layer
- `backend` owns sessions, balances, markets, positions, settlement, realtime, and leaderboard reads
- `.agents/GUIDE.md` keeps the confirmed constraints and architectural decisions synchronized across future work
