# Tasks

## Purpose

I am turning the architecture and schema decisions into an execution-ready task list so implementation can proceed one task at a time without losing dependency order.

This file is intentionally organized by phase first, then by module-level work inside each phase.

---

## How I Will Use This File

Each task is designed to be implementation-ready.

Every task includes:

- a clear outcome
- the primary owning module
- dependencies
- concrete deliverables
- acceptance criteria

Task statuses should be updated as work progresses:

- `todo`
- `in_progress`
- `done`
- `blocked`

---

## Execution Strategy

I am using `phase -> module tasks inside each phase` instead of a pure module-by-module plan.

Why:

- it respects dependency order
- it avoids starting modules before their contracts exist
- it keeps the backend shippable in slices
- it reduces later rework across auth, balance, position, and settlement

---

## Phase 1: Foundations

Goal:

- make the backend boot
- define the shared boundaries
- prepare the project for later module implementation

### T1.1 Initialize Go Backend Project`r`n`r`nStatus:`r`n`r`n- `done`

Owner:

- `backend`

Dependencies:

- none

Deliverables:

- Go module initialized in `backend`
- baseline folder layout aligned with `structure.md`
- starter `main.go`
- baseline `.env.example`
- backend `.gitignore` review for env safety

Acceptance criteria:

- the Go module builds successfully
- `main.go` starts a minimal app without business logic
- no secrets are hardcoded

### T1.2 Add Runtime Configuration Layer`r`n`r`nStatus:`r`n`r`n- `done`

Owner:

- `config`

Dependencies:

- T1.1

Deliverables:

- config struct definitions
- environment loading logic
- startup validation for required variables

Acceptance criteria:

- startup fails fast when required config is missing
- all runtime config is loaded through one config package

### T1.3 Add Database Connection And Transaction Base`r`n`r`nStatus:`r`n`r`n- `done`

Owner:

- `storage`

Dependencies:

- T1.2

Deliverables:

- PostgreSQL connection setup
- transaction manager or transaction helper
- connection lifecycle wiring in `main.go`

Acceptance criteria:

- the app can connect to PostgreSQL with environment config
- connection startup and shutdown paths are explicit

### T1.4 Define Shared Domain Constants And Common Types

Status:

- `done`

Owner:

- `backend`

Dependencies:

- T1.1

Deliverables:

- common enum-like constants for market and position state
- shared ID and timestamp handling decisions
- common error patterns for domain validation

Acceptance criteria:

- supported domain values match `schema.md` and `architecture.md`
- modules can reuse common domain constants instead of redefining them

### T1.5 Define Repository Interface Contracts`r`n`r`nStatus:`r`n`r`n- `done`

Owner:

- `storage`

Dependencies:

- T1.3
- T1.4

Deliverables:

- repository interfaces per module boundary
- transaction-safe repository usage pattern

Acceptance criteria:

- repository contracts exist for auth, player, balance, market, position, and settlement concerns
- interfaces align with the schema and module boundaries already agreed

### T1.6 Define Service And Controller Wiring Pattern

Status:

- `todo`

Owner:

- `backend`

Dependencies:

- T1.2
- T1.5

Deliverables:

- dependency injection pattern for services and controllers
- route registration approach for `main.go`
- controller export conventions

Acceptance criteria:

- `main.go` remains the composition root
- modules expose controllers or handlers, not route ownership

### T1.7 Add Backend Dockerfile`r`n`r`nStatus:`r`n`r`n- `done`

Owner:

- `backend`

Dependencies:

- T1.1

Deliverables:

- backend `Dockerfile`

Acceptance criteria:

- the backend can be containerized with a reproducible build path

---

## Phase 2: Auth, Player, And Balance Core

Goal:

- make guest identity and virtual balance work first

### T2.1 Implement Auth Domain Contracts`r`n`r`nStatus:`r`n`r`n- `done`

Owner:

- `auth`

Dependencies:

- T1.4
- T1.5
- T1.6

Deliverables:

- session schemas
- auth service interface
- auth controller interface

Acceptance criteria:

- auth contracts cover guest session creation and validation
- auth remains separate from player identity ownership

### T2.2 Implement Player Domain Contracts`r`n`r`nStatus:`r`n`r`n- `done`

Owner:

- `player`

Dependencies:

- T1.4
- T1.5
- T1.6

Deliverables:

- player schemas
- player service interface
- player controller interface

Acceptance criteria:

- player contracts cover identity and profile reads only
- no balance or session ownership leaks into player

