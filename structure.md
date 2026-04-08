# Structure

## Purpose

I am scoping the repository around the backend first and intentionally leaving the frontend out for now.

This file maps the current project layout, shows where each backend concern belongs, and links to the README for every folder currently present in the repo.

## Current Folder Structure

```text
Pacific Prediction/
+-- .agents/
¦   +-- GUIDE.md
¦   +-- README.md
+-- backend/
¦   +-- cmd/
¦   ¦   +-- README.md
¦   ¦   +-- api/
¦   ¦       +-- README.md
¦   +-- internal/
¦   ¦   +-- README.md
¦   ¦   +-- auth/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- balance/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- config/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- domain/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- httpapi/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- market/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- pacifica/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- position/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- player/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- realtime/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- settlement/
¦   ¦   ¦   +-- README.md
¦   ¦   +-- storage/
¦   ¦       +-- README.md
¦   +-- migrations/
¦   ¦   +-- README.md
¦   +-- README.md
+-- Readme.md
+-- Resources.md
+-- Research.md
+-- architecture.md
+-- schema.md
+-- task.md
+-- structure.md
```

## High-Level Mapping

- Product and sponsor context live in [Readme.md](file:///C:/Hackathons/Pacific%20Prediction/Readme.md) and [Resources.md](file:///C:/Hackathons/Pacific%20Prediction/Resources.md).
- Research findings and API-driven constraints live in [Research.md](file:///C:/Hackathons/Pacific%20Prediction/Research.md).
- Confirmed architecture decisions live in [architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
- Schema decisions and extension rules live in [schema.md](file:///C:/Hackathons/Pacific%20Prediction/schema.md).
- Execution-ready implementation sequencing lives in [task.md](file:///C:/Hackathons/Pacific%20Prediction/task.md).
- Persistent project constraints live in [.agents/README.md](file:///C:/Hackathons/Pacific%20Prediction/.agents/README.md) and [.agents/GUIDE.md](file:///C:/Hackathons/Pacific%20Prediction/.agents/GUIDE.md).
- Backend planning and module ownership begin in [backend/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/README.md).

## Backend Logic Map

- To find backend entrypoint planning visit [backend/cmd/api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- To find auth and session logic visit [backend/internal/auth/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/auth/README.md).
- To find balance logic visit [backend/internal/balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- To find configuration logic visit [backend/internal/config/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/README.md).
- To find shared domain vocabulary visit [backend/internal/domain/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md).
- To find HTTP transport logic visit [backend/internal/httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).
- To find market lifecycle logic visit [backend/internal/market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find Pacifica integration logic visit [backend/internal/pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find position lifecycle logic visit [backend/internal/position/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/README.md).
- To find player identity logic visit [backend/internal/player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- To find realtime delivery logic visit [backend/internal/realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- To find settlement logic visit [backend/internal/settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find persistence and repository logic visit [backend/internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- To find database schema and migration planning visit [backend/migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md).

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
- [backend/internal/position/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/README.md)
- [backend/internal/player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md)
- [backend/internal/realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md)
- [backend/internal/settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md)
- [backend/internal/storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md)
- [backend/migrations/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/migrations/README.md)

## Backend Boundary Summary

I am keeping the backend split simple:

- `cmd` will hold executable entrypoints only.
- `internal/auth` will own guest sessions and future auth upgrades.
- `internal/balance` will own authoritative virtual balance rules.
- `internal/config` will own environment and runtime configuration.
- `internal/domain` will own shared enum-like values, ID aliases, UTC timestamp rules, and common validation errors.
- `internal/httpapi` will own HTTP handlers, request validation, and response shaping.
- `internal/market` will own market creation, listing, and resolution rules at the domain level.
- `internal/pacifica` will isolate all Pacifica REST and WebSocket integration.
- `internal/position` will own YES or NO placement contracts and player position history rules.
- `internal/player` will own player identity and profile rules.
- `internal/realtime` will push live updates from backend state to the client-facing stream.
- `internal/settlement` will own expiry scanning and deterministic resolution.
- `internal/storage` will isolate PostgreSQL persistence and transaction boundaries.
- `migrations` will hold schema evolution files once implementation starts.

## Tradeoffs

- I am documenting a deeper backend structure now because it will prevent a monolithic Go service later.
- I am not creating frontend folders yet because the current planning focus is backend-first.
- I am keeping the first split coarse enough to stay practical, but separated enough that Pacifica integration, settlement, auth, and balance logic do not collapse into one package.


