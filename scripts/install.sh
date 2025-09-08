#!/usr/bin/env bash
set -euo pipefail

# Install the latest supershell binary for macOS or Linux.
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/BasWilson/supershell/main/scripts/install.sh | bash

REPO_OWNER="BasWilson"
REPO_NAME="supershell"

if [[ "${EUID:-$(id -u)}" -ne 0 ]]; then
  SUDO="sudo"
else
  SUDO=""
fi

uname_s=$(uname -s)
case "$uname_s" in
  Darwin) GOOS=darwin ;;
  Linux) GOOS=linux ;;
  *) echo "Unsupported OS: $uname_s" >&2; exit 1 ;;
esac

uname_m=$(uname -m)
case "$uname_m" in
  x86_64|amd64) GOARCH=amd64 ;;
  arm64|aarch64) GOARCH=arm64 ;;
  *) echo "Unsupported ARCH: $uname_m" >&2; exit 1 ;;
esac

api_url="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"
echo "Finding latest release..."
resp=$(curl -fsSL "$api_url")

if command -v jq >/dev/null 2>&1; then
  TAG=$(printf '%s' "$resp" | jq -r .tag_name)
else
  TAG=$(printf '%s' "$resp" | sed -n 's/.*"tag_name" *: *"\([^"]\+\)".*/\1/p' | head -n1)
fi

if [[ -z "${TAG:-}" ]]; then
  echo "Failed to detect latest release tag" >&2
  exit 1
fi

asset_name="supershell_${TAG}_${GOOS}_${GOARCH}.tar.gz"
download_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${TAG}/${asset_name}"

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

echo "Downloading ${asset_name}..."
curl -fsSL "$download_url" -o "$tmpdir/${asset_name}"

echo "Extracting..."
tar -xzf "$tmpdir/${asset_name}" -C "$tmpdir"

dest="/usr/local/bin"
if [[ ! -w "$dest" ]]; then
  if [[ -n "$SUDO" ]]; then
    echo "Installing to $dest (requires sudo)"
    $SUDO install -m 0755 "$tmpdir/supershell" "$dest/supershell"
  else
    echo "No permission to write to $dest and sudo is unavailable" >&2
    exit 1
  fi
else
  install -m 0755 "$tmpdir/supershell" "$dest/supershell"
fi

echo "Installed $($dest/supershell 2>/dev/null || echo supershell) to $dest"
echo "Note: Ensure $dest is on your PATH."


