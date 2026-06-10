---
name: project-add-skip-segments
description: "Merged to master — chapter preview UI with add/remove, save/reset/re-run buttons, Chapter 0 prepend, skip icon rollup"
metadata: 
  node_type: memory
  type: project
  originSessionId: 99289808-2fe3-4e7f-80ba-a8f8cda2fada
---

Branch `add-skip-segments` — **merged to master** across all three repos (2026-06-10).
- `github.com/tartale/kmttg-plus` (parent)
- `github.com/tartale/kmttg-go` — Go backend
- `github.com/tartale/kmttg-webui` — React frontend

**What's done:**

- **Chapter 0 prepend**: `writeChaptersXML` and `writeChaptersXMLFromOverrides` both prepend a "Chapter 0" atom at `00:00:00` when the first marker is non-zero. Ensures VLC shows chapter navigation (requires ≥2 chapters).
- **Override chapters on non-adSkip shows**: `Chapters()` early-return guard changed from `!show.GetAdSkip()` to `!show.GetAdSkip() && !hasOverrides` — fixes silent data loss when user manually adds chapters to a show with no TiVo adSkip data.
- **Skip icon shows for chaptered files**: `TitleCell` shows `skip.svg` when `show.chapters?.length > 0 || show.adSkip`. Series row also shows icon if any episode has chapters (`episodes.some(ep => ep.chapters?.length > 0)`).
- **Chapter model**: startOffset-only — each "Chapter N" spans from its startOffset to the next (no endOffset, no Ad Break interleaving). `writeChaptersXML` and `writeChaptersXMLFromOverrides` both follow this model. `readChaptersFromFile` filters on `"Chapter"` prefix.
- **Chapter 1 adjustable**: `cutTimes` initialized from all chapters (no `.slice(1)`). `handleSave` maps `cutTimes[i]` directly to each chapter start.
- **Add/remove chapters in UI**: `ChapterStrip` accepts `cutTimes`/`onCutTimesChange` from `ActionCell`. `handleAdd` inserts midpoint between last chapter and duration; `handleRemove` filters by index. `AddChapterDivider` between panels and at bottom. Delete (×) button in `BoundaryPanel` header.
- **Scroll-to-bottom on add**: `bottomRef`/`prevLengthRef` refs in `ChapterStrip` trigger `scrollIntoView({ behavior: 'smooth' })` when `cutTimes.length` increases.
- **"No chapter data" screen**: when `show.chapters.length < 1 && cutTimes.length < 1`, shows fallback with "Re-run Chapter Detection" + "Add First Chapter" button (disabled styling: `opacity: 0.45` with blue color preserved via `'&.Mui-disabled'`).
- **Explicit Save button** (replaced auto-save): `isSaving` state + `LinearProgress` under button; closes dialog on CHAPTERS job COMPLETE; does NOT close on error.
- **Reset to defaults button**: `isRerunning` state + `LinearProgress`; saves empty overrides → triggers RerunChapters → clears state on job COMPLETE.
- **Re-run Chapter Detection button**: same `isRerunning` state + progress; "Add First Chapter" button disabled while rerunning.
- **Chapter preview UI** (`AdSkipStrip` → `BoundaryPanel` + `DockThumbnail` in `ShowRow.tsx`):
  - Per-boundary framer-motion dock-style scrub strip; `GET /api/thumbnail?showID=<id>&t=<seconds>` for frames
  - framer-motion v6 required (v10 uses `useId` which requires React 18; this project is React 17)
  - `SliderWithHover` component: hover shows amber vertical line + timestamp; click jumps cutMs

**Why:** personal TiVo DVR management tool; strip ads from downloaded recordings.
