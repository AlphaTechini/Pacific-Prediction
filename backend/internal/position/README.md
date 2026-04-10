# backend/internal/position

## Purpose

I use this folder for direct participant position placement and player-facing position history reads.

## Architectural Decisions And Tradeoffs

- I keep position contracts in their own package because position placement depends on market eligibility and balance locking, but it should not collapse those concerns into one module.
- I keep direct participant placement here even though the market module now handles the creator's opening position during market creation.
- I keep `potential_payout` on the position record because the architecture locks fixed-odds economics at entry time instead of recomputing old positions from newer payout rules.
- Position stake validation should reject fractional stake amounts so spendable balances and payouts stay easy to manage as whole-number values.
- Successful position placement should publish a backend-owned `market.updated` event after commit so subscribers can react to market activity without this package managing SSE details.
- The tradeoff is another package boundary, but it keeps position lifecycle rules easier to change without destabilizing auth, market, or balance code.

## Logic Tracking

- To find position record and request contracts visit [contracts.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/contracts.go).
- To find position ID generation visit [ids.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/ids.go).
- To find position validation rules visit [validator.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/validator.go).
- To find position service contracts visit [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/service.go).
- To find position service implementation visit [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/service_impl.go).
- To find position-placement realtime publishing visit [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/service_impl.go).
- To find position controller contracts visit [controller.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/controller.go).
- To find position controller implementation visit [controller_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/controller_impl.go).
- To find shared position enums and ID aliases visit [../domain/position.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/position.go) and [../domain/ids.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/ids.go).
- To find market eligibility rules that this package will depend on visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find balance locking rules that this package will depend on visit [../balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- To find creator auto-stake market creation flow visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).

## Component And Connection Map

- The position contract boundary can be found in [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/service.go) and [controller.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/controller.go).
- The position creation ID path can be found in [ids.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/ids.go).
- The position validation boundary can be found in [validator.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/validator.go).
- The position service orchestration can be found in [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/position/service_impl.go).
- The market activity stream handoff can be found in [../realtime/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/realtime/README.md).
- The shared position vocabulary can be found in [../domain/position.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/position.go).
- The market connection can be found in [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- The balance connection can be found in [../balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- The settlement status updates that later consume these records can be found in [../settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- The persistence connection can be found in [../storage/position_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/position_repository.go).
