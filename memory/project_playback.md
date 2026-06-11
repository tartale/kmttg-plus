---
name: project-playback
description: "Active branch: playback — new feature work off master/main"
metadata:
  type: project
---

Branch `playback` across three repos (started 2026-06-11 off master/main):
- `github.com/tartale/kmttg-plus` (parent)
- `github.com/tartale/kmttg-go` — Go backend
- `github.com/tartale/kmttg-webui` — React frontend

**What's done (WIP — streaming not yet working):**

- `pkg/streamer/stream.go`: session manager (`GetOrStart`, `Get`, `Stop`, `WaitForPlaylist`); HLS via ffmpeg — file remux (`-c:v copy -bsf:v h264_mp4toannexb`) or live TiVo decrypt+encode; `cmd.Dir=dir` + relative paths so m3u8 has relative URIs; clears session dir on start
- `cmd/kmttg.go`: `/api/stream/{showID}/playlist.m3u8`, `/api/stream/{showID}/{segment}`, `DELETE /api/stream/{showID}`; direct `os.ReadFile`+write (not `http.ServeFile`) to avoid Content-Type override/redirects
- `VideoPlayer.tsx`: hls.js Dialog, init via `TransitionProps.onEntered` (not useEffect — video ref not ready otherwise); error handler logs `Hls.Events.ERROR` details to console
- `ShowRow.tsx`: play button (`PlayArrowIcon`) on movie/episode rows only

**Known issues / next steps:**
- Browser still shows "no video with supported format and MIME type found" — need to check browser console for HLS error details and Go debug logs for playlist content/ffmpeg output

**Why:** tartale is building a TiVo DVR management tool; this branch is for playback-related features.
