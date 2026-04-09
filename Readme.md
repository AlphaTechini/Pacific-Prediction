# Pacifica Pulse

## Documentation

- [structure.md](./structure.md)
- [architecture.md](./architecture.md)
- [schema.md](./schema.md)
- [task.md](./task.md)
- [Research.md](./Research.md)
- [frontend-notes.md](./frontend-notes.md)

## Overview

I am building Pacifica Pulse as a read-only prediction game on top of Pacifica testnet market data.

The current implementation focus is the backend. It already supports guest sessions, authoritative virtual balances, market creation, position placement, and deterministic settlement for the v1 market types.

## Current Product Shape

In v1:

- I use guest accounts instead of wallet auth.
- I use virtual balances and virtual payouts only.
- I support `price_threshold`, `candle_direction`, and `funding_threshold` markets.
- I settle every market from one explicit Pacifica-derived source.
- I do not include AI summaries, real trading, withdrawals, or wallet-linked execution.

## Backend Status

The backend currently provides:

- `POST /api/v1/players/guest`
- `GET /api/v1/players/me`
- `GET /api/v1/players/me/balance`
- `POST /api/v1/markets`
- `GET /api/v1/markets`
- `GET /api/v1/markets/{market_id}`
- `POST /api/v1/markets/{market_id}/positions`
- `GET /api/v1/players/me/positions`

Important behavior already implemented:

- Market creation now includes the creator's required opening side and stake.
- `POST /api/v1/markets` creates the market and the creator's first position in one transaction.
- Additional players can join through the position-placement route.
- Price markets settle from Pacifica mark-price snapshots.
- Candle markets settle from Pacifica mark-price candles.
- Funding markets settle from Pacifica funding history.
- Settlement updates market result, position outcomes, and virtual balance effects in one transaction.

## Payout Model

I am using a simple fixed-odds house rule for v1:

- every position locks the chosen stake at entry time
- losers receive no further balance credit at settlement
- winners receive `2x stake`

This keeps the demo easy to explain and easy to verify.

## Why This Repo Looks The Way It Does

I am optimizing for correctness and clarity over feature sprawl.

That means:

- Go backend
- PostgreSQL persistence
- Pacifica integration kept read-only
- backend-owned settlement logic
- backend-owned balance authority

## Where To Look

- To find the current repository map visit [structure.md](./structure.md).
- To find the architecture decisions and tradeoffs visit [architecture.md](./architecture.md).
- To find schema reasoning visit [schema.md](./schema.md).
- To find the execution record and task status visit [task.md](./task.md).
- To find frontend flow notes for creator auto-stake visit [frontend-notes.md](./frontend-notes.md).
- To find backend package ownership visit [backend/README.md](./backend/README.md).
