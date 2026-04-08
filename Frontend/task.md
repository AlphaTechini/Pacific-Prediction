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
- [ ] **Root Layout (`src/routes/+layout.svelte`):**
  - Implement the fundamental grid/shell for the app (e.g., persistent top/side navigation, background color containment).
- [ ] **Atomic Components (`src/lib/components/`):**
  - Extract `Button.svelte` (Primary gradient, ghost/secondary, and tertiary variants).
  - Extract `StatusBadge.svelte` (Success, Danger, Warning signals).
  - Extract `Input.svelte` (Minimal visual borders, terminal cursor focus effects).
  - Extract `MarketCard.svelte` (Reusable block for standard prediction markets to prevent monolithic files).

## 4. Route-by-Route Conversion
*Objective: Systematically parse the raw `stitch_screens/*.html`, mapping them into correct endpoint directories using Svelte 5 syntax.*

### 4.1. Landing Page (`/`)
- [ ] **Path:** `src/routes/+page.svelte`
- [ ] **Task:** Translate the entry view (`Pacifica_Pulse_Landing_Page.html`). Remove static `style` blocks from HTML and replace them with Tailwind utility classes aligning to our configured design system.

### 4.2. Dashboard (`/dashboard`)
- [ ] **Path:** `src/routes/dashboard/+page.svelte`
- [ ] **Task:** Process `Pacifica_Pulse_Dashboard.html`.
- [ ] **Logic Integration:** Use Svelte 5 `$state` runes for local interaction states (e.g., toggling between "Trending" and "New" feeds).

### 4.3. Create Market (`/markets/create`)
- [ ] **Path:** `src/routes/markets/create/+page.svelte`
- [ ] **Task:** Process `Create_Market_Pacifica_Pulse.html`.
- [ ] **Logic Integration:** Implement controlled form inputs reflecting the recessed depths (`surface_container_lowest`) strategy detailed in the design system.

### 4.4. Market Details (`/markets/[id]`)
- [ ] **Path:** `src/routes/markets/[id]/+page.svelte`
- [ ] **Task:** Process `Market_Details_Pacifica_Pulse.html`.
- [ ] **Complexity Constraint:** Modularize the heavy data elements. The "Place Bet" action window must be separated, or structured via Svelte snippet blocks to reduce monolithic size. Render fluctuating odds using tabular numerical fonts.

### 4.5. Market Resolved (`/markets/[id]/resolved`)
- [ ] **Path:** `src/routes/markets/[id]/resolved/+page.svelte`
- [ ] **Task:** Process `Market_Resolved_Pacifica_Pulse.html`. Translate final states, visual "winner" emphasis using the active glow mechanics.

### 4.6. Portfolio (`/portfolio`)
- [ ] **Path:** `src/routes/portfolio/+page.svelte`
- [ ] **Task:** Process `Portfolio_Pacifica_Pulse.html`. 
- [ ] **Logic Integration:** Separate data tables vs individual holding cards into loops `{#each holdings as holding}`.

### 4.7. Leaderboard (`/leaderboard`)
- [ ] **Path:** `src/routes/leaderboard/+page.svelte`
- [ ] **Task:** Process `Leaderboard_Pacifica_Pulse.html`. Focus on the hierarchical layering of ranks.

## 5. Polish & Verification
*Objective: Guarantee type safety, linting, and correctness.*
- [ ] Run Svelte type checks (`svelte-check`) over all translated `.svelte` files.
- [ ] Remove all hardcoded `<style>` blocks imported from the raw HTML to prevent CSS clashes.
- [ ] Ensure all referenced image links point properly to `$lib/assets/` resolving via Vite.
- [ ] Validate layouts against mobile/desktop aspect ratios based on design limits.
