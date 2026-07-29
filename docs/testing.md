# Testing — 7 nhóm A-G

Mức kiểm chứng: **T1** tự động (unit/integration) · **T2** CI `windows-latest` · **T3** máy Windows 11 thật của user.

```powershell
# T1 - Go
cd go-demo; go test ./... -race -v

# T1 - Rust
cd rust-demo; cargo test --manifest-path logic\Cargo.toml
```

## A — App cơ bản

| Mã | Nội dung | Mức | Cách kiểm | Trạng thái |
|---|---|---|---|---|
| A01 | Clean build Go | T2 | workflow release-go | ✅ |
| A02 | Clean build Rust | T2 | workflow release-rust | ✅ |
| A03 | App hiện đúng version | T3 | mở app | ⏳ user |
| A04 | `--version` đúng | T3 | `.\go-demo.exe --version` | ⏳ user |
| A05 | App khởi động được khi offline | T3 | tắt mạng rồi mở app | ⏳ user |
| A06 | Không ghi secret vào log | T1 | log chỉ ghi key ID | ✅ |

## B — Luồng update thành công

| Mã | Nội dung | Mức | Cách kiểm | Trạng thái |
|---|---|---|---|---|
| B01 | Không có bản mới → `up-to-date` | T1 | `TestServiceUpToDate` | ✅ |
| B02 | Có bản mới → `update-available` | T1 | `TestServiceUpdateThanhCong` | ✅ |
| B03 | Tự download không cần bấm | T1 | cùng test trên | ✅ |
| B04 | Hash + chữ ký đúng → cài được | T1 | cùng test trên | ✅ |
| B05 | App tự restart | T1 | hook Spawn được gọi | ✅ |
| B06 | Sau restart thấy version mới | T3 | update thật từ release | ⏳ user |
| B07 | State ghi `lastKnownGood` mới | T1 | `newProc.Startup()` | ✅ |

## C — Lỗi mạng

| Mã | Nội dung | Mức | Trạng thái |
|---|---|---|---|
| C01 | DNS fail / không có mạng | T1 | ✅ `Check` trả lỗi, state `failed` |
| C02 | HTTP 404 manifest | T1 | ✅ `fetch manifest: HTTP 404` |
| C03 | HTTP 500 | T1 | ✅ cùng nhánh xử lý |
| C04 | Manifest không phải JSON | T1 | ✅ `parse manifest` |
| C05 | Download bị ngắt giữa đường | T1 | ✅ size mismatch → reject |
| C06 | Artifact lớn hơn size khai báo | T1 | ✅ `io.LimitReader(size+1)` |
| C07 | Timeout | T1 | ✅ `http.Client{Timeout: 30s}` |

## D — Bảo mật

| Mã | Nội dung | Mức | Trạng thái |
|---|---|---|---|
| D01 | SHA-256 sai → từ chối | T1 | ✅ |
| D02 | Chữ ký sai → từ chối | T1 | ✅ `TestDemoKeyKhopPinnedPublicKey` |
| D03 | Manifest bị sửa → từ chối | T1 | ✅ `DisallowUnknownFields` + validate |
| D04 | Version thấp hơn → từ chối | T1 | ✅ `allowDowngrade=false` |
| D05 | Sai OS/arch → từ chối | T1 | ✅ lọc `PlatformString()` |
| D06 | Private key không lọt vào log/artifact | T1 | ✅ |
| D07 | Client không cần GitHub token | T1 | ✅ không gửi Authorization |
| D08 | URL không phải HTTPS → từ chối | T1 | ✅ `Manifest.Validate` |

## E — Restart & rollback

| Mã | Nội dung | Mức | Trạng thái |
|---|---|---|---|
| E01 | App đang giữ file exe | T1 | ✅ rename-then-copy |
| E02 | Bản mới crash trước health-check | T1 | ✅ `TestServiceRollbackKhiBanMoiKhongLen` |
| E03 | Bản mới treo quá timeout | T1 | ✅ cùng test |
| E04 | Local service | — | N/A trong demo |
| E05 | Rollback về `lastKnownGood` | T1 | ✅ `TestEvaluateStartupRollback` |
| E06 | Không rollback lặp vô hạn | T1 | ✅ `TestEvaluateStartupKhongRollbackLapVoHan` |
| E07 | Tắt máy giữa lúc install | T1 | ✅ state atomic + health-check khi khởi động |
| E08 | Rollback không xoá dữ liệu user | T1 | ✅ chỉ chạm exe + state |

## F — Quyền Windows

| Mã | Nội dung | Mức | Trạng thái |
|---|---|---|---|
| F01 | Bản Go chạy không cần admin | T3 | ⏳ user |
| F02 | UAC khi cài bản Rust (NSIS) | T3 | ⏳ user |
| F03 | SmartScreen cảnh báo | T3 | ⏳ user |
| F04 | Antivirus không chặn update | T3 | ⏳ user |
| F05 | Ghi được vào `%LOCALAPPDATA%` | T3 | ⏳ user |

## G — Phát hành

| Mã | Nội dung | Mức | Trạng thái |
|---|---|---|---|
| G01 | Nhánh `release/v1.0.0` → tạo release | T2 | 🔄 |
| G02 | `release/v1.0.1` → artifact mới | T2 | 🔄 |
| G03 | Artifact Go và Rust không ghi đè nhau | T2 | ✅ |
| G04 | Upload checksum | T2 | ✅ |
| G05 | Manifest trỏ đúng artifact | T2 | ✅ `verify-artifact` chạy trong CI |
| G06 | Ký được cả khi không có secret | T2 | ✅ fallback demo key |
| G07 | Release có release notes | T2 | ✅ |