### T2.3 Implement Balance Domain Contracts`r`n`r`nStatus:`r`n`r`n- `done`

Owner:

- `balance`

Dependencies:

- T1.4
- T1.5
- T1.6

Deliverables:

- balance schemas
- balance service interface
- rules for debit, credit, lock, and unlock operations

Acceptance criteria:

- balance contracts support authoritative available and locked balance flows
- spendable-value mutation rules are centralized in balance

### T2.4 Create Initial Identity And Balance Migrations`r`n`r`nStatus:`r`n`r`n- `done`

Owner:

- `storage`

Dependencies:

- T2.1
- T2.2
- T2.3

Deliverables:

- migrations for `players`
- migrations for `player_sessions`
- migrations for `player_balances`

Acceptance criteria:

- migrations align with `schema.md`
- session tokens are stored as hashes, not raw values

### T2.5 Implement Guest Session Creation Flow

Status:

- `todo`

Owner:

- `auth`

Dependencies:

- T2.1
- T2.4

Deliverables:

- `POST /api/v1/players/guest`
- guest player creation or guest session issuance flow
- secure cookie behavior

Acceptance criteria:

- a guest session can be created from the API
- session lookup works on subsequent requests

### T2.6 Implement Current Player Profile Flow

Status:

- `todo`

Owner:

- `player`

Dependencies:

- T2.2
- T2.5

Deliverables:

- `GET /api/v1/players/me`

Acceptance criteria:

- the endpoint returns the authenticated guest player profile

### T2.7 Implement Current Balance Flow

Status:

- `todo`

Owner:

- `balance`

Dependencies:

- T2.3
- T2.4
- T2.5

Deliverables:

- `GET /api/v1/players/me/balance`

Acceptance criteria:

- the endpoint returns authoritative available and locked balances

### T2.8 Add Session Resolution Middleware

Status:

- `todo`

Owner:

- `auth`

Dependencies:

- T2.5

Deliverables:

- request auth-context resolution
- middleware or equivalent request binding

Acceptance criteria:

- protected handlers can access the authenticated player context consistently

---

## Phase 3: Market Core

Goal:

- allow market creation and market reads

### T3.1 Implement Market Domain Contracts

Status:

- `todo`

Owner:

- `market`

Dependencies:

- T1.4
- T1.5
- T1.6
- T2.8

Deliverables:

- market schemas
- market service interface
- market controller interface
- validation models for supported market types

Acceptance criteria:

- contracts support price threshold, candle direction, and funding threshold markets only

### T3.2 Create Market Migrations

Status:

- `todo`

Owner:

- `storage`

Dependencies:

- T3.1

Deliverables:

- migration for `markets`

Acceptance criteria:

- the market table stores enough information for deterministic settlement

### T3.3 Implement Market Validation Rules

Status:

- `todo`

Owner:

- `market`

Dependencies:

- T3.1

Deliverables:

- validation for supported symbols
- validation for supported market types
- validation for threshold and expiry inputs

Acceptance criteria:

- unsupported market types are rejected
- invalid expiry and threshold combinations are rejected

### T3.4 Implement Create Market Flow

Status:

- `todo`

Owner:

- `market`

Dependencies:

- T3.2
- T3.3
- T2.8

Deliverables:

- `POST /api/v1/markets`

Acceptance criteria:

- an authenticated guest can create a valid market
- persisted market state matches agreed schema

### T3.5 Implement List Markets Flow

Status:

- `todo`

Owner:

- `market`

Dependencies:

- T3.2

Deliverables:

- `GET /api/v1/markets`

Acceptance criteria:

- active and resolved market listing is available

### T3.6 Implement Market Detail Flow

Status:

- `todo`

Owner:

- `market`

Dependencies:

- T3.2

Deliverables:

- `GET /api/v1/markets/:id`

Acceptance criteria:

- a single market can be retrieved with enough data for frontend display

---

## Phase 4: Position Core

Goal:

- allow safe YES and NO participation with locked virtual balances

### T4.1 Implement Position Domain Contracts

Status:

- `todo`

Owner:

- `position`

Dependencies:

- T1.4
- T1.5
- T1.6
- T2.3
- T3.1

Deliverables:

- position schemas
- position service interface
- position controller interface

Acceptance criteria:

- contracts support creating and listing player positions

### T4.2 Create Position Migrations

Status:

- `todo`

Owner:

- `storage`

Dependencies:

- T4.1

Deliverables:

