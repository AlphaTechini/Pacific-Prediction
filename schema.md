# Schema

## Purpose

I am defining the database schema direction for Pacifica Pulse v1 so later migrations and repository code can be created confidently without re-deciding the model every time.

This document captures the schema choices, the reasons behind them, the tradeoffs I am accepting, and the extension path I expect to use later.

---

## 1. Chosen Schema Direction

I am choosing Option C as the schema strategy.

Option C is the middle ground between:

- an overly simple schema that ships fast but becomes painful to extend
- an overly normalized schema that is cleaner academically but too heavy for v1

Why I am choosing Option C:

- it is easy to build now
- it supports safe market settlement
- it keeps identity, balances, positions, and settlement audit separated
- it can grow mostly through additive tables instead of destructive redesign

This is the schema shape I can confidently use for v1 implementation.

---

## 2. Design Goals

The schema should optimize for:

1. correctness
2. clarity
3. maintainability
4. extension without rewrites

In practice that means:

- every market must be settleable from stored fields without guesswork
- balances must be server-authoritative
- guest identity must be persistent enough to survive page reloads and future auth upgrades
- settlement must be auditable
- I should avoid tables that only make sense for features we explicitly cut from v1

---

## 3. Core v1 Tables

### `players`

Purpose:

- persistent user identity for guest players today
- future anchor point for wallet auth or external auth later

Recommended fields:

- `id`
- `display_name`
- `created_at`
- `updated_at`

Reasoning:

- I do not want balance or market ownership tied directly to sessions.
- If wallet auth arrives later, I can attach new identity records to `players` instead of rebuilding ownership across the app.

Tradeoff:

- this adds one more table than a session-only prototype
- I accept that because it prevents identity sprawl later

### `player_sessions`

Purpose:

- guest login continuity with backend-controlled session state

Recommended fields:

- `id`
- `player_id`
- `session_token_hash`
- `expires_at`
- `created_at`

Reasoning:

- I only want to store a hash of the session token, not the raw token
- this supports revocation and safer server-side auth checks
- it is much safer than keeping the whole identity model in browser storage

Tradeoff:

- session persistence adds server state
- I accept that because local-only identity would be trivial to tamper with

### `player_balances`

Purpose:

- authoritative balance state for each player

Recommended fields:

- `player_id`
- `available_balance`
- `locked_balance`
- `updated_at`

Reasoning:

- `available_balance` tracks funds the player can still use
- `locked_balance` tracks funds reserved for unresolved markets
- this lets the backend prevent double-spending cleanly

Tradeoff:

- keeping separate balance columns means I need strict transactional updates
- I accept that because recalculating spendable balance from unresolved positions on every request would be slower and easier to get wrong

### `markets`

Purpose:

- stores both market definition and market lifecycle state

Recommended fields:

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

Reasoning:

- `market_type` keeps v1 market families explicit
- `source_type` describes the Pacifica-derived source used for settlement
- `source_interval` matters for candle and funding-style markets
- `reference_value` helps explain what the creator saw when the market was created
- `settlement_value` stores what actually resolved the market
- `resolution_reason` helps with auditability and edge cases

Why I am keeping definition and lifecycle in one table for v1:

- most reads need both anyway
- it keeps backend queries simple
- the product is not complex enough yet to justify splitting definition and result history into separate primary tables

Tradeoff:

- the `markets` table will carry more responsibility than a heavily normalized design
- I accept that because it is the right balance for v1 speed and clarity

### `positions`

Purpose:

- records each YES or NO entry for a player in a market

Recommended fields:

- `id`
- `player_id`
- `market_id`
- `side`
- `stake_amount`
- `potential_payout`
- `status`
- `created_at`
- `settled_at`

Reasoning:

- each row represents one participation event
- `stake_amount` captures player exposure
- `potential_payout` should be stored at entry time because v1 uses fixed odds
- `status` and `settled_at` make lifecycle handling explicit

Why storing `potential_payout` is important:

- it freezes the economic meaning of the position at entry time
- it protects old positions if payout rules ever change in a later version

Tradeoff:

- this duplicates some derived information
- I accept that because reproducible payout math matters more than strict normalization here

### `market_settlements`

Purpose:

- stores the settlement audit record for each resolved market

Recommended fields:

- `id`
- `market_id`
- `pacifica_source`
- `source_timestamp`
- `raw_payload`
- `settlement_value`
- `result`
- `created_at`

Reasoning:

- this is where I capture exactly what Pacifica-derived value was used
- `raw_payload` is valuable for debugging and dispute resolution
- the settlement worker should write this record in the same transactional flow as market and position resolution

Tradeoff:

