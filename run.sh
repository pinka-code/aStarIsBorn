#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"

cd "$ROOT_DIR"

echo "Running GitHub random fetch..."

go run ./cmd/daily/main.go 2>> "$ROOT_DIR/error.log"

git config user.name "github-actions"
git config user.email "actions@github.com"

git add data/history.jsonl error.log

if ! git diff --cached --quiet; then
    echo "Committing changes..."
    git commit -m "feat(astarisborn): daily update ($(date +%F))"

    echo "Setting remote with token..."
    git remote set-url origin https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git

    echo "Pushing..."
    git push
else
    echo "No changes to commit"
fi

echo "Done at $(date)"
