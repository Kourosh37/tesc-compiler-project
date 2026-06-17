#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BIN="$ROOT/bin"

mkdir -p "$BIN"

cd "$ROOT"
GOOS=linux GOARCH="${GOARCH:-amd64}" go build -o "$BIN/tes" ./cmd/tes
GOOS=linux GOARCH="${GOARCH:-amd64}" go build -o "$BIN/tesc" ./cmd/teslang
GOOS=linux GOARCH="${GOARCH:-amd64}" go build -o "$BIN/tesvm" ./cmd/tesvm
chmod +x "$BIN/tes" "$BIN/tesc" "$BIN/tesvm"

echo "Built $BIN/tes"
echo "Built $BIN/tesc"
echo "Built $BIN/tesvm"
echo "Run:          ./bin/tes ./path/to/file.tes"
echo "Compile only: ./bin/tesc ./path/to/file.tes"
echo "VM only:      ./bin/tesvm ./target/tesvm/path/to/file.tesvm"
