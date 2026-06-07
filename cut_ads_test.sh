#!/bin/bash
# Test ad removal for Abbott Elementary - Candygrams
# Preserves the original file; writes output to ./output/comskip/
#
# OFFSET NOTE: TiVo's clipMetadata timestamps are content-time positions
# relative to when the *network* started broadcasting, not when the recording
# started. For this episode, the network started ~24s before TiVo's scheduled
# start time, so content-time leads video-time by 24116ms.
#
# video_pos = (content_ms - RECORDING_OFFSET_MS) / 1000
#
# Without sync-mark audio fingerprint matching, this offset must be measured
# manually. Segments corrected against observed ad positions (9:37, 12:52,
# 18:28, 22:00, 27:30, 31:09).

set -e

INPUT="${PWD}/output/download/Abbott Elementary/Abbott Elementary - Candygrams - E05S13.mkv"
TMPDIR="${PWD}/output/comskip/Abbott Elementary"
OUTPUT="$TMPDIR/Abbott Elementary - Candygrams - E05S13.mkv"
FFMPEG="${FFMPEG:-ffmpeg}"

mkdir -p "$TMPDIR"

if [ -f "$OUTPUT" ]; then
  echo "Output already exists: $OUTPUT"
  exit 0
fi

echo "Input:  $INPUT"
echo "Output: $OUTPUT"
echo ""

# Content segments with recording offset of 24116ms applied.
# video_pos = (content_ms - 24116) / 1000, clamped to [0, video_duration]
#
# Seg 1: content  0       -> 601580  =>  video  0.000  -> 577.464
# Seg 2: content  797593  -> 1133173 =>  video  773.477 -> 1109.057
# Seg 3: content  1345143 -> 1674187 =>  video  1321.027 -> 1650.071
# Seg 4: content  1889020 -> EOF     =>  video  1864.904 -> EOF

echo "=== Extracting segment 1 (0:00 -> 9:37) ==="
"$FFMPEG" -y \
  -ss 0.000 -to 577.464 \
  -i "$INPUT" -c copy -avoid_negative_ts make_zero \
  "$TMPDIR/segment_00.mkv"

echo "=== Extracting segment 2 (12:53 -> 18:29) ==="
"$FFMPEG" -y \
  -ss 773.477 -to 1109.057 \
  -i "$INPUT" -c copy -avoid_negative_ts make_zero \
  "$TMPDIR/segment_01.mkv"

echo "=== Extracting segment 3 (22:01 -> 27:30) ==="
"$FFMPEG" -y \
  -ss 1321.027 -to 1650.071 \
  -i "$INPUT" -c copy -avoid_negative_ts make_zero \
  "$TMPDIR/segment_02.mkv"

echo "=== Extracting segment 4 (31:05 -> end) ==="
"$FFMPEG" -y \
  -ss 1864.904 \
  -i "$INPUT" -c copy -avoid_negative_ts make_zero \
  "$TMPDIR/segment_03.mkv"

echo "=== Concatenating ==="
cat > "$TMPDIR/concat.txt" <<EOF
file '$TMPDIR/segment_00.mkv'
file '$TMPDIR/segment_01.mkv'
file '$TMPDIR/segment_02.mkv'
file '$TMPDIR/segment_03.mkv'
EOF

"$FFMPEG" -y -f concat -safe 0 -i "$TMPDIR/concat.txt" -c copy "$OUTPUT"

echo "=== Cleaning up ==="
rm "$TMPDIR"/segment_0?.mkv "$TMPDIR/concat.txt"

echo ""
echo "Done: $OUTPUT"
