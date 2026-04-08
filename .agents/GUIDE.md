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
- Pacifica integration is read-only in v1 through REST and WebSocket.
- The frontend should talk only to our backend, not directly to Pacifica.
- Player identity in v1 uses guest accounts with backend-issued sessions.
- Balances must be server-authoritative.

## Settlement Decisions

- Every market must resolve from one explicit Pacifica-derived source.
- Price markets settle from mark-price data.
- Candle markets settle from mark-price candles.
- Funding markets settle from funding records.
- I should prefer a final authoritative REST fetch for settlement instead of relying only on in-memory live stream data.

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

