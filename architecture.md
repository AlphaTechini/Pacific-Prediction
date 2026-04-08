# Architecture

## Purpose

I am defining the v1 architecture for Pacifica Pulse as a prediction-game layer built on top of Pacifica testnet market data.

This document turns the decisions from [`Readme.md`](./Readme.md), [`Resources.md`](./Resources.md), and [`Research.md`](./Research.md) into an implementation boundary I can actually build against without scope drift.

Research baseline used for this document: April 7, 2026.

---

## 1. Confirmed Product Boundary

### What v1 is

Pacifica Pulse v1 is a read-only prediction platform that:

- consumes Pacifica testnet market data
- lets users create short-duration prediction markets
- lets users take virtual YES / NO positions
- settles markets from Pacifica data
- pays out virtual winnings mathematically

### What v1 is not

I am explicitly not building these in v1:

- AI summaries
- real trading
- wallet authentication
- withdrawals
- Pacifica signed write operations
- builder-code integration
- API Agent Keys
- token transfer flows

### Virtual value model

The app will use virtual balances, but those balances will still be tied to real Pacifica-derived market outcomes.

That means:

- users do not trade real assets
- users do not need a wallet to participate
- winnings are calculated from the prediction engine
- the game reflects what would have happened under the chosen prediction rules

---

## 2. Prediction Types

### Supported in v1

I am locking v1 to the market types that Pacifica supports most cleanly with the highest settlement confidence.

#### 1. Price threshold markets

Examples:

- Will BTC mark price be above 105000 in 30 minutes?
- Will ETH mark price be below 4200 at expiry?

Why this is in v1:

- Pacifica exposes live price data
- Pacifica exposes historical mark-price candle data
- settlement is easy to explain and easy to verify

#### 2. Candle outcome markets

Examples:

- Will SOL close the next 5m candle bullish?
- Will BTC close the next 15m candle bearish?

Why this is in v1:

- Pacifica exposes candle streams and mark-price candle streams
- candle open and close are clear settlement points
- the UI can show the exact candle interval being judged

#### 3. Funding direction or threshold markets

Examples:

- Will BTC funding remain positive at the next funding settlement?
- Will ETH funding be above 0.00001 at the next interval?

Why this is in v1:

- Pacifica exposes funding-related data
- funding markets are differentiated enough from plain price bets
- they fit the product story without needing derived AI logic

### Explicitly deferred from v1

#### Open interest change markets

I am moving this out of v1.

Reason:

- Pacifica exposes live open interest
- I did not confirm a dedicated historical OI endpoint in the reviewed docs
- reliable settlement would require our backend to store reference and expiry snapshots

This is the easiest deferred market type to bring back later because it only needs internal snapshot support, not an entirely new product model.

#### Trade-flow markets

I am not building these in v1.

Reason:

- they require continuous trade aggregation by our backend
- settlement depends on our own derived rollups rather than one clean Pacifica field
- they are harder to explain during a demo

#### Volume-derived markets

I am not building these in v1.

Reason:

- rolling volume metrics are less intuitive for short-horizon predictions
- they add complexity without improving the core demo
- they are weaker than price, candle, and funding markets for clarity

### Deferred roadmap

I only want to carry forward the deferred items that are still relatively safe to implement later.

#### v1.5 candidate

- open interest change markets, after I add backend snapshot capture and explicit OI settlement rules

#### not prioritized yet

- trade-flow markets
- volume-derived markets

---

## 3. Architecture Direction

### Chosen stack

- Frontend: SvelteKit + Tailwind
- Backend: Go
- Database: PostgreSQL
- Realtime transport to UI: server-sent events or WebSocket from our backend
- Pacifica integration: REST + WebSocket, read-only

### Why I am choosing Go for the backend

Go is the best fit for the hard parts of v1:

- long-lived Pacifica WebSocket connections
- deterministic settlement jobs
- typed data models
- a clean single-service deployment path

### Tradeoff versus Fastify

