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

**Debugging session 2026-06-11 (root cause found & fixed):**
- Live-path bug: transcode runs ~1800x realtime, so `-hls_list_size 10 -hls_flags delete_segments` deleted segments before the player could fetch them (playlist started at MEDIA-SEQUENCE 24+). Fixed: `-hls_list_size 0 -hls_playlist_type event`, no delete_segments; session dir cleanup stays in Stop().
- Hardening: failed ffmpeg runs now remove the session from the manager so retry works (previously dead sessions lingered until dialog close → DELETE).
- Verified end-to-end in headless Chromium with real TiVo data (`test/data/01-decrypted.ts` via git-lfs, mpeg2video+ac3): both file-remux and live-transcode HLS output plays (MANIFEST_PARSED → FRAG_BUFFERED → VIDEO PLAYING) through handlers identical to the Go ones. Codecs, MIME types, relative URIs all correct.
- Added ginkgo tests: `pkg/streamer/stream_test.go` (file streaming, Stop cleanup, failed-session retry).
- Downloaded MKVs are always H.264/AAC (pkg/encoder), so stream-copy remux is browser-safe.

**Firefox bufferAppendError fix (2026-06-11, later same day):**
- User's Firefox console showed `mediaError bufferAppendingError` → fatal `bufferAppendError` on `sourceBufferName: "audio"` — the downloaded MKV has AC-3 audio, which Firefox cannot decode via MSE (video copied fine, audio SourceBuffer.appendBuffer threw). Note: Chromium in sandbox *does* decode AC-3, so this is Firefox-specific.
- Fix: `streamFromFile` now probes codecs with ffprobe (`probeCodec`) and only stream-copies h264 video / aac audio; anything else (ac3, mpeg2video, unknown) gets transcoded. ffprobe failure → transcode (safe default). ffprobe on MPEG-TS returns the codec twice (program + top-level), so probeCodec takes the first line.
- Verified fixed AC-3→AAC output plays in headless Chromium; AC-3 regression test added to stream_test.go.

**Real-Firefox verification (2026-06-11, portable Firefox 151 in sandbox):**
- Installed portable Firefox tarball at /tmp/firefox (no root: extracted via python lzma; apt needs root, xz missing). MOZ_HEADLESS + --profile $(mktemp -d), page beacons logs to test server via fetch("/log?m=").
- BOTH new pipelines' output, built from the real TiVo fixture with the exact new ffmpeg args, PLAY in real Firefox: file path (encoder MKV → probe → h264/aac copy) and live path (mpeg2/ac3 → libx264/aac event playlist). Only benign non-fatal bufferStalledError/bufferSeekOverHole at start.
- Old AC-3-copy output reproduces user's exact error signature. Conclusion: user's running server must still serve old-code (AC-3) segments.
- User rebuilt natively on the mac but said docker image rebuild wasn't done — CRA proxies to localhost:7676; if the docker container (old image) owns 7676, mac rebuild changes nothing. `.go-version` stays 1.25.10 (user installing it).
- VideoPlayer.tsx error handler now logs sourceBufferName + error.message (tsc clean).

**TRUE ROOT CAUSE found & fixed (2026-06-11): 5.1 AAC**
- User's mac checkout is the SAME filesystem as /workspace — their real cache (`go/output/cache/stream/`) and downloads (`go/output/download/`) are directly readable in the sandbox.
- Their cached segments were h264 + **aac 6-channel (5.1)** — Firefox cannot decode multichannel AAC via MSE (stereo is fine, Chrome handles 5.1). Served their actual segments to sandbox Firefox: reproduced the exact fatal `bufferAppendError sourceBuffer=audio`. AC-3 was never the issue for this file; the downloaded MKV (Abbott Elementary E05S13) is h264 + aac 5.1, so old code's `-c:a copy` kept 5.1.
- Fix in stream.go: audio copied only when aac AND 1–2 channels; otherwise `-c:a aac -ac 2` (stereo downmix). Live path also gets `-ac 2`. probeCodec generalized to `probeStream(ctx, path, selector, entry)`.
- Verified: new pipeline output from user's actual MKV plays in sandbox Firefox. Tests: DescribeTable covers AC-3 stereo and AAC 5.1 → both must yield aac/2ch segments. 5/5 pass; make go-build passes (container now has go 1.25.10).

**STATUS: WORKING (2026-06-11).** User rebuilt on macOS and confirmed playback works in Firefox. All changes committed & pushed on `playback` branches (go 5bed9e7, webui 79cd302, parent a46575b).

**Possible follow-ups (not started):**
- React console warning: "non-boolean attribute `indent`" from IconCell→TableCell in ShowRow.tsx.
- Optional enhancement: capability-based audio (5.1 passthrough for Chrome/Safari instead of always downmixing to stereo).
- `make go-build` requires Go 1.25.10 via goenv; sandbox has 1.26.4, used `go build ./...` instead.

**Why:** tartale is building a TiVo DVR management tool; this branch is for playback-related features.
