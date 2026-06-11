---
name: project-ui-enhancements
description: "Active branch: ui-enhancements — episode title, skip icon cleanup, blue accent bracket for series/episode rows"
metadata:
  type: project
---

Branch `ui-enhancements` across three repos (started 2026-06-10 off master/main):
- `github.com/tartale/kmttg-plus` (parent)
- `github.com/tartale/kmttg-go` — Go backend
- `github.com/tartale/kmttg-webui` — React frontend

**What's done:**

- **Episode title in title column** (2026-06-10):
  - `getTitleExtension` (showListingHelpers.tsx) returns `episodeTitle` for episodes; blank if missing
  - Added `episodeTitle` to the `episodes { ... }` selection in the `GetShows` query

- **Skip icon never on series rows** (2026-06-10):
  - `hasSkipIndicator` in `TitleCell` now gates on `show.kind !== ShowKind.Series`

- **Blue accent bracket for open series + episode rows** (2026-06-11):
  - Constants: `ACCENT_COLOR='#0077B6'`, `ACCENT_LEFT='10px'`, `ACCENT_W='4px'` — all in `ShowRow.tsx` `rowSx`
  - **Series row (open):** `borderLeft: 4px` on first cell; `::before` top accent across all cells (height 4px, top -1px, zIndex 1); first-cell `::before` extends left with `calc(-1 * ACCENT_W)`; `::after` connector at bottom with `left: calc(-1 * ACCENT_W), width: calc(ACCENT_LEFT + ACCENT_W)` — overlaps episode accent to avoid notch
  - **Episode rows (indent):** `::before` vertical accent at `left: 10px` (clear of thumbnail at paddingLeft 20px), `top: -1px`, `bottom: calc(-1 * ACCENT_W)` — extends into last row's bottom accent area
  - **Last episode row:** `::after` bottom accent on all cells (left: 0, right: 0); first cell overrides `left: ACCENT_LEFT` so accent starts at the vertical accent's indentation position

- **App header with TiVo device tabs** (2026-06-11):
  - `AppHeader.tsx` (new): MUI `AppBar` with TiVo logo + "kmttg+" name; `Tabs` row queries TiVo names, auto-selects first on load; no "All" tab
  - `Home.tsx`: flex column layout (`height: 100vh`), owns `selectedTivo` state; content box uses `minHeight: 0` so `TableContainer` is bounded
  - `ShowListing.tsx`: accepts `selectedTivo` prop, filters `data.tivos` before merging; `TableContainer` is scroll container (`height: 100%, overflow: auto`)
  - `ShowRow.tsx` `ShowHeader`: sticky header cells use `backgroundColor: #0d1f35` to prevent rows bleeding through on scroll; bottom border matches accent color

**Why:** tartale is building a TiVo DVR management tool; these are UX polish changes to the show listing.
