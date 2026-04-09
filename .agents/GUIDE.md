# Project Guide

## Confirmed Constraints

- Pacifica Pulse v1 is a read-only prediction game on top of Pacifica testnet data.
- v1 does not include AI summaries.
- v1 does not include real trading.
- v1 does not include wallet authentication.
- v1 does not include withdrawals.
- v1 uses virtual balances and virtual payouts only.
- Pacifica APIs are required, but wallet integration is not required in v1.
- Testnet support is required.

## Confirmed Market Scope

- Price threshold markets are in v1.
- Candle outcome markets are in v1.
- Funding direction or threshold markets are in v1.
- Open interest markets are deferred to a later version that includes internal snapshot support.
- Trade-flow markets are out of scope for v1.
- Volume-derived markets are out of scope for v1.

## Architecture Decisions

- Frontend stack: SvelteKit + Tailwind.
- Backend stack: Go.
- Database: PostgreSQL.
- Redis is not part of v1.
- Pacifica integration is read-only in v1 and should stay REST-first by default.
- The frontend should talk only to our backend, not directly to Pacifica.
- Player identity in v1 uses guest accounts with backend-issued sessions.
- Balances must be server-authoritative.
- Market creators must choose a side and stake when creating a market, and that initial creator participation should be handled by the backend as part of one flow.
- The frontend should prefer one backend route for "create market with creator auto-stake" instead of coordinating separate market-create and position-create requests.

## Settlement Decisions

- Every market must resolve from one explicit Pacifica-derived source.
- Price markets settle from mark-price data.
- Candle markets settle from mark-price candles.
- Funding markets settle from funding records.
- I should prefer a final authoritative REST fetch for settlement instead of relying only on in-memory live stream data.
- Price-threshold settlement should use one batched REST price request at expiry time and only settle when Pacifica's returned timestamp is at or after the market expiry.
- If the returned price snapshot still predates expiry, the settlement worker should retry briefly instead of guessing.
- Candle markets should fetch historical mark-price candles on demand at resolution time rather than polling continuously.
- Candle markets should only be valid when `expiry_time` lands exactly on the selected candle interval boundary so settlement can resolve one explicit finished candle.
- Funding markets should fetch historical funding records on demand at resolution time rather than polling continuously.
- Funding markets should resolve from the first Pacifica funding record whose settlement timestamp is at or after the market expiry.
- Always-on Pacifica WebSocket subscriptions are not the default settlement path in v1.
- Stake placement removes spendable balance up front.
- On settlement, losers receive no further balance change and winners receive their fixed `potential_payout`.
- `POST /api/v1/markets` should own the transactional creator flow instead of splitting creator market creation and creator opening position into separate frontend steps.

## Deferred Roadmap

- v1.5 candidate: open interest markets backed by internal snapshot capture.
- v2 candidate: wallet auth and optional reward redemption or testnet token mapping.
- v3 candidate: real trade bridge, builder code support, and API Agent Keys.

## Schema Decisions

- I am using Option C as the database strategy.
- The v1 core schema uses six primary tables: players, player_sessions, player_balances, markets, positions, and market_settlements.
- market_snapshots stays optional in true v1 and becomes the clean bridge to OI markets later.
- The schema is intended to scale through additive tables for wallets, rewards, leaderboards, and OI support rather than destructive redesign.


## Module Decisions

- The backend module split is auth, player, balance, market, position, settlement, pacifica, realtime, config, and storage.
- The backend also uses a shared internal domain package for enum-like values, ID aliases, UTC timestamp normalization, and reusable validation errors.
- main.go owns explicit route registration.
- Modules export controllers or handlers, while main.go constructs dependencies and assigns routes.
- auth owns sessions, player owns identity, and balance owns spendable-value mutations.
- Storage repository contracts are split by module concern and should be consumed through a transaction-scoped repository provider instead of ad hoc SQL access.
- The HTTP wiring pattern uses an `httpapi.Application` container for shared dependencies and controller references, plus an `httpapi.Router` helper that main.go uses to register method-aware routes explicitly.
- Guest sessions use opaque random tokens stored only in secure cookies, while PostgreSQL stores the token hash and the initial virtual balance is provisioned during guest creation from environment-backed config.
- Market validation should use Pacifica `/api/v1/info` metadata through the pacifica module, with backend-owned caching instead of a hardcoded symbol list.

## Settlement Contract Decisions

- T6.1 excludes payout application and defers payout concerns to T6.7.
- Settlement keeps one shared orchestration service and one worker contract.
- Price, candle, and funding settlement use separate resolver interfaces.
- Price, candle, and funding settlement use separate output schemas because their settlement truth differs by source and timestamp semantics.
- Any shared settlement audit shape should exist only after resolver-specific outputs are mapped into persistence-friendly fields.
- The expiry scanner uses the market repository to discover due active markets by `expiry_time` and hands off only orchestration-level attempts until per-market settlement logic is implemented.
- Price settlement should group due markets into batched Pacifica price fetches instead of creating per-market cron jobs.
- The settlement worker should maintain a configurable near-expiry lookahead window and retry interval for price-fetch planning.
- T6.4 resolves due price-threshold markets through the settlement service, validates Pacifica timestamps before expiry, writes a settlement audit row, and marks the market resolved in one transaction.

