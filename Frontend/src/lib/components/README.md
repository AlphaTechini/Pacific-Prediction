# Frontend/src/lib/components

## Purpose

I use this folder for reusable UI components shared across the app routes.

## Architectural Decisions And Tradeoffs

- I keep common buttons, inputs, cards, badges, and navigation here instead of duplicating them per page.
- I keep them presentation-focused and leave data orchestration to route files and data helpers.
- I keep the top navigation brand visual pointed at the static `/favicon.png` asset so the nav and browser favicon use the same mark.

## Logic Tracking

- To find navigation chrome and the static brand icon usage visit [TopNavBar.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/TopNavBar.svelte).
- To find the shared button component visit [Button.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/Button.svelte).
- To find the shared input component visit [Input.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/Input.svelte).
- To find the market summary card visit [MarketCard.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/MarketCard.svelte).
- To find the shared status badge visit [StatusBadge.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/StatusBadge.svelte).

## Component And Connection Map

- The top navigation component and brand icon connection can be found in [TopNavBar.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/TopNavBar.svelte).
- The shared market-card connection can be found in [MarketCard.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/MarketCard.svelte).
