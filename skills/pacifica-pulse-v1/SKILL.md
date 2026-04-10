---
name: pacifica-pulse-v1
description: Help design, implement, document, and review Pacifica Pulse or similar Pacifica-backed read-only prediction game systems. Use when the task involves Pacifica market data, guest-session gameplay, virtual balances, deterministic settlement, leaderboard derivation, SvelteKit frontend flow, or Go backend orchestration for prediction markets.
---

# Pacifica Pulse v1

Use this skill when the user is building or maintaining Pacifica Pulse, or when the target system is a very similar read-only prediction game on top of Pacifica testnet data.

## Keep the product boundary tight

- Treat the product as a prediction game, not a trading client.
- Keep guest accounts, virtual balances, and virtual payouts in v1.
- Exclude wallet authentication, withdrawals, real trading, and AI summaries unless the user explicitly expands scope.
- Keep Pacifica integration read-only by default.

## Use the repo's architecture defaults

- Use SvelteKit plus Tailwind CSS on the frontend.
- Use Go on the backend.
- Keep the frontend talking only to the backend, never directly to Pacifica.
- Keep balances, settlement, and market state server-authoritative.
- Prefer REST-first Pacifica integration unless a specific realtime requirement justifies more complexity.

## Preserve the supported market scope

- Support `price_threshold`, `candle_direction`, and `funding_threshold` markets in v1.
- Defer open-interest markets until internal snapshot support exists.
- Keep trade-flow and volume-derived markets out of v1 unless the user changes scope.

## Preserve market creation behavior

- Require the market creator to choose a side and stake during creation.
- Keep market creation as one transactional backend flow that creates the market and the creator's opening position together.
- Prefer one product-shaped backend route for create-plus-auto-stake behavior instead of splitting the flow across multiple frontend requests.
- Keep stake amounts as whole numbers unless the product rules change intentionally.

## Follow the settlement contract

- Resolve every market from one explicit Pacifica-derived source.
- Use a final authoritative REST fetch at or after expiry instead of guessing from stale in-memory data.
- Settle price markets from mark-price data.
- Settle candle markets from mark-price candles and only allow them when expiry lands on the selected candle boundary.
- Settle funding markets from the first funding record whose settlement timestamp is at or after expiry.
- If the fetched source timestamp still predates expiry, or Pacifica returns a temporary failure, retry instead of guessing.

## Preserve payout and leaderboard behavior

- Lock stake at entry time.
- Give losers no additional balance credit at settlement.
- Give winners fixed `2x stake` payouts in v1.
- Keep the leaderboard as a derived PostgreSQL read model instead of adding Redis or another cache layer unless scale changes materially.

## Use the expected realtime shape

- Prefer backend-owned SSE for the first realtime transport.
- Emit typed events such as `market.created`, `market.updated`, and `market.settled`.
- Publish realtime events only after durable writes succeed.

## Keep implementation reviews grounded

- When fixing a bug caused by a contract mismatch, query shape, or type mismatch, search adjacent routes, repositories, handlers, and sibling modules for the same pattern before stopping.
- Keep repo documentation synchronized with architecture changes.
- Prefer changes that improve correctness and maintainability over shortcuts that blur the product boundary.

## Sanity-check new work against these questions

- Does the change preserve the read-only Pacifica integration boundary?
- Does the backend still own balances, payouts, and settlement truth?
- Does the frontend still talk only to the backend?
- Does the new behavior stay within v1 market scope unless the user explicitly expands it?
- Does the change keep settlement deterministic and source-backed?
