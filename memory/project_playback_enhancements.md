---
name: project-playback-enhancements
description: "playback-enhancements branch: COMPLETE, merged to master/main 2026-06-12"
metadata:
  type: project
---

Branch `playback-enhancements` merged to master/main 2026-06-12.

**What shipped:**
- Migrated video player from custom hls.js to `@vidstack/react` (v1.15.6) with `DefaultVideoLayout`
- Chapter scrubber markers via inline VTT `Track` component
- TiVo chirp sounds on seek buttons (slot-based `ChirpSeekButton` using `useMediaPlayer`) and keyboard shortcuts
- Arrow key handler uses capture phase + `stopPropagation` so chirps work in fullscreen
- Upgraded React 17 → 18; `ReactDOM.render` → `createRoot`
- Installed `media-icons` peer dep for `defaultLayoutIcons`
- Dialog window sized `min(95vw, (95vh-52px)*16/9)` to fill screen without overflow
- Removed verbose ffmpeg debug log on success in `pkg/streamer/stream.go`
- Fixed `AppHeader` tivos array stability with `useMemo`
- Removed `KMTTG_TOOLS_DIR` config field; ffmpeg/ffprobe now resolved via PATH

**Why:** Better built-in controls, YouTube-style chapter markers on scrubber, trackpad support, fullscreen keyboard shortcuts.

**How to apply:** Submodule convention: parent=master, subs=main. All three repos merged.
