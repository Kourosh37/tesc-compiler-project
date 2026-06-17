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

Default mode emits TSVM. With file inputs, output is written next to each source file as `.tsvm`.

### PowerShell

```powershell
go run ./cmd/teslang .\hello.tes
go run ./cmd/teslang -o .\build\hello.tsvm .\hello.tes
go run ./cmd/teslang --out-dir .\build .\src\a.tes .\src\b.tes
go run ./cmd/teslang --stdout .\hello.tes
go run ./cmd/teslang --check .\hello.tes
go run ./cmd/teslang --tokens .\hello.tes
```

Stdin still works:

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
go run ./cmd/teslang hello.tes
go run ./cmd/teslang -o build/hello.tsvm hello.tes
go run ./cmd/teslang --out-dir build src/a.tes src/b.tes
go run ./cmd/teslang --stdout hello.tes
go run ./cmd/teslang --check hello.tes
go run ./cmd/teslang --tokens hello.tes
```

Stdin still works:

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
.\scripts\build-windows.ps1
```

Run built compiler:

```powershell
.\bin\tesc.exe .\hello.tes
.\bin\tesc.exe -o .\build\hello.tsvm .\hello.tes
```

### Linux/macOS

```sh
./scripts/build-linux.sh
./scripts/build-macos.sh
```

Run built compiler:

```sh
./bin/tesc hello.tes
./bin/tesc -o build/hello.tsvm hello.tes
```

## Cross-Compile

PowerShell:

```powershell
New-Item -ItemType Directory -Force .\dist | Out-Null
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o .\dist\tesc-windows-amd64.exe .\cmd\teslang
$env:GOOS="linux";   $env:GOARCH="amd64"; go build -o .\dist\tesc-linux-amd64 .\cmd\teslang
$env:GOOS="darwin";  $env:GOARCH="arm64"; go build -o .\dist\tesc-darwin-arm64 .\cmd\teslang
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
```

Linux/macOS:

```sh
mkdir -p dist
GOOS=windows GOARCH=amd64 go build -o dist/tesc-windows-amd64.exe ./cmd/teslang
GOOS=linux   GOARCH=amd64 go build -o dist/tesc-linux-amd64 ./cmd/teslang
GOOS=darwin  GOARCH=arm64 go build -o dist/tesc-darwin-arm64 ./cmd/teslang
```

## Outputs

- `bin/tesc` or `bin/tesc.exe`: built compiler for your current OS
- `dist/*`: cross-compiled release binaries
- `*.tsvm`: generated TSVM intermediate code

`bin/`, `dist/`, and `*.tsvm` are ignored by Git.

## Important

`hello.tsvm` is not a native executable. It is intermediate code. To run it, you need a TSVM runtime/interpreter:

```sh
tsvm hello.tsvm
```
