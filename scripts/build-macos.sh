#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="$ROOT/bin"
COMPILER="$BIN/tesc"
VM="$BIN/tesvm"

mkdir -p "$BIN"

cd "$ROOT"
GOOS=darwin GOARCH="${GOARCH:-arm64}" go build -o "$COMPILER" ./cmd/teslang
GOOS=darwin GOARCH="${GOARCH:-arm64}" go build -o "$VM" ./cmd/tesvm
chmod +x "$COMPILER" "$VM"

echo "Built $COMPILER"
echo "Built $VM"
echo "Usage: ./bin/tesc ./path/to/file.tes"
echo "Run:   ./bin/tesvm ./target/tesvm/path/to/file.tesvm"
