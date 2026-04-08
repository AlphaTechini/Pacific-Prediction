# Research

## Scope

I reviewed the local project briefs first, then the external sources listed in [`Resources.md`](./Resources.md), and pulled out the details that directly affect planning Pacifica Pulse.

Research date: April 7, 2026

---

## 1. Project Goals From `Readme.md`

Source: [`Readme.md`](./Readme.md)

### What the project is trying to do

Pacifica Pulse is positioned as a short-duration prediction market layer that sits on top of Pacifica market data instead of trying to replace Pacifica itself.

The MVP goal is clear:

- let users create binary markets around live Pacifica signals
- let users take YES / NO positions
- resolve markets from Pacifica data
- keep the product demo-friendly and simple

### Core data signals the product depends on

The readme repeatedly centers the app around these Pacifica-derived inputs:

- mark price
- funding rate
- open interest
- candle data
- recent trades

### Important product constraints already defined

- Pacifica APIs must be the primary source of truth
- testnet support is required
- demo clarity matters more than sophisticated market mechanics
- the MVP should avoid building a full AMM or liquidity pool design

### Immediate architecture consequence

This project should treat Pacifica as the external market-data and settlement authority, while our own system owns:

- prediction market definitions
- user positions
- pooled-odds logic
- settlement scheduling
- historical UI state

---

## 2. Source Review Order

I reviewed the sources in this order:

1. [`Resources.md`](./Resources.md)
2. [`Readme.md`](./Readme.md)
3. Pacifica Builder Program
4. Pacifica API documentation
5. Pacifica Python SDK
6. Pacifica mainnet and testnet docs
7. Sponsor tool pages for Fuul, Rhino.fi, Privy, and Elfa AI

I could not review `Pacifica Discord API Channel: Channel Link` because the resource in `Resources.md` is only a placeholder, not a real URL.

---

## 3. Pacifica Builder Program Findings

Source: https://docs.pacifica.fi/builder-program

### What matters

- Pacifica lets builders attach a `builder_code` to supported order-creation requests.
- Users must approve a builder code before it can be used on their orders.
- That approval is signed and submitted to Pacifica through account endpoints.
- The builder program can reward growth-driving apps with Pacifica points.

### Useful details for this project

- The builder rewards program was extended through June 12, 2026.
- Builder codes can be attached to:
  - market orders
  - limit orders
  - stop orders
  - TP/SL requests
- Pacifica exposes builder approval, revoke, trade history, and leaderboard endpoints.

### Why this matters for Pacifica Pulse

This is not required for the prediction-market MVP, but it is directly relevant for the Phase 3 "signal-to-trade bridge" described in `Readme.md`.

If Pacifica Pulse later lets a user press "Trade this idea on Pacifica", then:

- builder-code support can create a monetization and hackathon-growth angle
- user authorization must be part of onboarding
- trade execution needs a signed request path, not just read-only market data

### Planning implication

I should keep the trading bridge isolated from the core prediction-market flow. The MVP should not depend on builder-code approval to function.

---

## 4. Pacifica API Findings

Sources:

- https://docs.pacifica.fi/api-documentation/api
- https://pacifica.gitbook.io/docs/api-documentation/api/websocket
- https://pacifica.gitbook.io/docs/api-documentation/api/rest-api/markets/get-market-info
- https://pacifica.gitbook.io/closed-alpha/api-documentation/api/rest-api/markets/get-prices
- https://pacifica.gitbook.io/docs/api-documentation/api/rest-api/markets/get-mark-price-candle-data
- https://pacifica.gitbook.io/docs/api-documentation/api/rest-api/markets/get-recent-trades
- https://pacifica.gitbook.io/closed-alpha/api-documentation/api/rest-api/markets/get-historical-funding
- https://pacifica.gitbook.io/docs/api-documentation/api/market-symbols
- https://pacifica.gitbook.io/docs/api-documentation/api/signing/operation-types

### The read-only API surface is enough for the MVP

The Pacifica docs already expose the exact data classes the readme needs:

- market metadata via `GET /api/v1/info`
- prices via `GET /api/v1/info/prices`
- historical mark-price candles via `/api/v1/kline/mark`
- recent trades via `GET /api/v1/trades`
- historical funding via `GET /api/v1/funding_rate/history`
- live subscriptions for `prices`, `candle`, `mark_price_candle`, and `trades`

One gap matters for scope:

- I found live open-interest data in Pacifica price surfaces, but I did not find a dedicated historical open-interest endpoint in the reviewed docs

### WebSocket details that matter

