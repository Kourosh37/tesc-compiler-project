# TesLang Compiler Runbook

## Requirements

- Go 1.22+

Check:

```sh
go version
```

## Test

```sh
go test ./...
```

## Run From Source

The compiler reads source code from standard input.

### PowerShell

```powershell
Get-Content .\testdata\codegen_sample.tes | go run ./cmd/teslang --tokens
Get-Content .\testdata\codegen_sample.tes | go run ./cmd/teslang --check
Get-Content .\testdata\codegen_sample.tes | go run ./cmd/teslang --emit-tsvm
```

Save generated TSVM:

```powershell
Get-Content .\hello.tes | go run ./cmd/teslang --emit-tsvm > .\hello.tsvm
```

### Linux, macOS, Git Bash, or CMD

```sh
go run ./cmd/teslang --tokens < testdata/codegen_sample.tes
go run ./cmd/teslang --check < testdata/codegen_sample.tes
go run ./cmd/teslang --emit-tsvm < testdata/codegen_sample.tes
```

Save generated TSVM:

```sh
go run ./cmd/teslang --emit-tsvm < hello.tes > hello.tsvm
```

## Build

### Windows

```powershell
New-Item -ItemType Directory -Force .\bin | Out-Null
go build -o .\bin\teslang.exe .\cmd\teslang
```

Run built compiler:

```powershell
Get-Content .\hello.tes | .\bin\teslang.exe --emit-tsvm > .\hello.tsvm
```

### Linux/macOS

```sh
mkdir -p bin
go build -o bin/teslang ./cmd/teslang
```

Run built compiler:

```sh
./bin/teslang --emit-tsvm < hello.tes > hello.tsvm
```

## Cross-Compile

PowerShell:

```powershell
New-Item -ItemType Directory -Force .\dist | Out-Null
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o .\dist\teslang-windows-amd64.exe .\cmd\teslang
$env:GOOS="linux";   $env:GOARCH="amd64"; go build -o .\dist\teslang-linux-amd64 .\cmd\teslang
$env:GOOS="darwin";  $env:GOARCH="arm64"; go build -o .\dist\teslang-darwin-arm64 .\cmd\teslang
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
```

Linux/macOS:

```sh
mkdir -p dist
GOOS=windows GOARCH=amd64 go build -o dist/teslang-windows-amd64.exe ./cmd/teslang
GOOS=linux   GOARCH=amd64 go build -o dist/teslang-linux-amd64 ./cmd/teslang
GOOS=darwin  GOARCH=arm64 go build -o dist/teslang-darwin-arm64 ./cmd/teslang
```

## Outputs

- `bin/teslang` or `bin/teslang.exe`: built compiler for your current OS
- `dist/*`: cross-compiled release binaries
- `*.tsvm`: generated TSVM intermediate code

`bin/`, `dist/`, and `*.tsvm` are ignored by Git.

## Important

`hello.tsvm` is not a native executable. It is intermediate code. To run it, you need a TSVM runtime/interpreter:

```sh
tsvm hello.tsvm
```