Fastify would reduce context switching because the frontend is already TypeScript-based.

I am still choosing Go because v1 reliability depends more on:

- ingestion stability
- job scheduling
- settlement correctness
- clear service boundaries

than on sharing types with the frontend.

---

## 4. System Boundary

### Pacifica owns

- live market data
- historical market data
- market metadata
- funding and candle truth used for settlement

### Pacifica Pulse owns

- prediction market definitions
- virtual user balances
- virtual positions
- payout calculations
- market countdowns
- settlement orchestration
- resolved-market history
- UI-facing aggregation

### Why the frontend should not talk to Pacifica directly

I want the backend in the middle because it gives me:

- one place to manage Pacifica rate limits
- one place to maintain WebSocket heartbeats and reconnect logic
- one settlement source of truth
- simpler frontend code
- cleaner future extension into auth or rewards

---

## 5. High-Level Component Design

### Frontend

The SvelteKit frontend will contain:

- dashboard page for active and resolved markets
- market creation flow
- market detail page
- portfolio page for virtual positions and balance
- live countdown and result UI

Frontend responsibilities:

- rendering live market state
- collecting market creation inputs
- placing virtual positions
- showing virtual PnL and resolved outcomes

Frontend non-responsibilities:

- no direct Pacifica API access
- no settlement logic
- no authoritative balance updates

### Backend API

The Go backend will expose:

- guest session endpoints
- player profile endpoints
- balance endpoints
- market creation endpoints
- market listing and detail endpoints
- virtual position endpoints
- internal settlement workers
- realtime event feed for frontend updates

### Data layer

PostgreSQL will store:

- player records
- player balances
- markets
- market snapshots we choose to retain
- positions
- payouts
- settlement audit records

I am not adding Redis in v1.

Why:

- the MVP does not need distributed caching yet
- PostgreSQL is enough for correctness and demo scale
- adding Redis now would increase moving parts without solving a proven bottleneck

---

## 6. Identity Without Wallet Auth

### Chosen approach

I will use guest player identities in v1.

Each player gets:

- a backend-issued player ID
- an opaque session token stored in a secure cookie
- an optional display name for leaderboard-style UI later

### Why this is the right v1 choice

- no wallet is required
- balances can still persist across page refreshes
- the backend can enforce authoritative balance updates
- this leaves room for future account linking without rebuilding the whole data model

### Tradeoff versus local-only storage

Local-only storage would be faster to prototype, but I am not choosing it because:

- users would lose progress across devices
- virtual balances would be trivial to tamper with
- later migration to real accounts would be messy

---

## 7. Market Definition Rules

### Market schema intent

Every market must carry enough data to settle without interpretation drift.

Required fields:

- `id`
- `title`
- `symbol`
- `market_type`
- `condition_operator`
- `threshold_value`
- `source_type`
- `source_interval`
- `reference_value`
- `expiry_time`
- `status`
- `result`
- `settlement_value`
- `resolved_at`
- `resolution_reason`

### Market type mapping

#### Price threshold

- `market_type = price_threshold`
- `source_type = mark_price`
- `source_interval = null`

#### Candle outcome

- `market_type = candle_direction`
- `source_type = mark_price_candle`
- `source_interval = 1m | 5m | 15m | ...`

#### Funding threshold

- `market_type = funding_threshold`
- `source_type = funding_rate`
- `source_interval = funding_epoch`

### Why I am standardizing candle settlement on mark-price candles

Pacifica exposes both trade-price candles and mark-price candles.

I am choosing mark-price candles for v1 candle settlement because:

- they align better with the rest of the prediction framing
- they reduce ambiguity when price spikes or low-liquidity prints distort traded candles
- they make the settlement source more consistent

Tradeoff:

- trade-price candle fans may expect raw traded closes instead
- I accept that tradeoff because v1 needs consistency more than market microstructure nuance

---

## 8. Settlement Rules

### Core settlement principle

Every market must resolve from one authoritative Pacifica-derived value at one exact timestamp boundary.

