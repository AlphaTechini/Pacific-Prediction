# backend/internal/player

## Purpose

I use this folder for player identity and profile rules only.

## Architectural Decisions And Tradeoffs

- I want player identity separated from auth and balance because identity, session state, and spendable value are different concerns.
- This package should own the canonical player record and profile-facing reads only.
- The player read path depends on auth middleware for player context, which keeps session parsing out of profile logic.
- The tradeoff is a narrower module, but it keeps future auth and payout changes from distorting player ownership.

## Logic Tracking

- To find player identity logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- To find auth and session logic visit [../auth/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/auth/README.md).
- To find balance logic visit [../balance/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/balance/README.md).
- To find HTTP exposure of player state visit [../httpapi/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/httpapi/README.md).

## Component And Connection Map

- The canonical player record can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/player/README.md).
- The player persistence connection can be found in [../storage/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/storage/README.md).
