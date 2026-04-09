# Architecture

## Purpose

I use this document to describe the current Pacifica Pulse implementation boundary, not just the intended v1 shape.

This version reflects the code that exists in the repository today across the SvelteKit frontend, the Go backend, PostgreSQL persistence, the settlement worker, and the derived leaderboard flow.

---

## 1. Current Product Boundary

### What Pacifica Pulse is today

Pacifica Pulse is a read-only prediction platform that:

- consumes Pacifica testnet market data
- lets guest players create short-duration prediction markets
- lets players take virtual YES / NO positions
- settles markets from Pacifica data
- pays out virtual winnings with fixed transparent math
- exposes a public leaderboard derived from stored activity

### What it still is not

I am explicitly not building these in the current product:

- AI summaries
- real trading
- wallet authentication
- withdrawals
- Pacifica signed write operations
- reward redemption
- internal open-interest snapshot markets

### Virtual value model

The app uses virtual balances tied to real Pacifica-derived outcomes.

That means:

- users do not trade real assets
- users do not need a wallet to participate
- the backend stays authoritative for all balances and settlements
- the product behaves like a prediction game, not like an exchange client

---

## 2. Implemented Market Scope

### Supported market types

The implementation currently supports these market types:

- `price_threshold`
- `candle_direction`
- `funding_threshold`

### Deferred market types

These remain out of the live implementation:

- open-interest markets
- trade-flow markets
- volume-derived markets

Open interest is still the most realistic v1.5 extension because it mainly needs internal snapshot capture and explicit settlement rules.

---

## 3. Stack And System Shape

### Chosen stack

- Frontend: SvelteKit + Svelte 5 runes + Tailwind CSS
- Backend: Go
- Database: PostgreSQL
- Realtime transport: SSE from our backend
- Upstream data strategy: Pacifica REST-first, with WebSocket support isolated behind the backend

### Why PostgreSQL without Redis still works here

I am intentionally keeping Redis out of the current system because:

- current scale does not justify another moving part
- authoritative state already lives in PostgreSQL
- the leaderboard can be handled as derived queries plus focused indexes
- correctness and operational simplicity matter more than speculative cache layers right now

### High-level request flow

The live app now looks like this:

1. The browser hits SvelteKit routes.
2. SvelteKit proxies `/api/*` calls to the Go backend.
3. The Go backend owns sessions, balance rules, markets, positions, settlements, leaderboard reads, and SSE.
4. PostgreSQL stores identity, balances, markets, positions, and settlement audit records.
5. Pacifica stays behind the backend for context reads, settlement reads, and optional realtime ingestion.

---

## 4. Current Frontend Architecture

### Current pages

The frontend currently includes these route groups:

- landing page at `/`
- dashboard at `/dashboard`
- create-market flow at `/markets/create`
- market detail at `/markets/[id]`
- resolved market view at `/markets/[id]/resolved`
- portfolio at `/portfolio`
- leaderboard at `/leaderboard`

### Frontend responsibilities

The frontend is responsible for:

- presenting real market state
- creating or reusing guest sessions before protected actions
- collecting market creation and position-placement input
- rendering balance, position, and market detail reads
- rendering the leaderboard snapshot returned by the backend

### Frontend boundary rules

I keep these rules in place:

- the frontend should not talk directly to Pacifica
- the frontend should not compute settlement or balance authority
- the frontend should not rank leaderboard results itself
- the frontend should prefer one backend request per page snapshot where that stays practical

### Frontend data-loading pattern

The codebase currently uses two frontend data-loading styles:

- lightweight client-side fetch flows for dashboard, create-market, market detail, and portfolio
- a route-level page load for the leaderboard snapshot so the route renders from one backend response

This is a practical middle ground for the current repo. I have not forced every page into one loading strategy yet.

### Frontend proxy layer

SvelteKit exposes a catch-all `/api/[...path]` proxy route so the browser always talks to the app shell and forwards cookies cleanly to the Go backend.

That gives me:

- a single browser-facing origin
- cookie continuity
- less frontend environment branching
- a cleaner place to centralize backend URL configuration

---

## 5. Current Backend Architecture

### Backend modules

The backend is currently split into these modules:

- `auth`
- `player`
- `balance`
- `leaderboard`
- `market`
- `position`
- `settlement`
- `pacifica`
- `realtime`
- `config`
- `storage`

### Route wiring

`backend/cmd/api/main.go` remains the composition root.

That means:

- modules export controllers or services
- `main.go` wires dependencies and repositories
- routes are registered explicitly in one place
- the settlement worker is started beside the HTTP server

### Current API surface

The implemented backend routes are:

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

### Public versus session-protected reads

The current split is deliberate:

- guest identity, balance, positions, and market creation need a valid session
- market listings, market detail, realtime stream, and leaderboard reads are public

That keeps the discovery surface simple while still protecting player-specific state.

---

## 6. Market Creation And Position Flow

### Creator flow

Market creation is no longer metadata-only.

Every created market requires:

- market definition
- creator side selection
- creator stake amount

The backend creates the market and the creator's opening position in one transaction.

### Additional participant flow

