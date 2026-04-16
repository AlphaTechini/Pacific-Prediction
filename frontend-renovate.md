# Frontend Renovate Report

## Purpose

I use this report to keep the frontend cleanup list honest against the current implementation.

The product pages are now mostly connected to backend truth, so this document no longer treats the whole frontend as a static concept build. It tracks the remaining polish and product-alignment work after the live data pass.

## Current Product Reality

- The dashboard reads real markets, balance, and positions.
- The create-market page submits one backend-owned creator flow with side and stake.
- Market detail and resolved-market pages read backend market state by id.
- The portfolio reads real player balance and positions.
- The leaderboard is backed by a backend snapshot route.
- The top navigation uses the static PNG brand icon from `Frontend/static/favicon.png`.
- Guest continuity now uses a localStorage cache for player id and display name, while protected backend calls still depend on the secure session cookie.

## Remaining Cleanup

- The landing page still carries the most aspirational copy and should be tightened to the actual v1 product: virtual prediction markets, Pacifica read-only data, and leaderboard play.
- Any remaining references to AI insight, real trading, wallet claims, or unsupported signal types should be removed unless the backend supports them.
- The frontend still uses `@sveltejs/adapter-auto`; deployment alignment should switch it to `@sveltejs/adapter-vercel` when deployment work resumes.
- Some pages still fetch after mount instead of using route-level loads. That is acceptable for now, but the loading strategy should be unified when speed or SEO starts to matter.
- Realtime support exists in the backend, but the frontend should only consume it where it improves visible market state instead of adding decorative motion.

## Page Notes

- The landing page should become a truthful product introduction instead of a concept deck.
- The dashboard should stay focused on active/resolved markets and current player state.
- The create-market page should keep one backend request for market creation plus creator opening stake.
- The market-detail page should keep settlement and balance math out of the browser.
- The portfolio page should remain a read-only account surface unless deeper history filters are added.
- The leaderboard page should stay a public read-mostly surface until player profile drill-downs exist.

## Component Notes

- [Frontend/src/lib/components/TopNavBar.svelte](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/components/TopNavBar.svelte) owns the current nav chrome and brand icon usage.
- [Frontend/src/lib/guest-session.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/lib/guest-session.ts) owns guest bootstrapping plus the localStorage guest cache.
- [Frontend/static/favicon.png](file:///C:/Hackathons/Pacific%20Prediction/Frontend/static/favicon.png) is the current shared favicon and navigation brand mark.
- [Frontend/src/routes/api/[...path]/+server.ts](file:///C:/Hackathons/Pacific%20Prediction/Frontend/src/routes/api/[...path]/+server.ts) keeps browser requests on the app origin and proxies them to the backend.

## Product-Alignment Recommendation

I should keep the strong visual direction, but every visible claim should map to a real backend capability. If a section has no backend support and no immediate product value, I should defer it instead of keeping it as decorative noise.
