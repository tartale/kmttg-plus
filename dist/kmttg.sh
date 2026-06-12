#!/usr/bin/env bash

THIS_SCRIPT_DIR="$(cd $(dirname ${BASH_SOURCE[0]}); pwd)"

if [[ -f "${THIS_SCRIPT_DIR}/.bashrc" ]]; then
  source "${THIS_SCRIPT_DIR}/.bashrc"
fi

exec ${THIS_SCRIPT_DIR}/kmttg