I will not settle any market using vague wording like:

- around expiry
- near close
- recent average

### Settlement by market type

#### Price threshold markets

Settlement rule:

- fetch the first mark-price value at or immediately after expiry using the selected Pacifica source path
- compare it against the market threshold

Result examples:

- `mark_price > threshold` => YES
- `mark_price <= threshold` => NO

#### Candle outcome markets

Settlement rule:

- resolve against the close of the selected mark-price candle interval
- compare candle close against candle open

Result examples:

- `close > open` => bullish => YES
- `close <= open` => bearish-or-flat => NO

I am intentionally treating flat as NO in v1 because it removes draw handling complexity.

#### Funding markets

Settlement rule:

- resolve using the first funding record whose settlement timestamp is at or after expiry
- compare the final funding value with the threshold or sign rule

### Settlement execution model

I will use a backend worker that:

- polls for markets approaching expiry
- fetches authoritative Pacifica settlement data
- computes result
- updates balances and positions in one database transaction
- records the raw settlement value and source used

### Why I am not settling from only live WebSocket memory

Live subscription state is useful for UI updates, but it should not be the only settlement source.

I want a final fetch step because it gives me:

- cleaner recovery after reconnects
- more deterministic settlement
- better auditability

---

## 9. Virtual Balance and Payout Model

### Chosen v1 model

I am choosing a simple fixed-odds prediction game model with virtual stakes.

Each player has:

- a starting virtual balance
- stakeable amounts per market
- winnings or losses applied after settlement

### Payout approach

For v1, I want fixed transparent math over financial realism.

Recommended v1 rule:

- each position stakes a chosen virtual amount
- if the prediction is correct, payout = stake + profit
- if the prediction is wrong, loss = full stake

### Simpler profit option

The cleanest first version is fixed-odds by side at entry time or an even simpler house rule.

Recommended default:

- YES and NO both start at even odds
- payout is `2x stake` on win and `0` on loss, minus no platform fee in v1

Why I prefer this first:

- easy to explain in a demo
- easy to verify
- no AMM design
- no pool-balancing complexity

### Deferred alternative

A pooled-odds model can come later if I want:

- changing odds
- crowd sentiment visuals
- richer market dynamics

I am not choosing it now because it adds balancing and display complexity without improving the first demo enough.

---

## 10. Backend Modules

### Module list

I am splitting the backend into these modules:

- `auth`
- `player`
- `balance`
- `market`
- `position`
- `settlement`
- `pacifica`
- `realtime`
- `config`
- `storage`

This keeps the upgradeability benefits of a more layered design without turning the backend into a scattered set of tiny services.

### `auth`

Owns:

- guest session creation
- session validation
- session lookup
- future wallet or external auth upgrades

Possesses:

- auth-facing request and response schemas
- session models
- auth service and session service
- auth controller exports for route wiring in `main.go`

### `player`

Owns:

- player identity and profile only

Possesses:

- player schemas
- player service
- player controller exports for profile access

Boundary rule:

- `player` does not own balances or sessions

### `balance`

Owns:

- virtual balance state
- balance locking and unlocking
- debits, credits, and payout application

Possesses:

- balance schemas
- balance service
- balance controller exports for read access

Boundary rule:

- all spendable-value mutations must pass through `balance`

### `market`

Owns:

- market creation
- market validation
- market listing and market detail
- supported market-type rules

Possesses:

- market schemas
- market service
- market validation service
- market controller exports

### `position`

Owns:

- placing YES or NO positions
- reading a player's position history

Possesses:

- position schemas
- position service
- position controller exports

Boundary rule:

- `position` depends on `market` for market eligibility and on `balance` for stake locking

### `settlement`

Owns:

- expiry scanning
- deterministic market resolution
- payout orchestration
- settlement audit creation

Possesses:

- settlement worker
- settlement service
- payout coordination service
- internal settlement schemas

Boundary rule:

