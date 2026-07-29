# Luồng update

## State machine

9 trạng thái, dùng chung cho cả Go và Rust:

```
checking ──► up-to-date
   │
   └──► update-available ──► downloading ──► verifying ──► installing ──► restarting
                                  │              │             │              │
                                  └──────────────┴─────────────┴──► failed    │
                                                                              ▼
                                                        health-check OK ──► up-to-date
                                                        health-check FAIL ► rolled-back
```

Go: `updater/check.go` (hằng `State*`) · `updater/service.go` (chuyển trạng thái).
Rust: `logic/src/state.rs` + `src-tauri/src/updater/`.

## Chi tiết bản Go (Wails v2 — tự viết)

1. **Đọc version build** từ `internal/version.Version` (set bằng `-ldflags -X` lúc build) — nguồn duy nhất.
2. **Đọc policy** (`updater/policy.go`): `autoCheckOnStartup`, `autoDownload`, `autoInstall`, `restartAutomatically`, `allowDowngrade=false`, `healthCheckTimeoutSeconds=30`, `maxRollbackAttempts=1`.
3. **Startup health-check** (`Service.Startup`) — xem [rollback.md](rollback.md).
4. **Fetch manifest** qua HTTPS: `https://github.com/<repo>/releases/latest/download/manifest-go.json`. Bắt buộc HTTPS, giới hạn body 1 MiB.
5. **Validate manifest** (`updater/manifest.go`): `publishedAt` RFC3339, `sha256` đúng 64 hex, `size > 0`, URL scheme `https`, có `signature`. JSON có field lạ → từ chối (`DisallowUnknownFields`).
6. **Lọc** channel (`stable`) + platform (`windows-x86_64`, suy ra từ `runtime.GOOS/GOARCH`).
7. **So semver** (`Masterminds/semver`): chỉ update khi bản mới lớn hơn; `allowDowngrade=false` chặn hạ cấp.
8. Không có bản mới → `up-to-date`. Có → `update-available` + release notes (cắt 200 ký tự).
9. **Download** vào `%LOCALAPPDATA%\go-demo\update-tmp\<version>\` (ghi file `.part` rồi rename).
10. **Verify**: size khớp manifest → SHA-256 khớp → Minisign khớp public key **pin trong code** (`updater/keys.go`). Sai bất kỳ bước nào → xoá file tạm, state `failed`.
11. **Ghi state** `pending=<version mới>`, `updatedFrom=<version cũ>` vào `state.json` (ghi tạm rồi rename → atomic).
12. **Install** (`InstallOverSelf`): rename exe đang chạy thành `*.backup-<version cũ>` (Windows cho rename file đang chạy, chỉ không cho xoá), rồi copy artifact mới vào đúng đường dẫn cũ.
13. **Restart**: spawn `exe --updated-from <version cũ>`. Tiến trình cũ **ở lại làm watchdog**.
14. Bản mới khởi động → ghi `health/<version>.ok`. Watchdog thấy marker trong `healthCheckTimeoutSeconds` → thoát code 0. Không thấy → restore backup, spawn `exe --rolled-back`, thoát code 1.

## Chi tiết bản Rust (Tauri 2 — plugin chính thức)

1. `tauri-plugin-updater` đọc `plugins.updater.endpoints` → `latest-rust.json`.
2. Plugin so version với `tauri.conf.json.version`, verify Minisign bằng `pubkey` trong config.
3. `installMode: passive` → NSIS installer chạy im lặng, tự thay app, không hỏi user.
4. `app.restart()` khởi động lại; adapter tự viết `src-tauri/src/updater/health.rs` + `rollback.rs` lo health-check và rollback.

## Khác biệt cần biết

| | Go | Rust |
|---|---|---|
| Hash trong manifest | có `sha256`, verify cả hash và chữ ký | Tauri chỉ dùng chữ ký |
| Channel / minSupportedVersion | client tự xử lý | Tauri không có khái niệm này |
| Install | thay exe tại chỗ + backup | NSIS installer |
| Rollback | đổi lại file exe (tức thì) | chạy lại installer bản cũ (chậm hơn) |
