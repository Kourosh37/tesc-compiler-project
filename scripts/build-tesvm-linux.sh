#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="$ROOT/bin"
OUTPUT="$BIN/tesvm"

mkdir -p "$BIN"

cd "$ROOT"
GOOS=linux GOARCH="${GOARCH:-amd64}" go build -o "$OUTPUT" ./cmd/tesvm
chmod +x "$OUTPUT"

echo "Built $OUTPUT"
echo "Usage: ./bin/tesvm ./target/tesvm/path/to/file.tesvm"
