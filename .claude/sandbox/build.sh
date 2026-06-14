#!/usr/bin/env bash

THIS_SCRIPT_DIR="$(cd $(dirname ${BASH_SOURCE}); pwd)"

PLUGINS="${THIS_SCRIPT_DIR}/plugin.sh" CS_IMAGE_TAG=kmttg \
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/tartale/claude-sandbox/refs/heads/main/build-image.sh)"
