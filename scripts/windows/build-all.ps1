$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$Bin = Join-Path $Root "bin"

New-Item -ItemType Directory -Force $Bin | Out-Null

Push-Location $Root
try {
    go build -o (Join-Path $Bin "tes.exe") .\cmd\tes
    go build -o (Join-Path $Bin "tesc.exe") .\cmd\teslang
    go build -o (Join-Path $Bin "tesvm.exe") .\cmd\tesvm

    Write-Host "Built $(Join-Path $Bin "tes.exe")"
    Write-Host "Built $(Join-Path $Bin "tesc.exe")"
    Write-Host "Built $(Join-Path $Bin "tesvm.exe")"
    Write-Host "Run:          .\bin\tes.exe .\path\to\file.tes"
    Write-Host "Compile only: .\bin\tesc.exe .\path\to\file.tes"
    Write-Host "VM only:      .\bin\tesvm.exe .\target\tesvm\path\to\file.tesvm"
} finally {
    Pop-Location
}
