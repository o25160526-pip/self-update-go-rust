# Rollback & health-check

## Nguyên tắc

Update chỉ được coi là thành công khi **bản mới khởi động được**. Nếu không, app phải tự về `last-known-good` mà **không mất dữ liệu/cấu hình của user**.

## State trên đĩa

`%LOCALAPPDATA%\go-demo\state.json` (`updater/state.go`):

```json
{
  "current": "1.0.1",
  "pending": "",
  "lastKnownGood": "1.0.1",
  "updatedFrom": "1.0.0",
  "rollbackAttempts": 0,
  "status": "up-to-date",
  "updatedAt": "2026-07-29T08:40:00Z"
}
```

Ghi atomic: ghi `state.json.tmp` rồi `os.Rename`. Máy tắt giữa lúc ghi thì file cũ vẫn nguyên.

## Hai lớp bảo vệ

### Lớp 1 — watchdog (tiến trình cũ)

Sau khi cài, tiến trình **cũ chưa thoát**: nó spawn bản mới rồi chờ file `health/<version>.ok`.

- Thấy marker trong `healthCheckTimeoutSeconds` (mặc định 30s) → thoát code 0, update xong.
- Không thấy (bản mới crash hoặc treo) → `RestoreBackup` đưa `*.backup-<version cũ>` về lại, spawn `exe --rolled-back`, thoát code 1.

Bắt được cả 2 ca: crash ngay lúc start và treo quá timeout.

### Lớp 2 — health-check lúc khởi động

Mọi lần app start đều chạy `Service.Startup()` → `EvaluateStartup()`:

| State trên đĩa | Version đang chạy | Kết quả |
|---|---|---|
| `pending` rỗng | bất kỳ | `normal` — đồng bộ `current` |
| `pending = 1.0.1` | `1.0.1` | `healthy` — commit `lastKnownGood=1.0.1`, xoá `pending`, reset `rollbackAttempts` |
| `pending = 1.0.1` | `1.0.0` | `rolled-back` — `current = lastKnownGood`, `rollbackAttempts++` |

Lớp này bắt cả trường hợp **máy tắt giữa lúc install** (E07): lần khởi động sau, state còn `pending` nhưng version đang chạy là bản cũ → rollback sạch sẽ.

### Chặn rollback vô hạn

`Rollback()` từ chối khi `rollbackAttempts >= maxRollbackAttempts` (mặc định 1). App vẫn chạy, state `failed`, log ghi rõ lý do — không loop restart.

## Không xoá dữ liệu user

Rollback chỉ chạm: file `.exe`, `state.json`, `health/`, `update-tmp/`. Không xoá config hay dữ liệu người dùng.

## Bản Rust

Tauri 2 **không có** built-in health-check/rollback (giới hạn framework, ghi trong [architecture.md](architecture.md)). Adapter tự viết `src-tauri/src/updater/health.rs` + `rollback.rs`, dùng cùng state machine; rollback bằng cách chạy lại installer bản trước đã cache.

## Test liên quan

```powershell
cd go-demo
go test ./updater/ -run "Rollback|Startup|Health" -v
```