- Mainnet WebSocket URL: `wss://ws.pacifica.fi/ws`
- Testnet WebSocket URL: `wss://test-ws.pacifica.fi/ws`
- Pacifica closes idle connections after 60 seconds without a message.
- Long-lived connections are capped at 24 hours.
- The app should send `{"method":"ping"}` heartbeats to keep connections alive.

### Data-field observations that affect our schema

From the public docs:

- `prices` streams include `funding`, `mark`, `open_interest`, `oracle`, `volume_24h`, and `timestamp`
- `candle` streams include open, close, high, low, volume, interval, and start/end timestamps
- `mark_price_candle` streams are available with the same candle shape and explicit interval selection
- recent-trade responses include side, cause, amount, price, and a global ordering field

### Important correctness details

- Pacifica market symbols are case-sensitive
- market metadata exposes tick size, lot size, leverage, and order-size constraints
- those constraints matter if we later bridge into real trading
- for prediction-only logic, they still help validate which symbols are safe to expose in the UI

### Signing and auth implications

- GET requests and WebSocket subscriptions do not require signing
- POST requests require Ed25519-style signed payloads
- Pacifica documents operation types for every signed trading/account mutation
- API Agent Keys exist for delegated trading flows

### Planning implication

The MVP should stay read-only against Pacifica:

- subscribe to live price feeds
- fetch candles and funding history for market creation and settlement
- do not make signed Pacifica requests in Phase 1 unless the scope expands into real trade execution

For open-interest-based markets specifically, I should either:

- defer them out of the first MVP, or
- settle them using snapshots our backend captures over time instead of relying on a Pacifica historical OI endpoint

---

## 5. Rate Limit, Ordering, and Error Findings

Sources:

- https://pacifica.gitbook.io/closed-alpha/api-documentation/api/rate-limits
- https://pacifica.gitbook.io/closed-alpha/api-documentation/api/error-codes
- https://pacifica.gitbook.io/closed-alpha/api-documentation/api/last-order-id

### Rate limits

Pacifica uses a credit-based rolling 60-second rate limit.

The docs currently show:

- unidentified IP base quota: 125 credits per 60 seconds
- valid API config key base quota: 300 credits per 60 seconds
- WebSocket max: 300 concurrent connections per IP
- WebSocket max: 20 subscriptions per channel per connection

### Ordering

Pacifica documents `last_order_id` as an exchange-wide ordering primitive across trading-related responses and streams.

### Errors

The docs explicitly list:

- `403` forbidden, including cases such as no access code or restricted region
- `429` too many requests
- `422` business-logic errors

### Planning implication

The backend should sit between the frontend and Pacifica instead of every browser opening its own Pacifica subscriptions. That gives us:

- one controlled WebSocket session pool
- centralized retry and heartbeat handling
- better rate-limit protection
- consistent settlement and event ordering

---

## 6. Pacifica Python SDK Findings

Source: https://github.com/pacifica-fi/python-sdk

### What is available

The SDK repository is not a full packaged framework. It is an examples repo with separate `rest` and `ws` folders plus shared helper code.

The README says it includes examples for:

- obtaining market data
- monitoring account information
- placing orders
- cancelling orders

### Important caution

The example instructions tell developers to modify a `PRIVATE_KEY` directly in example files before running them.

That is acceptable for sample code, but it is not acceptable for this repo.

### Planning implication

I should treat the Python SDK as a reference implementation for:

- request signing
- WebSocket payload formats
- account and order request shapes

I should not mirror its secret-handling pattern. If we ever use any of its logic, secrets must come from environment variables only.

---

## 7. Environment and Access Findings

Sources:

- https://test-app.pacifica.fi/
- https://pacifica.gitbook.io/docs/pacifica/testnet-guide
- https://pacifica.gitbook.io/docs/pacifica/close-beta-guide

### What is confirmed

- `Resources.md` points to the Pacifica testnet app and says to use code `Pacifica`
- Pacifica maintains separate testnet and mainnet environments
- the current docs still expose dedicated testnet WebSocket URLs

### Planning implication

The app should support environment switching from day one:

- Pacifica API base URL
- Pacifica WebSocket URL
- feature flags for testnet-only behavior

This should be config-driven, not hardcoded.

---

## 8. Sponsor Tool Fit

Sources:

- https://www.fuul.xyz/
- https://rhino.fi/
- https://rhino.fi/use-case-onchain-finance
- https://www.privy.io/
- https://www.elfa.ai/
- https://docs.elfa.ai/

### Privy

What it provides:

