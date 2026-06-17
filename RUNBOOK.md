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

Default mode emits TESVM. With file inputs, output is written under `target/tesvm`.

### PowerShell

```powershell
go run ./cmd/teslang .\hello.tes
go run ./cmd/teslang -o .\build\hello.tesvm .\hello.tes
go run ./cmd/teslang --out-dir .\build .\src\a.tes .\src\b.tes
go run ./cmd/teslang --stdout .\hello.tes
go run ./cmd/teslang --check .\hello.tes
go run ./cmd/teslang --tokens .\hello.tes
```

Stdin still works:

```powershell
Get-Content .\testdata\codegen_sample.tes | go run ./cmd/teslang --tokens
Get-Content .\testdata\codegen_sample.tes | go run ./cmd/teslang --check
Get-Content .\testdata\codegen_sample.tes | go run ./cmd/teslang --emit-tesvm
```

Run generated TESVM from source:

```powershell
"5`n7`n" | go run ./cmd/tesvm .\target\tesvm\testdata\codegen_sample.tesvm
```

Save generated TESVM:

```powershell
Get-Content .\hello.tes | go run ./cmd/teslang --emit-tesvm > .\hello.tesvm
```

### Linux, macOS, Git Bash, or CMD

```sh
go run ./cmd/teslang hello.tes
go run ./cmd/teslang -o build/hello.tesvm hello.tes
go run ./cmd/teslang --out-dir build src/a.tes src/b.tes
go run ./cmd/teslang --stdout hello.tes
go run ./cmd/teslang --check hello.tes
go run ./cmd/teslang --tokens hello.tes
```

Stdin still works:

```sh
go run ./cmd/teslang --tokens < testdata/codegen_sample.tes
go run ./cmd/teslang --check < testdata/codegen_sample.tes
go run ./cmd/teslang --emit-tesvm < testdata/codegen_sample.tes
```

Run generated TESVM from source:

```sh
printf "5\n7\n" | go run ./cmd/tesvm ./target/tesvm/testdata/codegen_sample.tesvm
```

Save generated TESVM:

```sh
go run ./cmd/teslang --emit-tesvm < hello.tes > hello.tesvm
```

## Build

### Windows

```powershell
.\scripts\build-windows.ps1
.\scripts\build-vm-windows.ps1
```

Run built compiler:

```powershell
.\bin\tesc.exe .\hello.tes
.\bin\tesc.exe -o .\build\hello.tesvm .\hello.tes
.\bin\tesvm.exe .\target\tesvm\hello.tesvm
```

### Linux/macOS

```sh
./scripts/build-linux.sh
./scripts/build-macos.sh
./scripts/build-vm-linux.sh
./scripts/build-vm-macos.sh
```

Run built compiler:

```sh
./bin/tesc hello.tes
./bin/tesc -o build/hello.tesvm hello.tes
./bin/tesvm ./target/tesvm/hello.tesvm
```

## Cross-Compile

PowerShell:

```powershell
New-Item -ItemType Directory -Force .\dist | Out-Null
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o .\dist\tesc-windows-amd64.exe .\cmd\teslang
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o .\dist\tesvm-windows-amd64.exe .\cmd\tesvm
$env:GOOS="linux";   $env:GOARCH="amd64"; go build -o .\dist\tesc-linux-amd64 .\cmd\teslang
$env:GOOS="linux";   $env:GOARCH="amd64"; go build -o .\dist\tesvm-linux-amd64 .\cmd\tesvm
$env:GOOS="darwin";  $env:GOARCH="arm64"; go build -o .\dist\tesc-darwin-arm64 .\cmd\teslang
$env:GOOS="darwin";  $env:GOARCH="arm64"; go build -o .\dist\tesvm-darwin-arm64 .\cmd\tesvm
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
```

Linux/macOS:

```sh
mkdir -p dist
GOOS=windows GOARCH=amd64 go build -o dist/tesc-windows-amd64.exe ./cmd/teslang
GOOS=windows GOARCH=amd64 go build -o dist/tesvm-windows-amd64.exe ./cmd/tesvm
GOOS=linux   GOARCH=amd64 go build -o dist/tesc-linux-amd64 ./cmd/teslang
GOOS=linux   GOARCH=amd64 go build -o dist/tesvm-linux-amd64 ./cmd/tesvm
GOOS=darwin  GOARCH=arm64 go build -o dist/tesc-darwin-arm64 ./cmd/teslang
GOOS=darwin  GOARCH=arm64 go build -o dist/tesvm-darwin-arm64 ./cmd/tesvm
```

## Outputs

- `bin/tesc` or `bin/tesc.exe`: built compiler for your current OS
- `bin/tesvm` or `bin/tesvm.exe`: built VM for your current OS
- `dist/*`: cross-compiled release binaries
- `target/tesvm/**/*.tesvm`: default generated TESVM intermediate code
- `*.tesvm`: generated TESVM files from custom output paths

`bin/`, `dist/`, `target/`, and `*.tesvm` are ignored by Git.

## Important

`hello.tesvm` is not a native executable. Run it with the bundled TESVM runtime:

```sh
tesvm hello.tesvm
```