- migration for `positions`

Acceptance criteria:

- the position table supports fixed-odds v1 participation

### T4.3 Implement Position Placement Rules

Status:

- `todo`

Owner:

- `position`

Dependencies:

- T4.1
- T3.3
- T2.3

Deliverables:

- validation for open market state
- validation for stake amount
- validation for sufficient available balance

Acceptance criteria:

- invalid positions are rejected before persistence

### T4.4 Implement Place Position Flow

Status:

- `todo`

Owner:

- `position`

Dependencies:

- T4.2
- T4.3
- T2.8

Deliverables:

- `POST /api/v1/markets/:id/positions`
- balance locking integration

Acceptance criteria:

- creating a position locks the player stake correctly
- the stored position includes the correct potential payout

### T4.5 Implement Player Position History Flow

Status:

- `todo`

Owner:

- `position`

Dependencies:

- T4.2
- T2.8

Deliverables:

- `GET /api/v1/players/me/positions`

Acceptance criteria:

- the player can retrieve their positions across open and resolved markets

---

## Phase 5: Pacifica Integration

Goal:

- connect the backend to live and settlement-ready Pacifica data

### T5.1 Implement Pacifica REST Client Contracts

Status:

- `todo`

Owner:

- `pacifica`

Dependencies:

- T1.2
- T1.4

Deliverables:

- REST client interface
- normalized response models for market info, prices, candles, and funding

Acceptance criteria:

- Pacifica REST concerns are isolated from domain modules

### T5.2 Implement Pacifica WebSocket Client Contracts

Status:

- `todo`

Owner:

- `pacifica`

Dependencies:

- T1.2
- T1.4

Deliverables:

- WebSocket client interface
- normalized live event models

Acceptance criteria:

- live Pacifica event handling is isolated behind the pacifica module

### T5.3 Implement REST Data Fetches

Status:

- `todo`

Owner:

- `pacifica`

Dependencies:

- T5.1

Deliverables:

- market metadata fetch
- settlement-ready price fetch path
- candle fetch path
- funding fetch path

Acceptance criteria:

- required read-only Pacifica endpoints work through the backend integration layer

### T5.4 Implement WebSocket Subscription Manager

Status:

- `todo`

Owner:

- `pacifica`

Dependencies:

- T5.2

Deliverables:

- subscription manager
- heartbeat handling
- reconnect handling

Acceptance criteria:

- Pacifica live subscriptions can recover from disconnects without crashing the app

### T5.5 Integrate Pacifica Validation Inputs Into Market Module

Status:

- `todo`

Owner:

- `market`

Dependencies:

- T5.3
- T3.3

Deliverables:

- symbol validation against Pacifica metadata
- source-aware market creation context

Acceptance criteria:

- market creation can validate supported Pacifica symbols safely

---

## Phase 6: Settlement Engine

Goal:

- resolve markets deterministically and update payouts safely

### T6.1 Implement Settlement Domain Contracts

Status:

- `todo`

Owner:

- `settlement`

Dependencies:

- T1.4
- T1.5
- T1.6
- T3.1
- T4.1
- T5.1

Deliverables:

- settlement schemas
- settlement service interface
- worker interface

Acceptance criteria:

- settlement contracts cover price, candle, and funding resolution paths

### T6.2 Create Settlement Migrations

Status:

- `todo`

Owner:

- `storage`

Dependencies:

- T6.1

Deliverables:

- migration for `market_settlements`

Acceptance criteria:

- settlement audit storage aligns with `schema.md`

### T6.3 Implement Expiry Scanner

Status:

- `todo`

Owner:

- `settlement`

Dependencies:

- T6.1
- T3.2

Deliverables:

- scheduled or loop-driven expiry scan flow

Acceptance criteria:

- expiring markets can be discovered reliably

### T6.4 Implement Price Market Settlement

Status:

- `todo`

Owner:

- `settlement`

Dependencies:

- T6.3
- T5.3

Deliverables:

- price threshold settlement logic

Acceptance criteria:

- price markets resolve from the agreed Pacifica-derived source and timestamp rule

### T6.5 Implement Candle Market Settlement

Status:

- `todo`

Owner:

- `settlement`

Dependencies:

- T6.3
- T5.3

Deliverables:

- candle direction settlement logic

Acceptance criteria:

- candle markets resolve from mark-price candle data only

### T6.6 Implement Funding Market Settlement

Status:

- `todo`

Owner:

- `settlement`

Dependencies:

