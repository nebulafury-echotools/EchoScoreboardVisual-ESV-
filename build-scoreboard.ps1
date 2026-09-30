$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $projectRoot
$buildDirectory = Join-Path $projectRoot 'builds'
$exePath = Join-Path $buildDirectory 'scoreboard.exe'
New-Item -ItemType Directory -Path $buildDirectory -Force | Out-Null
Write-Host "Building $exePath..."
& go build -o $exePath .\cmd\scoreboard
if ($LASTEXITCODE -ne 0) {
    throw "Build failed."
}
Write-Host "Build complete: $exePath"
