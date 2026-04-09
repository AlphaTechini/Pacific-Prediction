# Frontend Renovate Report

## Purpose

I wrote this report to point out frontend pieces that look polished but are not helping the current product yet.

I am not changing wording or layout here. I am only listing what feels unnecessary, fake, unsupported, or not connected to real behavior.

## General Rule

If a section does not do anything, does not show real data, or suggests a feature that v1 does not support, I should remove it or defer it instead of dressing it up.

## Landing Page

- `Launch Terminal` and `View Markets` look important, but they do not currently drive a real user flow.
- The hero data card shows made-up values and confidence numbers.
- The large decorative images are not tied to the product and can distract from the real flow.
- The `AI Pulse Insight` block does not match the current v1 scope.
- The `Signal Lifecycle` section promises things like smart oracle-style flow in a way that is stronger than the actual frontend experience right now.
- Footer links like `API Documentation`, `Risk Disclosure`, `Privacy Policy`, and `System Status` currently point nowhere useful.

## Top Navigation

- `Log In` does not match the current guest-session product shape.
- `Launch Terminal` appears again as a major action, but it still needs a real destination and purpose.
- `Leaderboard` is visible in the main nav even though it is not backed by real data yet.

## Dashboard

- The featured market cards are hardcoded sample cards.
- The `ALL`, `CRYPTO`, and `MACRO` filter buttons only change local button state right now.
- The left sidebar items like `Live Feed`, `High Volatility`, `New Markets`, `Ending Soon`, and `Watchlist` all look interactive, but they do not do real work.
- The `Trade Now`, `Support`, and `API Documentation` actions in the sidebar are placeholders.
- The summary strip uses made-up counts like active markets, resolving soon, participation, and win rate.
- `Advanced Liquidity Analysis` is a decorative callout, not a real feature.
- The `Live Signal Panel` is sample data.
- `Your Active Markets` is sample data.
- `Recent Results` is sample data.
- `Leaderboard Preview` is sample data.
- `Trending Categories` is sample data and the links do not lead to real filtered views.
- The floating action button in the lower right looks important but does not do anything.

## Create Market Page

- `Open-interest-based` appears as a market type even though that is not in v1.
- `PUBLIC` and `PRIVATE` visibility adds a choice that the backend does not support right now.
- The settlement-source card is fixed text instead of being tied to the selected market setup.
- The preview card shows fake time, fake odds, and fake trader count.
- The `AI Intelligence` panel is not part of the current product scope.
- The form is visually rich, but several fields still behave like a static mock instead of a real creation flow.

## Market Detail Page

- The page does not really use the route id yet. It shows one hardcoded market.
- The stat blocks for price, funding, open interest, and volume are sample values.
- The `Pulse AI Intelligence` block is out of scope for v1.
- The market activity feed is fake.
- The related-market section is hardcoded.
- `Slippage Tolerance` feels unnecessary unless the backend actually supports that concept for this product.
- `Place Trade` looks live, but it is not connected yet.

## Resolved Market Page

- The full resolved view is still hardcoded.
- `Claim Funds to Wallet` does not match the current v1 product, since wallet flow is not part of this version.
- `Join Similar Markets` and `View Portfolio` need to be real routes or should be treated as placeholders.
- The event lifecycle is sample content.
- The market-intelligence block is sample content.

## Portfolio Page

- The holdings table is sample data.
- The stat cards are sample data.
- The tabs look useful, but they do not lead to real filtered content yet.
- `Filter By Asset` is decorative right now.
- `Predictive Insights` is outside the current v1 scope.
- `Performance Matrix` is sample content.
- `Export Performance Audit` is a placeholder action.

## Leaderboard Page

- The whole page is currently a visual concept, not a real product page.
- The tabs are local-only and do not load anything real.
- Podium, ranking table, feed, and vitals are all sample content.
- `Load Detailed Registry` is a placeholder action.
- Since leaderboard support is not part of the current backend work, this page should probably be deferred instead of lightly polished.

## Shared Component Notes

- [Frontend\src\lib\components\TopNavBar.svelte](c:\Hackathons\Pacific Prediction\Frontend\src\lib\components\TopNavBar.svelte) still pushes a login-style flow instead of a guest-session flow.
- [Frontend\src\lib\components\MarketCard.svelte](c:\Hackathons\Pacific Prediction\Frontend\src\lib\components\MarketCard.svelte) is built around sample labels like odds, timer, open interest, and AI insight, so it will need trimming once real backend data is used.
- Some pages contain odd text artifacts from copied content. Those should be cleaned when the real data pass begins.

## Product-Alignment Recommendation

If I follow the product-alignment path, I should do this:

- Keep the strong visual layout.
- Remove or hide anything that implies unsupported features.
- Replace fake content with real backend-backed content.
- If a section has no backend support and no immediate product value, defer it instead of keeping it as decorative noise.
