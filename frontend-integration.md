# Frontend Integration Report

## Goal

I want this document to describe the current frontend integration state, not just the original plan.

## Current Integration Status

The frontend is now meaningfully connected to the backend:

- guest session flow exists
- guest player id and display name are cached in localStorage to avoid silent duplicate guest provisioning on revisit
- dashboard reads real markets, balances, and positions
- create-market reads backend context and submits the creator's opening position in one request
- market detail reads a real market and can place a real position
- portfolio reads real player state
- leaderboard reads a real backend snapshot
- `/api/[...path]` proxies browser calls to the backend

## What Is Working Well

- The frontend does not call Pacifica directly.
- The backend owns settlement, balances, and ranking logic.
- The localStorage guest cache is only a frontend continuity layer; protected backend calls still depend on the secure cookie session.
- The leaderboard is efficient because the page renders from one snapshot response.
- The create-market flow stays clean because the creator-side logic is backend-owned.

## What Is Still Uneven

- The landing page still behaves more like concept marketing than the live product.
- Some pages use route loads and others fetch on mount, so the loading strategy is not fully unified.
- Realtime support exists at the backend boundary, but the frontend is still conservative about how much live-stream behavior it consumes.

## Current Backend Routes The Frontend Depends On

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

## Recommended Next Cleanup Pass

1. Bring the landing page copy back in line with the supported v1 features.
2. Decide whether more app pages should move to route-level server loads.
3. Switch the frontend adapter from `adapter-auto` to `@sveltejs/adapter-vercel`.
4. Tighten realtime usage only where it improves the real product loop.

## Final Note

The frontend is no longer pretending to be the product. Most of the core app loop is now connected to backend truth, and the remaining work is mainly cleanup, consistency, and deployment alignment.