Other players join existing markets through:

- `POST /api/v1/markets/{market_id}/positions`

That route validates:

- session ownership
- market eligibility
- side input
- stake amount
- balance availability

### Why I am keeping this product-shaped

I prefer this contract because:

- the frontend stays simple
- the backend owns the transactional rules
- creator participation is guaranteed at creation time
- balance locking stays authoritative

---

## 7. Settlement Architecture

### Core principle

Every market resolves from one explicit Pacifica-derived source.

I do not let the app settle from vague or approximate values.

### Current settlement paths

- price markets settle from Pacifica mark-price snapshots
- candle markets settle from Pacifica mark-price candles
- funding markets settle from Pacifica funding records

### Current execution model

The settlement worker:

- scans for active markets nearing expiry
- batches price-fetch planning
- retries if a returned price snapshot predates the actual expiry boundary
- resolves the market through the correct resolver path
- writes settlement audit data
- updates markets, positions, and balances in one transaction

### Why REST-first is still the default

I continue to prefer a final authoritative REST fetch instead of settling purely from in-memory live subscriptions because:

- reconnect recovery is simpler
- auditability is better
- the logic is easier to reason about under hackathon scale

---

## 8. Leaderboard Architecture

### Current leaderboard design

The leaderboard is now a real product feature, not a placeholder page.

It is implemented as:

- one public backend snapshot route at `GET /api/v1/leaderboard`
- one Go service that fans out category reads in parallel
- one PostgreSQL aggregation layer that derives rankings from stored data
- one SvelteKit page load that renders from the backend snapshot

### Categories currently exposed

- top predictors
- top creators
- best streaks
- most active

### Why I chose a derived query layer first

This is the right current tradeoff because:

- current scale does not justify Redis
- current scale does not justify a leaderboard table or materialized view
- the queries stay transparent and auditable
- the rankings are derived from already authoritative data

### Performance controls already in place

- focused leaderboard indexes in the migrations
- one snapshot endpoint instead of multiple category endpoints
- parallel category reads in Go
- cache headers on the public leaderboard route

---

## 9. Data Model

### Core persistent tables

The current implementation relies on these primary tables:

- `players`
- `player_sessions`
- `player_balances`
- `markets`
- `positions`
- `market_settlements`

### What PostgreSQL owns today

PostgreSQL currently stores:

- guest player identity
- secure session lookup material
- available and locked balances
- market definitions and lifecycle state
- player positions and payout meaning
- settlement audit records

### What is still derived instead of stored separately

These are still query-derived layers:

- leaderboard rankings
- dashboard grouping of active versus resolved markets
- player-facing summaries built from markets and positions

That is intentional. I do not want to add write-time denormalization before there is a real scale need.

---

## 10. Realtime And Read Performance

### Realtime path

The realtime surface currently uses SSE from the backend.

The stream is meant to carry backend-owned event types such as:

- `market.created`
- `market.updated`
- `market.settled`

### Read-performance choices already implemented

I have already made these platform-level choices:

- backend-owned Pacifica access instead of direct browser fan-out
- PostgreSQL-first reads with no Redis
- one proxy layer in SvelteKit for browser API calls
- one snapshot endpoint for leaderboard
- focused Postgres indexes for leaderboard reads

### Where I have intentionally not optimized yet

- I have not introduced Redis
- I have not introduced materialized leaderboard views
- I have not made every page fully server-rendered
- I have not added aggressive background prefetching

Those are conscious deferrals, not omissions by accident.

---

## 11. Current Implementation Gaps

The codebase is in a stronger state than the earlier planning docs implied, but a few gaps still matter:

- the landing page still contains unsupported marketing language such as AI and signals that do not match the current v1 scope
- the frontend still uses `adapter-auto` instead of `@sveltejs/adapter-vercel`
- some frontend flows still rely on post-mount fetching rather than route-level loads
- the leaderboard is live, but deeper player profile drill-downs are not part of the current product

I want the documentation to reflect these gaps clearly so the repo does not overstate what is shipped.

---

## 12. Current Roadmap

### Current v1 baseline

- guest player accounts
- virtual balances
- price threshold markets
- candle outcome markets
- funding markets
- Pacifica read-only integration
- deterministic settlement
- public leaderboard snapshot

### v1.5 candidates

- open-interest markets backed by internal snapshot storage
- more complete landing-page cleanup
- broader server-side page loading where it improves perceived speed

### v2 candidates

- wallet auth
- reward mapping or redemption
- richer leaderboard and player-profile surfaces

### v3 candidates

- real Pacifica trade bridge
- builder-code integration
- API Agent Keys

---

## 13. Final Position

The current architecture is a practical v1 game platform, not just a backend prototype.

The key decisions that now define the implementation are:

- keep Pacifica integration read-only
- keep balances and settlement server-authoritative
- keep Redis out until scale justifies it
- use product-shaped backend routes
- derive leaderboard reads from PostgreSQL instead of building a separate cache system
- keep the frontend pointed at our backend and proxy layer, not at Pacifica directly

That gives me a product that is small enough to move quickly, but structured enough to grow without rewriting the foundations.
