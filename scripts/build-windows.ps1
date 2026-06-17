$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$Bin = Join-Path $Root "bin"
$Compiler = Join-Path $Bin "tesc.exe"
$VM = Join-Path $Bin "tesvm.exe"

New-Item -ItemType Directory -Force $Bin | Out-Null

Push-Location $Root
try {
    go build -o $Compiler .\cmd\teslang
    go build -o $VM .\cmd\tesvm
    Write-Host "Built $Compiler"
    Write-Host "Built $VM"
    Write-Host "Usage: .\bin\tesc.exe .\path\to\file.tes"
    Write-Host "Run:   .\bin\tesvm.exe .\target\tesvm\path\to\file.tesvm"
} finally {
    Pop-Location
}