- `settlement` coordinates resolution, but `balance` performs balance mutations and `pacifica` provides source data

### `pacifica`

Owns:

- Pacifica REST integration
- Pacifica WebSocket integration
- subscription and heartbeat management
- normalization of vendor data into internal data shapes

Possesses:

- Pacifica DTOs
- REST and WebSocket client services
- subscription service

### `realtime`

Owns:

- SSE stream delivery
- backend event fan-out to the client-facing stream

Possesses:

- event schemas
- stream controller export
- realtime hub and publisher services

### `config`

Owns:

- environment loading
- config validation
- runtime configuration structs

### `storage`

Owns:

- database connection
- transactions
- repository wiring
- persistence implementations

### Route wiring rule

Routes will remain assigned in `main.go`.

That means:

- modules export controllers or handler functions
- `main.go` constructs dependencies
- `main.go` assigns routes explicitly

I am choosing this because it keeps startup and route ownership easy to scan.

## 11. API Shape

### External app endpoints by module

#### `auth`

- `POST /api/v1/players/guest`

#### `player`

- `GET /api/v1/players/me`

#### `balance`

- `GET /api/v1/players/me/balance`

#### `market`

- `POST /api/v1/markets`
- `GET /api/v1/markets`
- `GET /api/v1/markets/:id`

#### `position`

- `POST /api/v1/markets/:id/positions`
- `GET /api/v1/players/me/positions`

#### `realtime`

- `GET /api/v1/stream`

### Internal-only jobs

- market-expiry scanner
- settlement worker
- Pacifica subscription manager

### Why I am not exposing Pacifica-shaped endpoints to the frontend

I want our API to be product-shaped, not vendor-shaped.

That keeps the frontend focused on:

- markets
- positions
- balances
- outcomes

instead of making the UI understand Pacifica transport details.

---

## 12. Database Design

### Chosen schema direction

I am choosing Option C as the database direction for v1.

That means I want a middle-ground schema that is:

- simple enough to ship quickly
- explicit enough to settle markets safely
- flexible enough to grow into v1.5 and v2 without a rewrite

Why I am choosing it:

- it keeps identity, balances, market state, positions, and settlement audit separate
- it supports guest users now and wallet linking later
- it supports fixed-odds v1 now and additive feature growth later

### Core tables

#### `players`

Purpose:

- persistent user identity independent of wallet auth

Core fields:

- `id`
- `display_name`
- `created_at`
- `updated_at`

#### `player_sessions`

Purpose:

- backend-controlled guest sessions without storing raw session tokens

Core fields:

- `id`
- `player_id`
- `session_token_hash`
- `expires_at`
- `created_at`

#### `player_balances`

Purpose:

- authoritative virtual balance state with lockable funds

Core fields:

- `player_id`
- `available_balance`
- `locked_balance`
- `updated_at`

Why the split matters:

- `available_balance` handles spendable value
- `locked_balance` prevents double-spending while markets are unresolved

#### `markets`

Purpose:

- market definition plus lifecycle state in one readable v1 record

Core fields:

- `id`
- `title`
- `symbol`
- `market_type`
- `condition_operator`
- `threshold_value`
- `source_type`
- `source_interval`
- `reference_value`
- `expiry_time`
- `status`
- `result`
- `settlement_value`
- `resolved_at`
- `resolution_reason`
- `created_by_player_id`
- `created_at`

Why this is enough for v1:

- the market can be settled later without guessing intent
- the read path stays simple because definition and outcome live together

#### `positions`

Purpose:

- each virtual YES or NO participation entry

Core fields:

- `id`
- `player_id`
- `market_id`
- `side`
- `stake_amount`
- `potential_payout`
- `status`
- `created_at`
- `settled_at`

Why I store `potential_payout`:

- fixed-odds values should be locked at entry time
- future payout-rule changes should not rewrite old position economics

#### `market_settlements`

Purpose:

- auditable settlement record tied to the exact Pacifica-derived source used

Core fields:

