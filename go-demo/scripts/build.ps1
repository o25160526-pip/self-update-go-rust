param([string]$Version = "1.0.0")
$ErrorActionPreference = "Stop"
Push-Location $PSScriptRoot/..
try {
  Push-Location frontend; npm ci; npm run build; Pop-Location
  wails build -clean -ldflags "-X go-demo/internal/version.Version=$Version"
} finally { Pop-Location }