- T6.3
- T5.3

Deliverables:

- funding threshold settlement logic

Acceptance criteria:

- funding markets resolve from the agreed Pacifica funding record rule

### T6.7 Implement Transactional Payout Application

Status:

- `todo`

Owner:

- `settlement`

Dependencies:

- T6.4
- T6.5
- T6.6
- T2.3
- T4.2

Deliverables:

- market status update flow
- position win/loss updates
- balance unlock and credit/debit updates
- settlement audit write path

Acceptance criteria:

- settlement completes as one consistent transactional operation

### T6.8 Add Settlement Retry And Failure Handling

Status:

- `todo`

Owner:

- `settlement`

Dependencies:

- T6.7

Deliverables:

- retry behavior for Pacifica fetch failures
- safe handling for temporary data unavailability

Acceptance criteria:

- the system delays or retries uncertain settlement instead of guessing

---

## Phase 7: Realtime

Goal:

- push useful live updates to the frontend

### T7.1 Implement Realtime Event Contracts

Status:

- `todo`

Owner:

- `realtime`

Dependencies:

- T1.4

Deliverables:

- stream event schemas
- event type definitions

Acceptance criteria:

- event types cover market updates and settlement updates needed in v1

### T7.2 Implement SSE Stream Endpoint

Status:

- `todo`

Owner:

- `realtime`

Dependencies:

- T7.1
- T1.6

Deliverables:

- `GET /api/v1/stream`

Acceptance criteria:

- clients can subscribe to a live backend event stream

### T7.3 Implement Event Publisher And Hub

Status:

- `todo`

Owner:

- `realtime`

Dependencies:

- T7.1
- T7.2

Deliverables:

- publisher interface
- in-process event hub

Acceptance criteria:

- backend modules can publish events without directly managing client connections

### T7.4 Publish Market And Settlement Events

Status:

- `todo`

Owner:

- `realtime`

Dependencies:

- T3.4
- T4.4
- T6.7
- T7.3

Deliverables:

- market creation events
- market state update events
- settlement result events

Acceptance criteria:

- subscribed clients receive relevant v1 updates

---

## Phase 8: Hardening

Goal:

- make the backend more reliable and safer for demo use

### T8.1 Add Request Validation Consistency

Status:

- `todo`

Owner:

- `httpapi`

Dependencies:

- T2.5
- T3.4
- T4.4

Deliverables:

- shared request validation pattern
- consistent error response behavior

Acceptance criteria:

- API validation failures are predictable and clear

### T8.2 Add Basic Rate Limiting

Status:

- `todo`

Owner:

- `httpapi`

Dependencies:

- T2.5
- T3.4
- T4.4

Deliverables:

- rate limiting for guest session creation
- rate limiting for market creation
- rate limiting for position placement

Acceptance criteria:

- spam paths are protected at the transport layer

### T8.3 Add Logging And Operational Visibility

Status:

- `todo`

Owner:

- `backend`

Dependencies:

- T1.6
- T6.7

Deliverables:

- request logging baseline
- worker logging baseline
- Pacifica connectivity logging baseline

Acceptance criteria:

- failures in startup, Pacifica integration, and settlement are visible in logs

### T8.4 Add Pacifica Delay Handling

Status:

- `todo`

Owner:

- `pacifica`

Dependencies:

- T5.4
- T6.8

Deliverables:

- delayed-data detection
- stale-data signaling for app flows

Acceptance criteria:

- the backend can detect and surface delayed upstream market data conditions

### T8.5 Add Readiness And Health Checks

Status:

- `todo`

Owner:

- `backend`

Dependencies:

- T1.3
- T5.3

Deliverables:

- health endpoint strategy
- readiness checks for DB and core dependencies

Acceptance criteria:

- the app can expose whether it is ready to serve requests safely

---

## Recommended First Execution Order

I should execute the first implementation slice in this order:

1. T1.1 Initialize Go Backend Project
2. T1.2 Add Runtime Configuration Layer
3. T1.3 Add Database Connection And Transaction Base
4. T1.4 Define Shared Domain Constants And Common Types
5. T1.5 Define Repository Interface Contracts
6. T1.6 Define Service And Controller Wiring Pattern
7. T1.7 Add Backend Dockerfile
8. T2.1 Implement Auth Domain Contracts
9. T2.2 Implement Player Domain Contracts
10. T2.3 Implement Balance Domain Contracts

That gives me the cleanest possible base before touching migrations and route behavior.


