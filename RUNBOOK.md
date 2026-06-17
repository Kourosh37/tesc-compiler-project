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

Default mode emits TSVM. With file inputs, output is written under `target/tsvm`.

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

Run generated TSVM from source:

```powershell
"5`n7`n" | go run ./cmd/tsvm .\target\tsvm\testdata\codegen_sample.tsvm
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

Run generated TSVM from source:

```sh
printf "5\n7\n" | go run ./cmd/tsvm ./target/tsvm/testdata/codegen_sample.tsvm
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
.\bin\tsvm.exe .\target\tsvm\hello.tsvm
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
./bin/tsvm ./target/tsvm/hello.tsvm
```

## Cross-Compile

PowerShell:

```powershell
New-Item -ItemType Directory -Force .\dist | Out-Null
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o .\dist\tesc-windows-amd64.exe .\cmd\teslang
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o .\dist\tsvm-windows-amd64.exe .\cmd\tsvm
$env:GOOS="linux";   $env:GOARCH="amd64"; go build -o .\dist\tesc-linux-amd64 .\cmd\teslang
$env:GOOS="linux";   $env:GOARCH="amd64"; go build -o .\dist\tsvm-linux-amd64 .\cmd\tsvm
$env:GOOS="darwin";  $env:GOARCH="arm64"; go build -o .\dist\tesc-darwin-arm64 .\cmd\teslang
$env:GOOS="darwin";  $env:GOARCH="arm64"; go build -o .\dist\tsvm-darwin-arm64 .\cmd\tsvm
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
```

Linux/macOS:

```sh
mkdir -p dist
GOOS=windows GOARCH=amd64 go build -o dist/tesc-windows-amd64.exe ./cmd/teslang
GOOS=windows GOARCH=amd64 go build -o dist/tsvm-windows-amd64.exe ./cmd/tsvm
GOOS=linux   GOARCH=amd64 go build -o dist/tesc-linux-amd64 ./cmd/teslang
GOOS=linux   GOARCH=amd64 go build -o dist/tsvm-linux-amd64 ./cmd/tsvm
GOOS=darwin  GOARCH=arm64 go build -o dist/tesc-darwin-arm64 ./cmd/teslang
GOOS=darwin  GOARCH=arm64 go build -o dist/tsvm-darwin-arm64 ./cmd/tsvm
```

## Outputs

- `bin/tesc` or `bin/tesc.exe`: built compiler for your current OS
- `bin/tsvm` or `bin/tsvm.exe`: built VM for your current OS
- `dist/*`: cross-compiled release binaries
- `target/tsvm/**/*.tsvm`: default generated TSVM intermediate code
- `*.tsvm`: generated TSVM files from custom output paths

`bin/`, `dist/`, `target/`, and `*.tsvm` are ignored by Git.

## Important

`hello.tsvm` is not a native executable. Run it with the bundled TSVM runtime:

```sh
tsvm hello.tsvm
```
