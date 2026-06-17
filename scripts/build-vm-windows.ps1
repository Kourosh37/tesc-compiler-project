$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$Bin = Join-Path $Root "bin"
$Output = Join-Path $Bin "tsvm.exe"

New-Item -ItemType Directory -Force $Bin | Out-Null

Push-Location $Root
try {
    go build -o $Output .\cmd\tsvm
    Write-Host "Built $Output"
    Write-Host "Usage: .\bin\tsvm.exe .\target\tsvm\path\to\file.tsvm"
} finally {
    Pop-Location
}
