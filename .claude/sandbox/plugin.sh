#!/usr/bin/env bash
set -euo pipefail

export LANGUAGE_VERSIONS="go-1.25.10"
${PLUGINS_DIR}/languages/go.sh
${PLUGINS_DIR}/languages/react.sh
apt-get update && apt-get install -y ffmpeg mkvtoolnix && rm -rf /var/lib/apt/lists/*
