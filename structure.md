# Structure

## Purpose

I use this file to map the current repository layout, show where the implemented backend logic lives, and link to the README for every backend folder that already has ownership notes.

The frontend folder exists, but the backend is still the most complete part of the repo right now, so this map is intentionally backend-heavy.

## Current Folder Structure

```text
Pacific Prediction/
+-- .agents/
|   +-- GUIDE.md
|   +-- README.md
+-- backend/
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
|   +-- README.md
+-- Frontend/
+-- Readme.md
+-- Resources.md
+-- Research.md
+-- architecture.md
+-- frontend-notes.md
+-- schema.md
+-- structure.md
+-- task.md
```

## High-Level Mapping

- Product-level overview and current repo entrypoint live in [Readme.md](file:///C:/Hackathons/Pacific%20Prediction/Readme.md).
- Architecture decisions and tradeoffs live in [architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
- Schema reasoning and extension boundaries live in [schema.md](file:///C:/Hackathons/Pacific%20Prediction/schema.md).
- Execution sequencing and task state live in [task.md](file:///C:/Hackathons/Pacific%20Prediction/task.md).
- Frontend flow notes for creator auto-stake live in [frontend-notes.md](file:///C:/Hackathons/Pacific%20Prediction/frontend-notes.md).
- Pacifica research notes live in [Research.md](file:///C:/Hackathons/Pacific%20Prediction/Research.md).
- Backend package ownership begins in [backend/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/README.md).

## Backend Logic Map

- To find backend entrypoint wiring visit [backend/cmd/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- To find auth and session logic visit [backend/internal/auth/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/auth/README.md).
- To find balance rules and settlement-time balance accounting visit [backend/internal/balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- To find configuration loading visit [backend/internal/config/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/README.md).
- To find shared domain vocabulary and decimal helpers visit [backend/internal/domain/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md).
- To find HTTP request and response shaping visit [backend/internal/httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- To find market creation and creator auto-stake orchestration visit [backend/internal/market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find Pacifica REST and WebSocket integration visit [backend/internal/pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find direct participant position placement visit [backend/internal/position/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/README.md).
- To find player profile logic visit [backend/internal/player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- To find realtime placeholder ownership visit [backend/internal/realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- To find deterministic settlement and payout application visit [backend/internal/settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find PostgreSQL repositories and transaction boundaries visit [backend/internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- To find migration files visit [backend/migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).

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
- [backend/internal/market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md)
- [backend/internal/pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md)
- [backend/internal/player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md)
- [backend/internal/position/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/README.md)
- [backend/internal/realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md)
- [backend/internal/settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md)
- [backend/internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md)
- [backend/migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md)

## Backend Boundary Summary

I am keeping the backend split simple and practical:

- `cmd` owns executable startup only.
- `auth` owns guest session lifecycle.
- `balance` owns authoritative balance concepts and mutation intent.
- `config` owns environment-backed runtime config.
- `domain` owns shared enums, IDs, time helpers, candle interval helpers, and decimal helpers.
- `httpapi` owns transport shaping only.
- `market` owns market creation and creator auto-stake orchestration.
- `pacifica` isolates upstream market data access.
- `position` owns participant position placement and history reads.
- `player` owns player identity reads.
- `settlement` owns resolver orchestration and transactional payout completion.
- `storage` owns PostgreSQL repositories and transaction boundaries.
- `migrations` owns schema history.