- `id`
- `market_id`
- `pacifica_source`
- `source_timestamp`
- `raw_payload`
- `settlement_value`
- `result`
- `created_at`

Why this table matters:

- settlement is the highest-risk backend operation
- the audit trail makes disputes and debugging much easier

### Optional v1 table

#### `market_snapshots`

I am keeping this optional in true v1.

Use cases:

- charting
- replay
- debugging
- future OI support

I only want to add it early if I decide one of these becomes a real requirement.

### Why this schema scales well

This schema is easy to extend because most future work is additive.

Examples:

- wallet auth can add `player_wallets` or `auth_identities`
- rewards can add a `reward_ledger`
- leaderboard views can derive from balances, positions, and settlements
- open interest support can add `market_snapshots`
- advanced payout systems can add odds history or pool tables later

### Where the schema would start to strain

I would revisit the design if the product grows into:

- dynamic odds or pooled liquidity
- highly custom multi-condition market logic
- event-sourced analytics at larger scale
- real trade execution with full order lifecycle tracking

### Supporting documentation

I am keeping the full schema reasoning in [`schema.md`](./schema.md) so implementation can rely on one dedicated source of truth when migrations start.

---

## 13. Pacifica Integration Design

### Read-only sources used in v1

- market info
- prices
- mark-price candles
- funding history
- live price subscriptions
- live mark-price candle subscriptions

### Pacifica integration modules

#### REST client

Responsibilities:

- fetch market metadata
- fetch settlement values
- fetch candle history for market creation context
- fetch funding data for funding market resolution

#### WebSocket client

Responsibilities:

- maintain Pacifica subscription connections
- send heartbeat pings
- reconnect safely
- fan live updates into backend state and frontend streams

### Failure strategy

If Pacifica connectivity drops:

- market creation should remain available only if recent source data is still fresh enough
- active market UIs should show delayed-data status
- settlement should retry instead of guessing

I would rather delay settlement briefly than settle from uncertain data.

---

## 14. Realtime Delivery To The Frontend

### Chosen direction

I prefer server-sent events first, with a later upgrade path to WebSocket if the UI truly needs bidirectional realtime.

Why SSE fits v1:

- simpler transport for live dashboards
- enough for pushing market updates, countdown state, and settlement results
- easier operationally than adding another bidirectional channel immediately

Tradeoff:

- WebSocket is more flexible long term
- SSE is enough for the first product shape because the frontend mostly consumes updates

---

## 15. Security and Abuse Notes

### v1 security focus

Because there is no wallet auth and no real asset movement, the main risks are:

- session abuse
- market spam
- virtual balance tampering attempts
- excessive Pacifica fan-out

### Controls I want in v1

- secure signed session cookies
- server-authoritative balances
- rate limiting on market creation
- rate limiting on position placement
- backend-only Pacifica access
- configuration from environment variables only

---

## 16. Versioned Roadmap

### v1

- guest player accounts
- virtual balances
- price threshold markets
- candle outcome markets
- funding markets
- Pacifica read-only integration
- deterministic settlement

### v1.5

- open interest markets backed by internal snapshot storage
- richer charts and market history
- better odds visualization if needed

### v2

- wallet auth
- optional reward redemption or testnet token mapping
- leaderboard and gamification polish

### v3

- real Pacifica trade bridge
- builder-code integration
- API Agent Keys
- optional AI assistance if the product still benefits from it

---

## 17. Final Decision Summary

I am building Pacifica Pulse v1 as a simple, high-confidence, read-only prediction game.

The key architecture decisions are:

- keep Pacifica integration read-only
- use Go for backend reliability
- use PostgreSQL without Redis in v1
- use guest identities instead of wallet auth
- support only price, candle, and funding markets
- settle every market from one explicit Pacifica-derived source
- use virtual balances and transparent payout math
- defer open interest to the next safe version
- leave trade-flow, volume-derived markets, AI, and real trading out of v1

This is the smallest architecture that still feels like the real product instead of a loose prototype.