- it introduces another write path during settlement
- I accept that because settlement is too important to leave unaudited

---

## 4. Optional v1 Table

### `market_snapshots`

Status:

- optional, not required for true v1

Potential fields:

- `id`
- `market_id`
- `symbol`
- `source_type`
- `observed_at`
- `mark_price`
- `funding_rate`
- `open_interest`
- `payload`

Why I am not forcing it into v1:

- price, candle, and funding markets can settle without it
- I do not want to add storage and retention complexity until I actually need historical replay or chart detail

Why it is still worth planning for:

- it is the cleanest bridge into open-interest markets later
- it helps with charting and debugging if richer market detail becomes important

---

## 5. What I Am Intentionally Not Modeling Yet

I am not creating dedicated schema for these in v1:

- wallets
- withdrawals
- reward redemptions
- leaderboards
- dynamic odds history
- liquidity pools
- real trade execution
- trade-flow aggregations
- volume-derived analytics

Reasoning:

- none of these are required for the agreed v1
- adding them now would make the schema look more complete while actually making implementation harder

This is deliberate scope control, not missing design.

---

## 6. Enum-Like Field Strategy

I do not need PostgreSQL enum types immediately to be confident in the schema design.

For v1, the important part is agreeing on the allowed values first.

### `market_type`

Recommended values:

- `price_threshold`
- `candle_direction`
- `funding_threshold`

### `source_type`

Recommended values:

- `mark_price`
- `mark_price_candle`
- `funding_rate`

### `condition_operator`

Recommended values:

- `gt`
- `gte`
- `lt`
- `lte`
- `bullish_close`
- `bearish_close`
- `positive`
- `negative`

### `market.status`

Recommended values:

- `active`
- `resolving`
- `resolved`
- `cancelled`

### `positions.status`

Recommended values:

- `open`
- `won`
- `lost`
- `cancelled`

Why I like this approach:

- it keeps the domain vocabulary explicit
- it makes validation easier in both the API and the database layer
- it leaves room to convert these into database enums later if I want stronger constraints

---

## 7. Indexing Strategy I Should Plan Early

I should plan these indexes as part of the first migration set:

- `markets(status, expiry_time)`
- `positions(player_id, created_at)`
- `positions(market_id, created_at)`
- `player_sessions(session_token_hash)`
- `market_settlements(market_id)`

Why these first:

- settlement workers need fast lookup of expiring markets
- player views need fast lookup of a player’s positions
- market detail views need fast lookup of market positions
- session lookup needs to be fast and exact
- settlement audit lookup should be cheap when showing resolved market detail

---

## 8. Why Option C Is Easy To Extend

This schema scales well because most future features can be added with new tables rather than rewiring the existing ones.

### Easy additions later

#### Wallet auth

Likely additive tables:

- `player_wallets`
- `auth_identities`

Why this is easy:

- `players` already exists as the canonical owner record

#### Reward systems or token-equivalent mapping

Likely additive tables:

- `reward_ledger`
- `reward_redemptions`

Why this is easy:

- balances and settlements already exist as clean event anchors

#### Leaderboards

Likely approach:

- derived query layer first
- table or materialized view later only if needed

Why this is easy:

- leaderboard data comes naturally from balances, positions, and settlements

#### Open interest markets

Likely additive table:

- `market_snapshots`

Why this is easy:

- the core market and position model does not need to change much
- I mainly need snapshot capture plus new settlement logic

---

## 9. Where Option C Would Eventually Need More Structure

I would revisit the schema if the product grows into any of these:

- dynamic odds and odds history
- pooled market mechanics or AMM-style liquidity
- highly custom multi-condition markets
- large-scale analytics or event-sourcing needs
- full real-trading integration with order lifecycle tracking

Why these are different:

- they introduce new domains, not just new fields
- forcing them into the current six-table core would eventually create overloaded tables and confusing write paths

That is not a problem today. It just means Option C is the right v1 and v2 shape, not necessarily the forever shape.

---

## 10. Practical Implementation Guidance

When I later write PostgreSQL migrations, I should follow these rules:

- use additive migrations first
- keep IDs stable and opaque
- store monetary-like virtual values with exact numeric types, not floating point
- keep balance updates and settlement writes inside transactions
- avoid embedding business rules only in SQL when they are easier to reason about in Go

I also want migration naming and schema changes to stay tightly aligned with this document so the implementation never drifts silently away from the design.

---

## 11. Final Schema Position

I am confident using this schema as the backend foundation for v1.

The core reason is simple:

- it is small enough to move fast
- it is explicit enough to be safe
- it is structured enough to grow later

That is exactly what this project needs right now.
