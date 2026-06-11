#!/bin/bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_DIR"

if ! git diff --quiet 2>/dev/null; then
  echo "changes detected, committing..."
  git add -A
  git commit -m "sync: $(date '+%Y-%m-%d %H:%M:%S')"
fi

git pull --rebase
git push
echo "valhalla synced."
