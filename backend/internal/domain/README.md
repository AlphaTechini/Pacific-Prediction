# backend/internal/domain

## Purpose

I use this folder for shared backend domain vocabulary that multiple modules need to agree on before they start implementing their own logic.

## Architectural Decisions And Tradeoffs

- I keep enum-like values, shared ID types, UTC timestamp handling, and reusable validation errors here so modules do not drift into conflicting copies.
- I prefer a dedicated domain package over putting shared constants in `storage` because domain rules should not depend on infrastructure ownership.
- The tradeoff is one more package, but it gives me a clean place to evolve common business language without creating circular dependencies.

## Logic Tracking

- To find shared ID aliases visit [ids.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/ids.go).
- To find shared market vocabulary visit [market.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/market.go).
- To find shared candle interval helpers visit [candle_interval.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/candle_interval.go).
- To find shared position vocabulary visit [position.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/position.go).
- To find shared UTC timestamp helpers visit [time.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/time.go).
- To find reusable validation errors visit [errors.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/errors.go).
- To find the storage contracts that depend on this vocabulary visit [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
- To find the architecture decisions behind these values visit [../../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).

## Component And Connection Map

- The shared backend domain vocabulary can be found in [market.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/market.go), [candle_interval.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/candle_interval.go), [position.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/position.go), and [ids.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/domain/ids.go).
- The PostgreSQL persistence boundary that reuses these types can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
