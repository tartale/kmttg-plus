---
name: project-audio-effects
description: "audio-effects branch: active as of 2026-06-12; scope TBD"
metadata:
  type: project
---

Branch `audio-effects` created 2026-06-12 off master/main across all three repos:
- `github.com/tartale/kmttg-plus` (parent)
- `github.com/tartale/kmttg-go` — Go backend
- `github.com/tartale/kmttg-webui` — React frontend

**What's done (2026-06-12):**
- Skip back 8s / skip forward 30s overlay buttons added to VideoPlayer (bottom-right), alongside existing chapter skip buttons (bottom-left)
- Left/right arrow keys trigger skip + audio; `preventDefault` suppresses native video seek
- TiVo-style chirp synthesized via Web Audio API: two sine tones at 620Hz/930Hz, 60ms each, 120ms gap
- Forward: each tone starts 6% sharp and settles (upward snap); backward: starts 6% flat and settles (droopy start) — direction-specific pitch overshoot gives distinct feel
- `AudioContext` created lazily on first button press, reused across clicks, closed on unmount

**Why:** TiVo-style audio feedback makes skip actions feel familiar to TiVo users.

**How to apply:** New work goes on the `audio-effects` branch; submodule convention is parent=master, subs=main. Not yet merged to master.
