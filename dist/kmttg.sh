#!/usr/bin/env bash

set -Eeuo pipefail

THIS_SCRIPT_DIR="$(cd $(dirname ${BASH_SOURCE[0]}); pwd)"

if [[ -f "${THIS_SCRIPT_DIR}/.bashrc" ]]; then
  source "${THIS_SCRIPT_DIR}/.bashrc"
fi

echo "Running ${THIS_SCRIPT_DIR}/kmttg $@"
exec ${THIS_SCRIPT_DIR}/kmttg $@
