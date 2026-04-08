# backend/internal/balance

## Purpose

I use this folder for virtual balance state, balance locking, unlocking, debits, credits, and payout application.

## Architectural Decisions And Tradeoffs

- I keep balance separate from player and position logic because anything that controls spendable value deserves a first-class boundary.
- This package should own available-balance rules, locked-balance rules, and authoritative mutation paths.
- The tradeoff is another module boundary, but it makes payout and stake handling much safer to evolve later.

## Logic Tracking

- To find balance logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- To find player identity logic visit [../player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- To find settlement payout coordination visit [../settlement/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).

## Component And Connection Map

- The authoritative balance system can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- The persistence boundary behind balance rules can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
