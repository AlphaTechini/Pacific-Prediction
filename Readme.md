# Pacifica Pulse

## Documentation

- [structure.md](./structure.md)
- [architecture.md](./architecture.md)
- [Research.md](./Research.md)
- [schema.md](./schema.md)
- [task.md](./task.md)

## Overview

Pacifica Pulse is a real-time prediction market overlay built on top of Pacifica’s perpetual trading infrastructure.

It allows users to create, participate in, and settle short-duration prediction markets using live Pacifica market data such as mark price, funding rate, open interest, and candle data.

The system is designed to transform raw trading data into actionable, gamified market insights while maintaining a clear path toward real trading integration.

---

## Why This Project Exists

Perpetual traders rely heavily on interpreting market signals such as price action, funding rates, and liquidity flow.

However:

* Signals are fragmented across tools
* Decision-making is manual and slow
* There is no native way to express or trade *beliefs* about short-term outcomes

Pacifica Pulse solves this by:

* Turning market expectations into structured prediction markets
* Using real Pacifica data for settlement
* Providing optional pathways to convert predictions into trading actions

---

## Core Features

### 1. Prediction Market Creation

Users can create markets based on predefined or custom conditions:

Examples:

* “Will BTC be above $105k in 30 minutes?”
* “Will ETH funding remain positive for the next interval?”
* “Will SOL close the next 5m candle bullish?”

Market Parameters:

* Asset (BTC, ETH, etc.)
* Condition type (price, funding, candle, OI)
* Threshold
* Expiry time
* Settlement source (mark price, candle close, etc.)

---

### 2. Real-Time Market Data Integration

Powered by Pacifica APIs:

* WebSocket price stream (mark price, funding, OI, volume)
* Mark price candles (1m–1d)
* REST historical data
* Recent trades

Used for:

* Market creation context
* Live UI updates
* Settlement logic

---

### 3. Market Participation

Users can:

* Take YES or NO positions
* View live odds (simple pooled model for MVP)
* Track active and resolved markets

---

### 4. Settlement Engine

Markets are resolved using Pacifica data:

Examples:

* Price-based → mark price at expiry
* Candle-based → close price of selected interval
* Funding-based → funding value at interval

Settlement flow:

1. Fetch final data point from Pacifica
2. Evaluate condition
3. Resolve outcome (YES/NO)
4. Update market + user positions

---

### 5. AI Market Intelligence Layer (Optional but Recommended)

AI is used to enhance understanding, not fabricate signals.

Capabilities:

* Generate market summaries using real data
* Explain why a market is interesting
* Classify market types (trend, mean-reversion, volatility, squeeze)

Example Output:
“BTC is rising with increasing open interest and positive funding, suggesting aggressive long positioning rather than short covering.”

---

### 6. Signal-to-Trade Bridge (Optional Extension)

Users can convert prediction outcomes into trading actions:

* “Trade this idea on Pacifica”
* Backend executes order using Pacifica API
* Supports future integration with API Agent Keys

---

## System Architecture

### High-Level Components

1. Frontend (SvelteKit + Tailwind)

   * Market creation UI
   * Live market dashboard
   * Real-time updates via WebSocket

2. Backend API (Go or Fastify)

   * Market lifecycle management
   * Settlement engine
   * AI integration layer
   * Pacifica API integration

3. Data Layer

   * PostgreSQL (markets, users, positions)
   * Redis (optional for real-time caching)

4. Pacifica Integration Layer

   * WebSocket client for live data
   * REST client for historical + settlement
   * Optional trading execution

---

## Data Model (Simplified)

### Market

* id
* symbol
* condition_type
* threshold
* expiry_time
* settlement_source
* status (active, resolved)
* result (yes, no)

### Position

* id
* user_id
* market_id
* side (yes/no)
* amount

### MarketSnapshot

* market_id
* timestamp
* mark_price
* funding_rate
* open_interest

---

## Market Types (MVP Scope)

1. Price-Based
2. Candle-Based
3. Funding-Based
4. Open Interest-Based (optional)

---

## API Design (Internal)

### POST /markets

Create a new market

### GET /markets

List active markets

### GET /markets/:id

Get market details

### POST /markets/:id/position

Place a YES/NO position

### POST /markets/:id/settle

Trigger settlement (cron or worker)

---

## Pacifica Integration

### WebSocket

* Subscribe to price stream
* Subscribe to mark price candles

### REST

* Fetch historical candles
* Fetch recent trades
* Fetch settlement data at expiry

### Optional

* Order execution endpoints
* API Agent Keys for delegated trading

---

## Development Plan

### Phase 1 (Core MVP)

* Market creation
* Real-time data ingestion
* Basic participation (YES/NO)
* Settlement engine

### Phase 2

* AI summaries
* Market analytics
* Improved UI/UX

### Phase 3

* Trade execution bridge
* Agent key integration
* Leaderboards / gamification

---

## Constraints

* Must use Pacifica APIs as primary data source
* Must support testnet environment
* Must prioritize demo clarity over complexity
* No overengineering of AMM or liquidity pools for MVP

---

## Future Expansion

* Support for additional asset classes (e.g., stocks)
* Advanced market types (multi-condition, range-based)
* Strategy backtesting engine
* Social features (copy markets, leaderboards)
* Automated trading strategies based on prediction outcomes

---

## Tech Stack

* Frontend: SvelteKit + Tailwind
* Backend: Go or Fastify (Node.js)
* Database: PostgreSQL
* Realtime: WebSocket
* AI: LLM via API (prompt-driven, no fine-tuning)

---

## Demo Flow

1. User opens dashboard
2. Creates a prediction market
3. System shows real-time Pacifica data
4. Users take positions
5. Countdown to expiry
6. Market settles using Pacifica data
7. Results displayed
8. Optional: execute trade based on outcome

---

## Key Differentiator

Pacifica Pulse is not just a prediction market.

It is a **market intelligence layer** that:

* uses real trading data
* enables structured belief expression
* connects sentiment to execution




