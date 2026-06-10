---
name: project-add-episode-poster
description: "Active branch: add-episode-poster — episode poster from MKV, TheTVDB poster search for movies/series"
metadata: 
  node_type: memory
  type: project
  originSessionId: de3786a1-2d75-421a-85b3-02b3c45276e0
---

Branch `add-episode-poster` across three repos (all pushed — as of 2026-06-10, latest 2026-06-10):
- `github.com/tartale/kmttg-plus` (parent)
- `github.com/tartale/kmttg-go` — Go backend
- `github.com/tartale/kmttg-webui` — React frontend

**What's done:**

- **Episode poster from MKV** (`/api/poster?showID=<id>`): extracts a frame via ffmpeg. `shows.PosterTimeS(show)` picks 1s into chapter 2 if ≥2 chapters exist, else falls back to 5 minutes. `IconCell` renders the poster for `status.downloaded` episodes; episode check moved before `indent` bail-out (fix for indent=true on all episodes).

- **Custom poster file storage + save fix** (2026-06-10):
  - Pasted images are POST'd to `POST /api/custom-poster?showID=<id>` as raw JPEG bytes; saved to `<CacheDir>/posters/<showID>.jpg`
  - `GET /api/custom-poster?showID=<id>` serves the file; this relative URL is stored as `CustomPosterURL`
  - `tivos.SetCustomPosterURL(showID, url)` helper in `tivos.go` replaces direct `UpdateShow` lambda in callers
  - Avoids passing large data URLs through GraphQL (root cause of original save-not-working bug)
  - HTTP URLs still use `setPosterURL` GraphQL mutation directly
  - **Optimistic UI**: `onPosterSaved(dataURL)` callback → `localPosterURL` state in `IconCell` → new image shows immediately without waiting for `refetchQueries`
  - `refetchQueries` fires without await (background cache sync)
  - `PosterSearchDialog` receives `currentImageURL` prop (= `localPosterURL || show.imageURL`) so reopening the dialog shows the already-saved image in the "Current" card

- **Poster paste/drop zone** (replaced TVDB search, 2026-06-10):
  - Dialog has a prominent dashed paste/drop zone; supports Ctrl+V (image or http URL) and drag-and-drop
  - Image blobs downscaled to 320px JPEG on canvas; http URLs stored as-is
  - All TVDB code removed: `pkg/tvdb/`, `KMTTG_TVDB_API_KEY`, schema entries, generated.go functions, model struct

- **TheTVDB poster search (removed)** — was:
  - `KMTTG_TVDB_API_KEY` config field (optional, no validation)
  - `pkg/tvdb/client.go`: token-cached client (29-day TTL), `Search(ctx, query, kind)` — filters results with no image
  - GraphQL: `tvdbSearch(query, kind) → [TvdbSearchResult]`, `setPosterURL(showID, url) → Boolean`, `tvdbAvailable → Boolean`
  - `DetailedShow.CustomPosterURL`: persisted in JSON cache, survives TiVo reloads via `preserveOverrides`; `WithImageURL` prefers it over `findBestImageURL`
  - `PosterSearchDialog`: click placeholder → dialog opens pre-searched with show title; re-searchable; poster grid (2:3 aspect, hover highlight); click to save via `setPosterURL` + `refetchQueries`
  - Search button disabled with `opacity: 0.45` + blue color when `tvdbAvailable` is false; shows config hint
  - `onClick={e => e.stopPropagation()}` on Dialog prevents row toggle (React portal synthetic event bubbling)

- **Poster dialog enhancements v2** (2026-06-10):
  - Ctrl+V paste: `document.addEventListener('paste')` while dialog open; dashed placeholder card replaces Paste button
  - `POSTER_CARD_H = 90px` shared constant for Current and paste placeholder cards; `objectFit: contain` on current image
  - Grid uses `alignItems: start` so cards don't stretch to match tallest neighbor

- **Poster dialog enhancements v1** (2026-06-10):
  - Any movie/series poster (image or placeholder) is now clickable — unified clickable `Box` wrapper with hover opacity dim
  - Dialog shows current poster as first card in results grid (natural aspect ratio, non-interactive, dimmed "Current" label)
  - "Paste from Clipboard" button: reads clipboard image, downscales to 320px JPEG on canvas, shows as preview card in grid
  - `tvdbSearch` query now skips when `!tvdbAvailable` (key not configured) — prevents red GraphQL error on open
  - `tvdbAvailable` defaults to `false` while loading so search never fires prematurely

**Known issue (potential follow-up):** same React portal click propagation fix should be applied to the chapter preview dialog in `ActionCell` — mentioned but not yet implemented.

**Why:** tartale is building a TiVo DVR management tool; episode posters from MKV content + ability to assign show art from TheTVDB.
