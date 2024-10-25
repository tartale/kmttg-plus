#!/usr/bin/env bash

set -Eeuo pipefail

THIS_SCRIPT_DIR="$(cd $(dirname ${BASH_SOURCE[0]}); pwd)"

usage() {
  echo "
usage: ${0} [-h|--help]

Environment variable            Purpose
KMTTG_CERT_ZIP_URI              Points to a location where the latest certificate
                                and password can be found (on the web, on disk, etc).
                                The file is expected to be a zip file containing two entries; 
                                cdata.p12 and cdata.password.

" >&2

    exit 1
}

curl -L "${KMTTG_CERT_ZIP_URI}" -o "${THIS_SCRIPT_DIR}/cdata.zip"
go test ${THIS_SCRIPT_DIR}/../pkg/certificate/...
