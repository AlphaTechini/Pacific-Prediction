# backend/internal/domain

## Purpose

I use this folder for shared backend domain vocabulary that multiple modules need to agree on before they start implementing their own logic.

## Architectural Decisions And Tradeoffs

- I keep enum-like values, shared ID types, UTC timestamp handling, and reusable validation errors here so modules do not drift into conflicting copies.
- I prefer a dedicated domain package over putting shared constants in `storage` because domain rules should not depend on infrastructure ownership.
- The tradeoff is one more package, but it gives me a clean place to evolve common business language without creating circular dependencies.

## Logic Tracking

- To find shared market and position vocabulary visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md).
- To find the storage contracts that depend on this vocabulary visit [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- To find the architecture decisions behind these values visit [../../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).

## Component And Connection Map

- The shared backend domain vocabulary can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/README.md).
- The PostgreSQL persistence boundary that reuses these types can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
