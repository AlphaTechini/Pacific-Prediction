# backend/internal/config

## Purpose

I use this folder for runtime configuration loading and validation.

## Architectural Decisions And Tradeoffs

- I want all environment reads in one package so secrets and deploy-time values stay centralized.
- This package should validate required config at startup instead of letting bad values fail later inside business logic.
- For local development, I allow the config package to load `.env` files from the backend working directory or the `cmd/api` working directory before validation runs.
- I now keep database pool tuning in config too, so local and deployed environments can reduce unnecessary reconnect churn without code changes.
- I also keep market-creation guardrails in config so fairness rules like the allowed price-threshold band can change without a code edit.
- The tradeoff is stricter startup behavior, but I prefer fast failure over hidden runtime drift.

## Logic Tracking

- To find configuration logic visit [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/README.md).
- To find runtime config loading and local `.env` support visit [config.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/config.go).
- To find database pool defaults and validation visit [config.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/config.go).
- To find market guardrail config loading and validation visit [config.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/config.go).
- To find deployment-related backend decisions visit [../../README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/README.md).

## Component And Connection Map

- The runtime configuration boundary can be found in [README.md](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/README.md).
- The market guardrail configuration can be found in [config.go](file:///C:/Hackathons/Pacific%20Prediction/backend/internal/config/config.go).
- The environment-driven system design can be found in [../../../architecture.md](file:///C:/Hackathons/Pacific%20Prediction/architecture.md).
