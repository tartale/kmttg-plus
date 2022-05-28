#!/usr/bin/env bash

set -Eeuo pipefail

THIS_SCRIPT_DIR="$(cd $(dirname ${BASH_SOURCE}); pwd)"

usage() {
  echo "
usage: ${0} [flags]
            [-h|--help]


Flags                           Purpose
<none>

Environment variable            Purpose
APP_DIR
MOUNT_DIR
TOOLS_DIR

INPUT_DIR
OUTPUT_DIR
COMSKIP_FILE
ENCODER_DIR
ENCODER_NAME

" >&2

    exit 1
}


function removeOverriddenEntries() {

  inputPath="${1}"
  outputPath="${2}"

  # https://stackoverflow.com/a/56977026/1258206
  #   Iterate through the file, only keep the first <xyz> entry that we see

  awk '
    /^<.*>$/{ 
      hdr = $0
      c[hdr]++
      next
    }
    
    a[hdr] == "" {
      print hdr
      a[hdr] = $0;
      print $0
      next
    }
    
    c[hdr] == 1 {
      a[hdr] = a[hdr] ORS $0
      print $0
      next
    }
  ' < "${inputPath}" > "${outputPath}"

}

function mergeIniFiles() {
  basePath="${1}"
  overridesPath="${2}"
  outputPath="${3}"

  if [[ -e "${overridesPath}" ]]; then
    export OVERRIDES=$(envsubst < "${overridesPath}")
  else
    unset OVERRIDES
  fi

  tmpOutputPath="${TMPDIR}/merged.ini"
  envsubst < "${basePath}" > "${tmpOutputPath}"
  removeOverriddenEntries "${tmpOutputPath}" "${outputPath}"  
}

if [[ -z "${APP_DIR+x}" ]] || [[ -z "${MOUNT_DIR+x}" ]] || [[ -z "${TOOLS_DIR+x}" ]]; then
  usage
fi

export TMPDIR="${TMPDIR:-${PWD}/tmp}"
export OUTPUT_DIR="${OUTPUT_DIR:-${MOUNT_DIR}/output}"
export INPUT_DIR="${INPUT_DIR:-${MOUNT_DIR}/input}"
export COMSKIP_FILE="${COMSKIP_FILE:-comskip.ini.us-ota}"
export ENCODER_DIR="${ENCODER_DIR:-${INPUT_DIR}/encoders}"
export ENCODER_NAME="${ENCODER_NAME:-none}"

umask 000

echo "creating required files/directories in ${MOUNT_DIR}"
mkdir -p "${TMPDIR}"
mkdir -p "${INPUT_DIR}/files"
mkdir -p "${ENCODER_DIR}"
mkdir -p "${OUTPUT_DIR}/download"
mkdir -p "${OUTPUT_DIR}/mpeg"
mkdir -p "${OUTPUT_DIR}/encode"
mkdir -p "${OUTPUT_DIR}/qsfix"
mkdir -p "${OUTPUT_DIR}/webcache"
mkdir -p "${OUTPUT_DIR}/mpegcut"
mkdir -p "${OUTPUT_DIR}/logs"
touch "${OUTPUT_DIR}/logs/auto.history"
touch "${OUTPUT_DIR}/logs/auto.log.0"

echo "merging configuration base and overrides files"
mergeIniFiles "${THIS_SCRIPT_DIR}/config.ini.base" "${INPUT_DIR}/config.ini.overrides" "${INPUT_DIR}/config.ini"
mergeIniFiles "${THIS_SCRIPT_DIR}/auto.ini.base" "${INPUT_DIR}/auto.ini.overrides" "${INPUT_DIR}/auto.ini"

echo "linking input/output files to app home directory"
ln -f -s "${INPUT_DIR}/config.ini" "${APP_DIR}/config.ini"
ln -f -s "${INPUT_DIR}/auto.ini" "${APP_DIR}/auto.ini"
ln -f -s "${APP_DIR}/${COMSKIP_FILE}" "${APP_DIR}/comskip.ini"
ln -f -s "${OUTPUT_DIR}/logs/auto.history" "${APP_DIR}/auto.history"
ln -f -s "${OUTPUT_DIR}/logs/auto.log.0" "${APP_DIR}/auto.log.0"

if [[ -e "${ENCODER_DIR}/${ENCODER_NAME}.enc" ]]; then
  ln -f -s "${ENCODER_DIR}/${ENCODER_NAME}.enc" "${APP_DIR}/encode/"
fi

echo "running kmttg"
${APP_DIR}/kmttg -a
