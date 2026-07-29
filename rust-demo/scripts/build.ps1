param([string]$Version = "1.0.0")
$ErrorActionPreference = "Stop"
Push-Location $PSScriptRoot/..
try {
  npm ci
  $env:CARGO_PKG_VERSION = $Version
  npm run tauri build -- --target x86_64-pc-windows-msvc
} finally { Pop-Location }
