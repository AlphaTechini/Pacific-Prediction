# backend/internal/market

## Purpose

I use this folder for market creation rules, market state transitions, and prediction-specific domain logic.

## Architectural Decisions And Tradeoffs

- I keep market rules separate from settlement execution because defining a market and resolving a market are related but not identical concerns.
- This package should own supported market types, threshold validation, and status transitions.
- The tradeoff is another package boundary, but it makes the business rules easier to reason about.

## Logic Tracking

- To find market lifecycle logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find settlement rule ownership visit [../settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find Pacifica market data sourcing visit [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).

## Component And Connection Map

- The market domain can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- The market persistence connection can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
