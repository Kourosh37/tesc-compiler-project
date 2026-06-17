$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$Bin = Join-Path $Root "bin"
$Output = Join-Path $Bin "tesvm.exe"

New-Item -ItemType Directory -Force $Bin | Out-Null

Push-Location $Root
try {
    go build -o $Output .\cmd\tesvm
    Write-Host "Built $Output"
    Write-Host "Usage: .\bin\tesvm.exe .\target\tesvm\path\to\file.tesvm"
} finally {
    Pop-Location
}
