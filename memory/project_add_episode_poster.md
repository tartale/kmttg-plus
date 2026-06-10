---
name: project-add-episode-poster
description: "Active branch: add-episode-poster — episode poster from MKV, TheTVDB poster search for movies/series"
metadata: 
  node_type: memory
  type: project
  originSessionId: de3786a1-2d75-421a-85b3-02b3c45276e0
---

Branch `add-episode-poster` across three repos (all pushed — as of 2026-06-10):
- `github.com/tartale/kmttg-plus` (parent)
- `github.com/tartale/kmttg-go` — Go backend
- `github.com/tartale/kmttg-webui` — React frontend

**What's done:**

- **Episode poster from MKV** (`/api/poster?showID=<id>`): extracts a frame via ffmpeg. `shows.PosterTimeS(show)` picks 1s into chapter 2 if ≥2 chapters exist, else falls back to 5 minutes. `IconCell` renders the poster for `status.downloaded` episodes; episode check moved before `indent` bail-out (fix for indent=true on all episodes).

- **TheTVDB poster search for Movies/Series**:
  - `KMTTG_TVDB_API_KEY` config field (optional, no validation)
  - `pkg/tvdb/client.go`: token-cached client (29-day TTL), `Search(ctx, query, kind)` — filters results with no image
  - GraphQL: `tvdbSearch(query, kind) → [TvdbSearchResult]`, `setPosterURL(showID, url) → Boolean`, `tvdbAvailable → Boolean`
  - `DetailedShow.CustomPosterURL`: persisted in JSON cache, survives TiVo reloads via `preserveOverrides`; `WithImageURL` prefers it over `findBestImageURL`
  - `PosterSearchDialog`: click placeholder → dialog opens pre-searched with show title; re-searchable; poster grid (2:3 aspect, hover highlight); click to save via `setPosterURL` + `refetchQueries`
  - Search button disabled with `opacity: 0.45` + blue color when `tvdbAvailable` is false; shows config hint
  - `onClick={e => e.stopPropagation()}` on Dialog prevents row toggle (React portal synthetic event bubbling)

**Known issue (potential follow-up):** same React portal click propagation fix should be applied to the chapter preview dialog in `ActionCell` — mentioned but not yet implemented.

**Why:** tartale is building a TiVo DVR management tool; episode posters from MKV content + ability to assign show art from TheTVDB.
