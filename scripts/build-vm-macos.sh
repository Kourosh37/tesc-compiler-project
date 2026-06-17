#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="$ROOT/bin"
OUTPUT="$BIN/tsvm"

mkdir -p "$BIN"

cd "$ROOT"
GOOS=darwin GOARCH="${GOARCH:-arm64}" go build -o "$OUTPUT" ./cmd/tsvm
chmod +x "$OUTPUT"

echo "Built $OUTPUT"
echo "Usage: ./bin/tsvm ./target/tsvm/path/to/file.tsvm"