- embedded wallet infrastructure
- multiple login methods
- hardware-backed wallet/security features
- server-side and wallet orchestration tooling

Fit for Pacifica Pulse:

- strong option for onboarding non-technical users
- useful if the product later needs wallet-linked identity or delegated trading
- optional for MVP unless wallet-based participation is a hard requirement

### Fuul

What it provides:

- referral and affiliate tooling
- points and leaderboard programs
- SDK/API support for tracking and rewards
- no-code options for simpler campaigns

Fit for Pacifica Pulse:

- good fit for Phase 2 or Phase 3 gamification
- especially relevant if we add leaderboards, referrals, or trading competitions
- not necessary for core market creation and settlement

### Rhino.fi

What it provides:

- stablecoin onboarding
- cross-chain routing
- deposit widgets
- automated post-deposit actions into balances, vaults, or trading destinations

Fit for Pacifica Pulse:

- useful if the app later needs easy funding into trading or margin environments
- much more relevant to the trade-bridge phase than to the initial prediction-market MVP

### Elfa AI

What it provides:

- crypto-focused data and AI surfaces
- APIs for trending tokens, mentions, narratives, and event summaries
- API-key-based REST access

Fit for Pacifica Pulse:

- the strongest optional match for the AI market-intelligence layer in the readme
- can enrich market cards with social context without inventing fake signals
- should remain secondary to Pacifica market data, not replace it

---

## 9. Architecture Implications For The MVP

### Recommended system boundary

The cleanest MVP boundary is:

- Pacifica supplies external truth data
- Pacifica Pulse owns prediction-market state and user interaction

### Suggested component split

#### Frontend

- SvelteKit app
- market dashboard
- market creation flow
- market detail page
- countdown and resolution views

#### Backend

Given the repo rules, Go is the stronger default for the backend.

Why Go fits:

- good fit for WebSocket clients, workers, and scheduled settlement
- clear typed models for market state
- simple deployment story for a single API service

Tradeoff versus Fastify:

- Go gives stronger service/process boundaries
- Fastify would speed up shared TypeScript development with SvelteKit
- Go is still the better match here because the core difficulty is event ingestion and settlement reliability, not UI-adjacent JSON shaping

#### Storage

- PostgreSQL for markets, positions, snapshots, and settlement records
- Redis only if we later need fan-out caching or lightweight pub/sub

### Suggested MVP modules

- Pacifica REST client
- Pacifica WebSocket client
- market-definition service
- pricing and snapshot ingestion service
- settlement evaluator
- position service
- API handlers for the frontend

### Data shape additions worth planning early

The simplified readme model is a good start, but the MVP should probably also include:

- `market.condition_operator`
- `market.reference_value`
- `market.settlement_value`
- `market.resolved_at`
- `market.resolution_reason`
- `market.source_symbol`
- `market.source_interval`
- `position.price_or_odds_at_entry`
- `snapshot.source`
- `snapshot.observed_at`

These extra fields reduce ambiguity during settlement and demos.

---

## 10. What Should Not Be Built Yet

Based on the readme constraints and the source review, I should avoid these in Phase 1:

- a custom AMM
- full liquidity-pool mechanics
- real-money routing complexity
- trading execution as a hard dependency
- overbuilt multi-chain funding flows

These can be staged later once the prediction-market core is stable.

---

## 11. Main Risks I See Early

### Settlement ambiguity

If a market says "BTC above X in 30 minutes", we need one exact settlement source and one exact timestamp rule. Otherwise the app will look unfair.

### Source drift

The readme mentions price, funding, candle, and open interest markets, but each of those settles differently. The schema must encode the source and evaluation method explicitly.

Open interest is the sharpest example. The reviewed docs expose live OI, but I did not confirm a dedicated historical OI endpoint, so OI markets need either internal snapshotting or a later release boundary.

### Demo complexity

If the first build includes wallet auth, real trading, AI summaries, and gamification together, the demo will become noisy and fragile.

### External dependency concentration

Pacifica is both the live data source and the settlement truth. If connectivity degrades, the app needs graceful fallback states instead of freezing the UX.

---

## 12. Recommended Next Planning Step

The research points to a safe sequencing model:

1. lock the MVP market types
2. define the exact settlement rules per type
3. choose final backend architecture
4. map the database schema
5. define API contracts between frontend and backend
6. leave builder-code trading, Fuul growth loops, Rhino funding, and Elfa AI as isolated extensions

This keeps the first architecture focused on the product promise in `Readme.md` instead of trying to ship every sponsor path at once.
