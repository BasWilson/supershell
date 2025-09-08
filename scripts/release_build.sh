#!/usr/bin/env bash
set -euo pipefail

# Cross-compile supershell for popular OS/ARCH targets and package artifacts.
# Usage:
#   scripts/release_build.sh v0.1.0
#   scripts/release_build.sh               # defaults to "dev"

VERSION="${1:-dev}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="${ROOT_DIR}/dist"
BUILD_DIR="${DIST_DIR}/build"
mkdir -p "${BUILD_DIR}"

targets=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
  "windows arm64"
)

echo "Building supershell ${VERSION}..."
rm -f "${DIST_DIR}"/* || true

for t in "${targets[@]}"; do
  read -r GOOS GOARCH <<<"${t}"
  echo "- ${GOOS}/${GOARCH}"
  outdir="${BUILD_DIR}/${GOOS}_${GOARCH}"
  mkdir -p "${outdir}"

  binname="supershell"
  ext=""
  if [[ "${GOOS}" == "windows" ]]; then
    ext=".exe"
  fi

  CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" \
    go build -trimpath -ldflags "-s -w" -o "${outdir}/${binname}${ext}" ./cmd/supershell

  artifact="supershell_${VERSION}_${GOOS}_${GOARCH}"
  if [[ "${GOOS}" == "windows" ]]; then
    (cd "${outdir}" && zip -q "${DIST_DIR}/${artifact}.zip" "${binname}${ext}")
    shasum -a 256 "${DIST_DIR}/${artifact}.zip" > "${DIST_DIR}/${artifact}.zip.sha256"
  else
    (cd "${outdir}" && tar -czf "${DIST_DIR}/${artifact}.tar.gz" "${binname}${ext}")
    shasum -a 256 "${DIST_DIR}/${artifact}.tar.gz" > "${DIST_DIR}/${artifact}.tar.gz.sha256"
  fi
done

echo "Artifacts in ${DIST_DIR}:"
ls -l "${DIST_DIR}" | sed -e 's/^/  /'


