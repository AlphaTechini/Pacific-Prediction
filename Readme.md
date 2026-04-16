# Pacifica Pulse

## Documentation

- [structure.md](./structure.md)
- [architecture.md](./architecture.md)
- [schema.md](./schema.md)
- [Research.md](./Research.md)
- [frontend-notes.md](./frontend-notes.md)
- [frontend-integration.md](./frontend-integration.md)

## Overview

I am building Pacifica Pulse as a read-only prediction game on top of Pacifica testnet market data.

The current codebase is no longer backend-only. The real product loop now exists across both stacks:

- SvelteKit pages for dashboard, market creation, market detail, resolved markets, portfolio, and leaderboard
- a Go API for guest sessions, balances, markets, positions, leaderboard reads, and SSE updates
- PostgreSQL-backed authoritative balances, positions, and settlement records
- a settlement worker that resolves price, candle, and funding markets from Pacifica data
- a SvelteKit static PNG favicon that also acts as the top-navigation brand mark

## Current Product Shape

In v1 today:

- I use guest accounts instead of wallet auth.
- I cache the guest player id and display name in browser localStorage so returning users do not silently create another guest identity on revisit.
- I use virtual balances and virtual payouts only.
- I support `price_threshold`, `candle_direction`, and `funding_threshold` markets.
- I settle every market from one explicit Pacifica-derived source.
- I expose a real leaderboard derived from stored player, market, and position data.
- I do not include AI summaries, real trading, withdrawals, or Pacifica write operations.

## Current Frontend Shape

The frontend currently includes:

- a landing page that still carries some older concept copy
- a real dashboard backed by markets, balances, and positions
- a create-market flow backed by `GET /api/v1/markets/context` and `POST /api/v1/markets`
- real market detail and resolved-market views
- a portfolio page backed by player balance and positions
- a leaderboard page backed by `GET /api/v1/leaderboard`
- a SvelteKit backend proxy under `/api/*` so browser calls stay pointed at our app, not directly at Pacifica
- a PNG brand icon served from `Frontend/static/favicon.png` and reused by the top navigation

## Current Backend API

The backend currently provides:

- `POST /api/v1/players/guest`
- `GET /api/v1/players/me`
- `GET /api/v1/players/me/balance`
- `GET /api/v1/players/me/positions`
- `GET /api/v1/leaderboard`
- `GET /api/v1/stream`
- `POST /api/v1/markets`
- `GET /api/v1/markets`
- `GET /api/v1/markets/context`
- `GET /api/v1/markets/{market_id}`
- `POST /api/v1/markets/{market_id}/positions`

## Important Behavior Already Implemented

- Market creation includes the creator's required opening side and stake.
- `POST /api/v1/markets` creates the market and the creator's first position in one transaction.
- Additional players can join through the position-placement route.
- Price markets settle from Pacifica mark-price snapshots.
- Candle markets settle from Pacifica mark-price candles.
- Funding markets settle from Pacifica funding history.
- Settlement updates market result, position outcomes, balances, and settlement audit records in one transaction.
- The leaderboard is derived from stored market and position history instead of a separate cache system.
- Guest continuity is assisted by a frontend localStorage cache, but backend authorization still depends on the secure session cookie.

## Payout Model

I am using a simple fixed-odds house rule for v1:

- every position locks the chosen stake at entry time
- losers receive no further balance credit at settlement
- winners receive `2x stake`

This keeps the product easy to explain and easy to verify.

## Known Implementation Gaps

- The landing page still contains aspirational copy and unsupported signal language from the earlier concept phase.
- The frontend still uses `adapter-auto` today even though the project standard calls for `@sveltejs/adapter-vercel`.
- Most app pages already use real backend data, but the dashboard and create-market flow still fetch after mount rather than loading everything through route-level server loads.

## Where To Look

- To find the current repository map visit [structure.md](./structure.md).
- To find the architecture decisions and tradeoffs visit [architecture.md](./architecture.md).
- To find schema reasoning visit [schema.md](./schema.md).
- To find frontend implementation notes visit [frontend-notes.md](./frontend-notes.md) and [frontend-integration.md](./frontend-integration.md).
- To find backend package ownership visit [backend/README.md](./backend/README.md).
- To find frontend package ownership visit [Frontend/README.md](./Frontend/README.md).
