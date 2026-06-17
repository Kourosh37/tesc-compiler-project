#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BIN="$ROOT/bin"
OUTPUT="$BIN/tes"

mkdir -p "$BIN"

cd "$ROOT"
GOOS=linux GOARCH="${GOARCH:-amd64}" go build -o "$OUTPUT" ./cmd/tes
chmod +x "$OUTPUT"

echo "Built $OUTPUT"
echo "Usage: ./bin/tes ./path/to/file.tes"
