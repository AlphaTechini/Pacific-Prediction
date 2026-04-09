# Frontend/src/routes/markets

## Purpose

I use this folder for market-related route groups: creation, detail, and resolved views.

## Architectural Decisions And Tradeoffs

- I keep market creation separated from market detail because they have different data and interaction needs.
- I keep the resolved view nested under the market-detail route so the route structure mirrors the product flow.

## Logic Tracking

- To find the create-market route visit [create/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/create/README.md).
- To find the market-detail route visit [market detail README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/[id]/README.md).
- To find the resolved-market route visit [resolved market README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/[id]/resolved/README.md).

## Component And Connection Map

- The market creation connection can be found in [create/README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/create/README.md).
- The market detail connection can be found in [market detail README.md](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/markets/[id]/README.md).
