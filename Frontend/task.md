# SvelteKit Screen Conversion & Migration Plan

This document outlines the systematic, detailed steps to transform the raw HTML exports from the Stitch design tool into a fully functional, production-ready SvelteKit web application. 

## 1. Asset Preparation & Organization
*Objective: Organize static assets so they can be bundled efficiently by SvelteKit (Vite).*
- [x] Create directory `src/lib/assets/mockups/` if it doesn't exist.
- [x] Move all downloaded `.png` files from `stitch_screens/` into `src/lib/assets/mockups/` for visual reference during development.
- [x] Extract any actual image assets (icons, logos) from the HTML files and save them to `src/lib/assets/icons/` to be used via direct imports in Svelte.

## 2. Design System & Global Styles Integration
*Objective: Implement the "Obsidian Pulse" design intelligence systematically across Tailwind CSS and global stylesheets.*
- [x] **Typography Setup:**
  - Inject Google Fonts import for `Space Grotesk` (Headings/Data Anchors) and `Inter` (Body/Labels) into `src/app.css`.
- [x] **Color Palette & Theme Variables (Tailwind CSS v4):**
  - Read `stitch_screens/design_system.json`.
  - Transfer `namedColors` into CSS Variables inside `src/app.css`.
  - Configure Tailwind directives referencing these variables to enforce strict adherence to the **"Tonal Layering"** and **"No-Line"** rules.
- [x] **Global Utilities Construction:**
  - Create global classes for the **"Glassmorphism"** mechanics (60% opacity with 20px backdrop blur on `#323538` top surfaces).
  - Add classes for the "Liquid Crystal" signature gradients (`#dbfcff` to `#00f0ff` at 135deg).

## 3. Base Layout & Component Extraction
*Objective: Adhere to strict modularity (1 file = 1 feature/major function) by carving out recurring UI elements before page assembly.*
- [x] **Root Layout (`src/routes/+layout.svelte`):**
  - Implement the fundamental grid/shell for the app (e.g., persistent top/side navigation, background color containment).
- [x] **Atomic Components (`src/lib/components/`):**
  - Extract `Button.svelte` (Primary gradient, ghost/secondary, and tertiary variants).
  - Extract `StatusBadge.svelte` (Success, Danger, Warning signals).
  - Extract `Input.svelte` (Minimal visual borders, terminal cursor focus effects).
  - Extract `MarketCard.svelte` (Reusable block for standard prediction markets to prevent monolithic files).
  - Extract `TopNavBar.svelte` (Persistent navigation with active page highlighting).

## 4. Route-by-Route Conversion
*Objective: Systematically parse the raw `stitch_screens/*.html`, mapping them into correct endpoint directories using Svelte 5 syntax.*

### 4.1. Landing Page (`/`)
- [x] **Path:** `src/routes/+page.svelte`
- [x] **Task:** Translated, using `TopNavBar`, `Button` components and `{#each}` loops.

### 4.2. Dashboard (`/dashboard`)
- [x] **Path:** `src/routes/dashboard/+page.svelte`
- [x] **Task:** Processed with `MarketCard` component and `{#each}` loops.
- [x] **Logic Integration:** Svelte 5 `$state` for filter toggling (ALL/CRYPTO/MACRO).

### 4.3. Create Market (`/markets/create`)
- [x] **Path:** `src/routes/markets/create/+page.svelte`
- [x] **Task:** All form inputs use `$state` runes (asset, expiry, visibility, question binding).
- [x] **Logic Integration:** Live preview panel reflects `$state` in real time.

### 4.4. Market Details (`/markets/[id]`)
- [x] **Path:** `src/routes/markets/[id]/+page.svelte`
- [x] **Task:** Place Bet panel uses `$state` for side selection (YES/NO), modularized from main content.

### 4.5. Market Resolved (`/markets/[id]/resolved`)
- [x] **Path:** `src/routes/markets/[id]/resolved/+page.svelte`
- [x] **Task:** Winner glow mechanics, event lifecycle timeline with `{#each}`.

### 4.6. Portfolio (`/portfolio`)
- [x] **Path:** `src/routes/portfolio/+page.svelte`
- [x] **Task:** Processed with `{#each holdings as h}` loops and tab switching via `$state`.

### 4.7. Leaderboard (`/leaderboard`)
- [x] **Path:** `src/routes/leaderboard/+page.svelte`
- [x] **Task:** Hierarchical podium, data table via `{#each}`, tab switching via `$state`.

## 5. Polish & Verification
*Objective: Guarantee type safety, linting, and correctness.*
- [x] Run Svelte type checks (`svelte-check`) — **0 errors, 4 minor a11y warnings (non-blocking)**.
- [x] Removed all hardcoded `<style>` blocks from raw HTML.
- [x] All pages use Tailwind design tokens from `app.css`.
- [x] Layouts are responsive (mobile/desktop grid breakpoints applied).
