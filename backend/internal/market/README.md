# backend/internal/market

## Purpose

I use this folder for market creation rules, market state transitions, and prediction-specific domain logic.

## Architectural Decisions And Tradeoffs

- I keep market rules separate from settlement execution because defining a market and resolving a market are related but not identical concerns.
- This package should own supported market types, threshold validation, and status transitions.
- The tradeoff is another package boundary, but it makes the business rules easier to reason about.

## Logic Tracking

- To find market record and filter contracts visit [market.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/market.go).
- To find market creation input contracts visit [create_input.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/create_input.go).
- To find supported market validation models visit [validation_models.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/validation_models.go).
- To find market service contracts visit [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/service.go).
- To find market controller contracts visit [controller.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/controller.go).
- To find settlement rule ownership visit [../settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find Pacifica market data sourcing visit [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).

## Component And Connection Map

- The market contract boundary can be found in [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/service.go) and [controller.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/controller.go).
- The supported market validation models can be found in [validation_models.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/validation_models.go).
- The market persistence connection can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
