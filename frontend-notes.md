# Frontend Notes

## Creator Stake Flow

- Market creation is no longer metadata-only.
- The market creator must choose a side (`yes` or `no`) and a stake amount as part of creating a prediction.
- The frontend should treat creator market creation as a guided action that also places the creator's first position automatically.

## UX Notes

- This should be handled carefully in the UI because the creator has more required inputs than a normal participant.
- A simple market-creation form may feel too short once creator side and stake are required.
- A longer structured creation form is acceptable if it stays easy to scan.
- A secondary confirmation card or modal can also work if it keeps the main form from feeling crowded.

## API Notes

- The frontend should prefer one backend route for "create market with creator auto-stake" instead of orchestrating two separate calls client-side.
- I want the backend to own the transaction that creates the market, creates the creator's first position, and locks the creator stake.
- That keeps the frontend flow easier to follow and avoids half-finished states where a market exists but the creator stake failed.

## Suggested Product Shape

- Creator flow:
  - define market
  - choose side
  - choose stake
  - submit once
- Participant flow:
  - view market
  - choose side
  - choose stake
  - place position
