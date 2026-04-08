# backend/internal/settlement

## Purpose

I use this folder for expiry scanning, settlement execution, payout calculation, and resolution auditing.

## Architectural Decisions And Tradeoffs

- I keep settlement separate because it is the most correctness-sensitive part of the backend.
- This package should own deterministic result computation and transactional updates across markets, positions, and balances.
- The tradeoff is that it depends on several other packages, but I prefer one explicit coordination point over scattered settlement code.

## Logic Tracking

- To find settlement logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- To find market rule inputs visit [../market/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/market/README.md).
- To find Pacifica settlement source access visit [../pacifica/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/pacifica/README.md).
- To find balance update ownership visit [../player/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).

## Component And Connection Map

- The settlement engine can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/settlement/README.md).
- The settlement audit persistence can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
