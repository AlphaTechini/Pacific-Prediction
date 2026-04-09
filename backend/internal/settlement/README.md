# backend/internal/settlement

## Purpose

I use this folder for expiry scanning, settlement execution, payout calculation, and resolution auditing.

## Architectural Decisions And Tradeoffs

- I keep settlement separate because it is the most correctness-sensitive part of the backend.
- This package should own deterministic result computation and transactional updates across markets, positions, and balances.
- In v1, price-threshold settlement should prefer batched Pacifica REST reads at expiry time, while candle and funding settlement should resolve from historical endpoints on demand.
- A price market should settle only when the Pacifica response timestamp is at or after expiry, and the worker should retry briefly instead of guessing if the first fetch is too early.
- The worker should plan price fetches in shared near-expiry batches instead of scheduling one timer or cron job per market.
- Settlement completion should update positions, clear locked stake accounting, and credit winners inside the same transaction as the market-resolution audit write.
- The tradeoff is that it depends on several other packages, but I prefer one explicit coordination point over scattered settlement code.

## Logic Tracking

- To find settlement orchestration contracts visit [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/service.go).
- To find settlement scan service behavior visit [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/service_impl.go).
- To find price fetch batch planning visit [price_fetch_plan.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/price_fetch_plan.go).
- To find worker lifecycle contracts visit [worker.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/worker.go).
- To find expiry scanner loop behavior visit [worker_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/worker_impl.go).
- To find price settlement resolver contracts visit [price_resolver.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/price_resolver.go).
- To find price settlement resolver behavior visit [price_resolver_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/price_resolver_impl.go).
- To find settlement ID generation visit [ids.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/ids.go).
- To find settlement error markers visit [errors.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/errors.go).
- To find candle settlement resolver contracts visit [candle_resolver.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/candle_resolver.go).
- To find candle settlement resolver behavior visit [candle_resolver_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/candle_resolver_impl.go).
- To find funding settlement resolver contracts visit [funding_resolver.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/funding_resolver.go).
- To find funding settlement resolver behavior visit [funding_resolver_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/funding_resolver_impl.go).
- To find settlement audit mapping models visit [audit.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/audit.go).
- To find settlement tests visit [price_settlement_test.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/price_settlement_test.go), [price_fetch_plan_test.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/price_fetch_plan_test.go), and [worker_e2e_test.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/worker_e2e_test.go).
- To find market rule inputs visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find Pacifica settlement source access visit [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find balance update ownership visit [../balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).

## Component And Connection Map

- The settlement orchestration can be found in [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/service.go).
- The settlement scan flow can be found in [service_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/service_impl.go) and [worker_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/worker_impl.go).
- The price settlement path can be found in [price_resolver.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/price_resolver.go) and [price_resolver_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/price_resolver_impl.go).
- The candle settlement path can be found in [candle_resolver.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/candle_resolver.go) and [candle_resolver_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/candle_resolver_impl.go).
- The funding settlement path can be found in [funding_resolver.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/funding_resolver.go) and [funding_resolver_impl.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/funding_resolver_impl.go).
- The settlement audit persistence can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
