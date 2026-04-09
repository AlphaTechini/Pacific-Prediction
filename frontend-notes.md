# Frontend Notes

## Main Product Reminders

- I should keep the frontend honest.
- If a feature is not supported in v1, I should not make it feel half-live.
- The frontend should talk only to our backend, not directly to Pacifica.
- Guest session is the current player flow. Wallet login is not the main flow for this version.

## Creator Stake Flow

- Market creation is no longer metadata-only.
- The market creator must choose a side (`yes` or `no`) and a stake amount as part of creating a prediction.
- The frontend should treat creator market creation as one guided action that also places the creator's first position automatically.
- I should keep this flow simple on the page: define market, choose side, choose stake, submit once.

## What I Should Avoid Right Now

- AI insight panels.
- Wallet-claim language.
- Open-interest market options.
- Private-market choices.
- Leaderboard-heavy product focus.
- Decorative controls that look clickable but do not do anything.

## UX Direction

- I should keep the current visual quality, but strip out fake power-user noise.
- If a section is not backed by real data, I should either remove it or clearly simplify it.
- I should avoid stuffing the UI with technical wording just because the backend is ready.
- I should prefer plain labels and simple user messages over frontend jargon.

## Integration Reminder

- The create-market page should use one backend create route for the full creator action.
- The trade page should let a player choose side, enter stake, and submit once.
- The dashboard should become the main live view for real markets instead of sample cards.
- The portfolio should show real balance and real positions before anything more advanced is added.

## Practical Notes

- The frontend builds today, but it still has a few accessibility warnings that should be cleaned during the integration pass.
- The frontend config still uses `adapter-auto`, but the project standard says `@sveltejs/adapter-vercel`.
- Some copied text has encoding issues and should be cleaned during the same pass.
- Browser-style interaction from this environment may be limited, so I may be better at code-side frontend work than live browser clicking in this setup.
