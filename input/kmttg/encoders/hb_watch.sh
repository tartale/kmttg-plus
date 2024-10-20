#!/bin/bash

set -xEeuo pipefail

input="${1}"
output="${2}"
srtFile="${3+x}"

subdir="TV"
if [[ "${input}" == *"("*")"* ]]; then
  subdir="movies"
fi

inputSubfile="${input#${OUTPUT_DIR}/}"
inputSubdir=$(dirname "${inputSubfile}")
inputSubdir="${inputSubdir#mpeg/}"
inputSubdir="${inputSubdir#mpegcut/}"
encodeDir=$(dirname "${output}")
watchDir="${OUTPUT_DIR}/watch/${subdir}/${inputSubdir}"

mkdir -p "${watchDir}"
cp -rf "${input}" "${watchDir}/"
cp -rf "${srtFile}" "${watchDir}/" || true
cp -rf "${encodeDir}"/* "${watchDir}/" || true
echo "success!" > "${output}"
