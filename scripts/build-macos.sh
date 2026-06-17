#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="$ROOT/bin"
COMPILER="$BIN/tesc"
VM="$BIN/tsvm"

mkdir -p "$BIN"

cd "$ROOT"
GOOS=darwin GOARCH="${GOARCH:-arm64}" go build -o "$COMPILER" ./cmd/teslang
GOOS=darwin GOARCH="${GOARCH:-arm64}" go build -o "$VM" ./cmd/tsvm
chmod +x "$COMPILER" "$VM"

echo "Built $COMPILER"
echo "Built $VM"
echo "Usage: ./bin/tesc ./path/to/file.tes"
echo "Run:   ./bin/tsvm ./target/tsvm/path/to/file.tsvm"
