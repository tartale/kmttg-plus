---
name: reference-tivo-clip-metadata
description: "How TiVo clipMetadata works — segment offsets are content-time (broadcast), not video-file positions; silence detection used for calibration"
metadata: 
  node_type: memory
  type: reference
  originSessionId: 99289808-2fe3-4e7f-80ba-a8f8cda2fada
---

**Source:** `message.ClipMetadata` in `go/pkg/message/metadata.go`

```go
type ClipMetadata struct {
    ClipMetadataID string
    SegmentType    SegmentType            // "adSkip" = content segments
    Segment        []ClipMetadataSegment  // content (non-ad) segments
    SyncMark       []ClipMetadataSyncMark // audio fingerprints (NOT decoded; see note)
}
type ClipMetadataSegment struct {
    StartOffset string  // milliseconds, content-time (broadcast clock)
    EndOffset   string  // milliseconds, content-time (broadcast clock)
}
type ClipMetadataSyncMark struct {
    Hash      string  // 32-bit unsigned int string; TiVo proprietary algorithm
    Timestamp string  // content-time ms
}
```

**Key insight:** `StartOffset`/`EndOffset` are **content-time positions** (broadcast clock), NOT positions in the downloaded video file. The recording may start before or after the content clock's t=0.

**Recording offset:** `video_pos = (content_ms - recordingOffsetMs) / 1000`
- `recordingOffsetMs` = content-time that corresponds to video t=0
- For pre-roll recordings: `recordingOffsetMs ≈ 0` and `firstSeg.startOffset = preroll_ms` (handles offset implicitly)
- For Candygrams edge case: `recordingOffsetMs = 24116ms` (broadcast started 24s before TiVo recording), `firstSeg.startOffset = 0`

**Calibration (implemented):** `calibrateRecordingOffset` in `go/pkg/jobs/comskip.go`
- Finds audio silence at the first content/ad boundary via `ffmpeg silencedetect`
- `offset = seg0.endOffset − silence_end_video_ms`
- Falls back to `firstSeg.startOffset` if no silence detected

**Sync marks (NOT usable without TiVo SDK):**
- `SyncMark.Hash` values are TiVo's proprietary audio fingerprint (32-bit unsigned int)
- Attempted reverse-engineering: CRC32, adler32, FNV-1a, sum-of-abs, spectral energy bands — none matched
- kmttg's Java `AutoSkip.java` does NOT decode these offline; it drives a physical TiVo device via RPC to get positions
- Silence detection is the practical offline alternative

**Candygrams (Abbott Elementary S05E13) corrected segments (recordingOffsetMs=24116):**
```
Seg 1: startOffset=0,       endOffset=601580   → video 0.000s  to 577.464s  (0:00–9:37)
Seg 2: startOffset=797593,  endOffset=1133173  → video 773.477s to 1109.057s (12:53–18:29)
Seg 3: startOffset=1345143, endOffset=1674187  → video 1321.027s to 1650.071s (22:01–27:30)
Seg 4: startOffset=1889020, endOffset=1980000  → video 1864.904s to EOF       (31:05–end)
```
Output confirmed correct at 21:46 duration.
