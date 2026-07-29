# self-update-go-rust

Reusable Windows 11 desktop auto-update template for Go/Wails and Rust/Tauri applications.

## Current status

✅ v1.0.3 Go and Rust validated by the user on Windows 11. Go update and Rust UI/update path are working. The repo is now being hardened as a reusable product template.

## What product forks inherit

Signed manifests, SHA-256 and Minisign verification, staged Go install, Tauri updater install, restart, health-check, rollback, Windows CI, release automation, and detailed deployment docs.

## Start a product

Read [docs/template-guide.md](docs/template-guide.md), copy [shared/update-server.example.json](shared/update-server.example.json), configure a host, rotate the demo key, then add business modules. Do not rewrite updater internals for product features.

## Update hosts

Go loads `update-server.json` beside the executable or the path in `GO_DEMO_UPDATE_SERVER_CONFIG`; invalid/missing config falls back to the built-in GitHub Releases endpoint. Rust endpoints are compiled into Tauri config and should be supplied through a product release config overlay.

## Documentation

- [Deployment](docs/deployment.md)
- [Update server config](docs/server-config.md)
- [Template guide](docs/template-guide.md)
- [Next roadmap](docs/roadmap-next.md)
- [Windows build](docs/windows-11-build.md)
- [Release](docs/release.md)
- [Update flow](docs/update-flow.md)
- [Rollback](docs/rollback.md)
- [Testing](docs/testing.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Architecture](docs/architecture.md)
- [Process](plan/PROGRESS.md)

## Demo releases

Download the latest Go portable executable or Rust NSIS installer from [Releases](https://github.com/o25160526-pip/self-update-go-rust/releases). The demo key is public and must be replaced before production.
