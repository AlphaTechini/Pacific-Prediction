# backend/internal/market

## Purpose

I use this folder for market creation rules, market state transitions, and prediction-specific domain logic.

## Architectural Decisions And Tradeoffs

- I keep market rules separate from settlement execution because defining a market and resolving a market are related but not identical concerns.
- This package should own supported market types, threshold validation, and status transitions.
- Supported symbols should come from Pacifica market metadata instead of a hardcoded list, so symbol validation can track the actual tradable catalog without code churn.
- Market creation should normalize and validate input before persistence so the write path stays deterministic and later HTTP handlers remain thin.
- Market creation in v1 should create the market and the creator's first staked position in one transaction so the frontend does not need to coordinate two writes.
- Successful market creation should publish a backend-owned `market.created` event only after the transaction commits so subscribers never see phantom markets.
- Candle-direction markets should support the Pacifica mark-price candle intervals I expose in validation and still require expiry times that land exactly on the chosen candle close boundary so settlement resolves one unambiguous candle.
- Funding-threshold markets should anchor to the next hourly funding epoch after submission so they do not become due immediately and settlement always targets the next real checkpoint.
- Price-threshold markets should capture the live mark price at creation, store it as `reference_value`, use Pacifica `tick_size` as the primary threshold increment, and keep user thresholds inside a config-backed band on the correct side of that reference so trivial free-win markets are rejected.
- Creator opening stake validation should reject fractional stake amounts so market-funded balances stay whole-number based.
- Market creation context should come from Pacifica-backed symbol and price inputs so the UI can guide valid source choices before the user submits a market.
- Market listing should expose active and resolved catalogs without forcing the frontend to piece together multiple transport calls for the default dashboard view.
- Market detail should be readable through a single route-backed lookup so the frontend can hydrate a dedicated market page without reusing the list payload as a surrogate detail source.
- The tradeoff is another package boundary, but it makes the business rules easier to reason about.

## Logic Tracking

- To find market record and filter contracts visit [market.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/market.go).
- To find market creation input contracts visit [create_input.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/create_input.go).
- To find market creation context contracts visit [create_context.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/create_context.go).
- To find supported market validation models visit [validation_models.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/validation_models.go).
- To find dynamic symbol and market-shape validation rules visit [validator.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/validator.go).
- To find market service implementation visit [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/service_impl.go).
- To find price-threshold creation reference enrichment visit [create_reference.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/create_reference.go).
- To find market creation realtime publishing visit [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/service_impl.go).
- To find market controller implementation visit [controller_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/controller_impl.go).
- To find catalog response models visit [market.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/market.go).
- To find market service contracts visit [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/service.go).
- To find market controller contracts visit [controller.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/controller.go).
- To find settlement rule ownership visit [../settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find Pacifica market data sourcing visit [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).

## Component And Connection Map

- The market contract boundary can be found in [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/service.go) and [controller.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/controller.go).
- The supported market validation models can be found in [validation_models.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/validation_models.go).
- The Pacifica-backed symbol validation path can be found in [validator.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/validator.go).
- The Pacifica-backed market creation context can be found in [create_context.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/create_context.go) and [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/service_impl.go).
- The price-threshold creation reference lookup can be found in [create_reference.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/create_reference.go).
- The market creation orchestration can be found in [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/service_impl.go).
- The creator auto-stake market creation flow can be found in [create_input.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/create_input.go) and [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/service_impl.go).
- The market creation stream event handoff can be found in [../realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- The market persistence connection can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
