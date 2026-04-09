# backend/internal/balance

## Purpose

I use this folder for virtual balance state, balance locking, unlocking, debits, credits, and payout application.

## Architectural Decisions And Tradeoffs

- I keep balance separate from player and position logic because anything that controls spendable value deserves a first-class boundary.
- This package should own available-balance rules, locked-balance rules, and authoritative mutation paths.
- The public balance endpoint is still read-only, but position placement and settlement now rely on repository-backed stake locking and settlement-time balance application.
- The tradeoff is another module boundary, but it makes payout and stake handling much safer to evolve later.

## Logic Tracking

- To find balance logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- To find balance service contracts visit [contracts.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/contracts.go).
- To find balance service behavior visit [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/service.go).
- To find player identity logic visit [../player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- To find settlement payout coordination visit [../settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).

## Component And Connection Map

- The authoritative balance system can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- The public balance read path can be found in [service.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/service.go).
- The settlement-time balance persistence connection can be found in [../storage/balance_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/balance_repository.go) and [../storage/balance_postgres_repository.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/balance_postgres_repository.go).
- The persistence boundary behind balance rules can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
