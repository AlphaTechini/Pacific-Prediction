# backend/cmd

## Purpose

I use `cmd` for executable entrypoints only. I do not want business logic living here.

## Architectural Decisions And Tradeoffs

- I keep `cmd` thin so bootstrapping, dependency wiring, and runtime startup stay readable.
- The tradeoff is a little more indirection, but it prevents the entrypoint from becoming the whole application.

## Logic Tracking

- To find backend process startup planning visit [api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- To find the backend package split visit [../README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/README.md).

## Component And Connection Map

- The executable entrypoint can be found in [api/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/cmd/api/README.md).
- The dependency ownership behind the entrypoint can be found in [../internal/README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/README.md).
