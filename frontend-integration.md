# Frontend Integration Report

## Goal

I want the frontend to stop acting like a static concept and start behaving like the real product.

I also want to keep this simple:

- no fake live data
- no unsupported v1 features
- no heavy frontend jargon
- no direct calls from the frontend to Pacifica

## Main Rule

The frontend should speak to our backend only.

That keeps the product simpler, keeps rules in one place, and stops the UI from drifting away from backend truth.

## First Change I Would Make

When a new visitor lands in the app, the frontend should quietly create or reuse a guest session first.

That gives the app a real player identity before the user tries to:

- create a market
- place a position
- view balance
- view portfolio

Backend route:

- `POST /api/v1/players/guest`

## Shared Integration Work

I would add one small frontend data layer that handles:

- backend base URL from environment
- sending cookies with requests
- one clear loading state style
- one clear error state style
- simple plain-English error messages

This should stay lightweight. I do not need a complicated frontend data system to get v1 working.

## Dashboard

What the page should become:

- a real list of active markets
- a real list of resolved markets if needed
- real player balance and participation summary where useful

Backend routes:

- `GET /api/v1/markets`
- `GET /api/v1/players/me`
- `GET /api/v1/players/me/balance`
- `GET /api/v1/players/me/positions`
- `GET /api/v1/stream`

What I would change:

- Replace the hardcoded market cards with real market data.
- Make card clicks open the correct market detail page.
- Remove filters that do not map to real backend-backed categories yet.
- Remove fake side panels until they are backed by real player data.
- Use the stream route to refresh market state when new markets are created or markets settle.

## Create Market Page

What the page should become:

- a real market-creation flow
- one submit action
- one backend request

Backend routes:

- `GET /api/v1/markets/context`
- `POST /api/v1/markets`

What I would change:

- Load available symbols and allowed setup choices from the backend context route.
- Keep only market types that the current product supports.
- Remove unsupported choices like open interest and private visibility.
- Submit the full creator flow in one request, including creator side and creator stake.
- After success, send the user to the new market page.

What I should keep especially simple:

- the field labels
- error messages
- submit feedback

## Market Detail Page

What the page should become:

- one real market page for the route id the user opened
- real market status
- real stake action

Backend routes:

- `GET /api/v1/markets/{market_id}`
- `POST /api/v1/markets/{market_id}/positions`
- `GET /api/v1/players/me/balance`

What I would change:

- Load the page from the real market id in the URL.
- Replace fake values with backend values only.
- Let the user choose a side and stake amount, then submit a position.
- Show simple failure messages if the backend rejects the request.
- If the market is already resolved, send the user to the resolved view or render the resolved state clearly.

## Resolved Market View

What the page should become:

- a clean read-only result screen
- a real settlement summary

Backend route:

- `GET /api/v1/markets/{market_id}`

What I would change:

- Use the real result, settlement value, and resolved time from the backend.
- Remove wallet-claim language for now.
- Keep only actions that are real, such as returning to markets or portfolio.

## Portfolio

What the page should become:

- real player profile
- real balance
- real positions

Backend routes:

- `GET /api/v1/players/me`
- `GET /api/v1/players/me/balance`
- `GET /api/v1/players/me/positions`
- `GET /api/v1/markets`

What I would change:

- Replace sample holdings with real positions.
- Show active and settled positions using backend status.
- Show real available balance and locked balance.
- If I keep a `Created Markets` view, I should only show it if I can derive it cleanly from real market data.
- Remove fake analytics until there is a real product reason to add them.

## Landing Page

What the page should become:

- a clean entry point into the real app

What I would change:

- Make the main action lead into the real dashboard flow.
- Remove fake confidence numbers and fake AI sections.
- Keep the story of the product simple and true to what the app can actually do today.

## Leaderboard

Current reality:

- there is no clear backend support for a real leaderboard in the current work

Product-alignment choice:

- I should defer this page
- or remove it from the main nav until it is real

I do not think it helps to spend frontend time polishing a page that is not connected to the current product loop.

## Stream Updates

The backend already has a stream route:

- `GET /api/v1/stream`

I would use it for a small set of simple updates:

- a new market appears
- a market changes
- a market settles

I would not turn this into a complicated live-data layer on day one. I only need enough to keep the dashboard and market views fresh.

## Unsupported Or Misleading UI I Would Cut During Integration

- AI insight sections
- wallet claim actions
- private-market controls
- open-interest market options
- fake filters
- fake leaderboard emphasis
- decorative buttons with no destination

## Recommended Order

1. Guest session flow.
2. Dashboard market list.
3. Create-market flow.
4. Market detail and place-position flow.
5. Portfolio with real balance and positions.
6. Resolved market page.
7. Stream updates.
8. Anything extra after the real product loop feels solid.

## Final Note

If I stay disciplined, the frontend pass should feel cleaner, not more technical.

The goal is not to make the UI sound smart.

The goal is to make the UI truthful, connected, and easy to use.
