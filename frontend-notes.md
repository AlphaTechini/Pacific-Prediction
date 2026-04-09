# Frontend Notes

## Current Product Reminders

- I should keep the frontend honest.
- If a feature is not supported in the current backend, I should not make it feel half-live.
- The frontend should talk to our backend only, not directly to Pacifica.
- Guest session is still the main player flow for this version.

## What Is Real Today

- The dashboard is backed by real market, balance, and position reads.
- The create-market page uses one backend flow that includes creator side and creator stake.
- Market detail and resolved pages are backed by real market reads.
- The portfolio uses real player balance and positions.
- The leaderboard is now a real backend-backed page, not a placeholder.

## What Is Still Drift, Not Product Truth

- The landing page still speaks in a more aspirational tone than the actual v1 feature set.
- Some old marketing copy still references unsupported concepts like AI or extra signal types.
- The frontend config still uses `adapter-auto`, but the project standard says `@sveltejs/adapter-vercel`.

## UX Direction

- I should keep the current visual quality, but make the copy and controls more truthful over time.
- If a section is not backed by real data, I should either remove it or clearly simplify it.
- I should avoid stuffing the UI with technical wording just because the backend is ready.
- I should prefer plain labels and simple user messages over frontend jargon.

## Integration Reminder

- The app already has a backend proxy route, so browser requests should keep using it.
- The create-market page should keep using one backend create route for the full creator action.
- The dashboard should stay focused on real markets and real player state.
- The leaderboard should stay a read-only performance surface until deeper profile features exist.

## Practical Notes

- The frontend now builds and type-checks with the live data pages in place.
- The leaderboard uses a route-level load because one snapshot response is cheaper than multiple client requests.
- Other app pages still use lightweight post-mount fetches, which is acceptable for the current scale but not the only possible future direction.
