# Source Map — Mapping Go ↔ Rust

> Mục đích: đảm bảo hai app có cấu trúc tương đương, dễ so sánh.
> Cập nhật khi layout thay đổi.

## Thư mục

| Ý nghĩa | Go (Wails v2) | Rust (Tauri 2) |
|---|---|---|
| Version nguồn duy nhất | `go-demo/internal/version/version.go` | `rust-demo/src-tauri/src/version.rs` |
| UI hello version | `go-demo/frontend/` | `rust-demo/frontend/` |
| Check update | `go-demo/updater/check.go` | `rust-demo/src-tauri/src/updater/check.rs` |
| Download/verify/install | `go-demo/updater/install.go` | `rust-demo/src-tauri/src/updater/install.rs` |
| Restart | `go-demo/updater/restart.go` | `rust-demo/src-tauri/src/updater/restart.rs` |
| Health-check | `go-demo/updater/health.go` | `rust-demo/src-tauri/src/updater/health.rs` |
| Rollback | `go-demo/updater/rollback.go` | `rust-demo/src-tauri/src/updater/rollback.rs` |
| Policy | `go-demo/updater/policy.go` | `rust-demo/src-tauri/src/updater/policy.rs` |
| Manifest parser | `go-demo/updater/manifest.go` | `rust-demo/src-tauri/src/updater/manifest.rs` |
| Mock provider | `go-demo/updater/mock.go` | `rust-demo/src-tauri/src/updater/mock.rs` |
| Launcher/installer helper | `go-demo/installer/` | `rust-demo/installer/` |
| Build script | `go-demo/scripts/build.ps1` | `rust-demo/scripts/build.ps1` |
| Release workflow | `.github/workflows/release-go.yml` | `.github/workflows/release-rust.yml` |
| App entrypoint | `go-demo/app/main.go` | `rust-demo/src-tauri/src/main.rs` |
| Tauri config | (N/A) | `rust-demo/src-tauri/tauri.conf.json` |
| Cargo.toml / go.mod | `go-demo/go.mod` | `rust-demo/src-tauri/Cargo.toml` |
| Capabilities | (N/A — Wails dùng manifest app) | `rust-demo/src-tauri/capabilities/default.json` |

## Cơ chế install/rollback khác nhau

| | Go (Wails v2) | Rust (Tauri 2) |
|---|---|---|
| Updater core | `creativeprojects/go-selfupdate v1.6.0` (check + download) + tự viết staged install/restart/rollback | `tauri-plugin-updater v2.10.1` (check + download + verify + install + restart) |
| Install mechanism | **Launcher + versioned dirs + pointer file** (tự viết) | **NSIS installer** do Tauri bundler tạo, `installMode: passive` |
| Tránh ghi đè exe đang chạy | Mỗi version 1 thư mục riêng, chỉ đổi con trỏ | NSIS: app thoát rồi installer mới chạy |
| Rollback | Đổi pointer về `last-known-good` (atomic, tức thì) | Chạy lại installer của version trước (đã cache) |
| Health-check | Tự viết | Tự viết (Tauri không có built-in) |
| Signature | Minisign (Ed25519) qua `go-minisign v0.2.7` | Minisign (Ed25519) qua `tauri-plugin-updater` |

## Giới hạn framework đã ghi nhận

| # | Giới hạn | Hành động |
|---|---|---|
| 1 | Wails v3 vẫn Alpha, không có updater service | Dùng Wails v2 stable + go-selfupdate |
| 2 | Tauri 2 không có built-in rollback | Tự viết adapter `rollback.rs`, có test |
| 3 | Tauri 2 không có built-in health-check | Tự viết adapter `health.rs`, có test |
