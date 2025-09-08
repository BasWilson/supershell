#!/usr/bin/env bash
set -euo pipefail

# Tag and push a release. Usage: scripts/release_tag.sh v0.1.0

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 vX.Y.Z" >&2
  exit 1
fi

VERSION="$1"

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Version must be in form vX.Y.Z" >&2
  exit 1
fi

git fetch --tags
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "Tag $VERSION already exists." >&2
  exit 1
fi

git add -A
git commit -m "chore(release): $VERSION" || true
git tag -a "$VERSION" -m "$VERSION"
git push origin HEAD --tags
echo "Pushed tag $VERSION"


