$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$Bin = Join-Path $Root "bin"
$Output = Join-Path $Bin "tesc.exe"

New-Item -ItemType Directory -Force $Bin | Out-Null

Push-Location $Root
try {
    go build -o $Output .\cmd\teslang
    Write-Host "Built $Output"
    Write-Host "Usage: .\bin\tesc.exe .\path\to\file.tes"
} finally {
    Pop-Location
}
