# PROGRESS.md — Theo dõi tiến độ thực hiện

> Cập nhật sau mỗi task hoàn thành.
> **Mirror trong ClickUp:** doc "History Work — self-update-go-rust" (trang 📊 PROGRESS + 🕒 Work Log + 📥 Yêu cầu từ user).
> Quy ước: ✅ Hoàn thành | 🔄 Đang làm | ⏳ Chưa bắt đầu | ❌ Bị block | ⏭️ Skipped

---

## Tổng quan

| Phase | Trạng thái | GATE | Ghi chú |
|---|---|---|---|
| P0 — Scaffolding & Môi trường | ✅ | GATE-0 | CI xanh cả Go + Rust trên windows-latest |
| P1 — App Hello Version | 🔄 | GATE-1 | Code xong; đang phát hành v1.0.2 để user chạy Windows 11 |
| P2 — Check Update + Manifest | ✅ | GATE-2 | v1.0.0/v1.0.1 đã có release + manifest |
| P3 — Download + Verify + Install + Restart | 🔄 | GATE-3 | T1 pass; T3 đang kiểm chứng bằng v1.0.2 |
| P4 — Rollback + Health-Check | 🔄 | GATE-4 | T1 pass; T3 cần test trên Windows 11 |
| P5 — Test + Docs + Báo cáo | 🔄 | GATE-5 | Docs hoàn chỉnh; chờ kết quả T3 |

## Releases đã tạo

| Release | Go | Rust | Mục đích | Trạng thái |
|---|---|---|---|---|
| v1.0.0 | ✅ | ✅ | Baseline đầu tiên | ✅ Published |
| v1.0.1 | ✅ | ✅ | Bản update trung gian | ✅ Published |
| v1.0.2 | 🔄 | 🔄 | Bản user tải về để kiểm chứng update thật | CI đang chạy |

## Chi tiết từng Task

| Task | Trạng thái | Kết quả kiểm chứng | Ghi chú |
|---|---|---|---|
| P0-T1..T5 | ✅ | Layout + CI Windows xanh | |
| P1-T1 Go hello version | ✅ | `--version` được build theo version release | |
| P1-T2 Rust hello version | ✅ | Tauri version được workflow cập nhật theo branch | |
| P1-T3 Architecture + source map | ✅ | 2 docs hoàn chỉnh | |
| P1-T4 GATE-1 Windows 11 | 🔄 | v1.0.2 sẽ là bản test | Chờ user xác nhận |
| P2-T1 Key pair Minisign | ✅ | Demo key ID `5756DB5A3509A8C1` trong `keys/` | Production dùng Secrets |
| P2-T2 Go check update | ✅ | HTTPS, channel, platform, semver, downgrade guard | |
| P2-T3 Rust check update | ✅ | Tauri updater + pinned public key | |
| P2-T4 Publish v1.0.0 | ✅ | Go/Rust assets + manifest/signature/checksum | |
| P2-T5 GATE-2 | ✅ | v1.0.0 → v1.0.1 release path tồn tại | Cần user chạy client cũ |
| P3-T1 Go download + verify | ✅ | Size + SHA-256 + Minisign, test T1 | |
| P3-T2 Rust download + install | ✅ | NSIS + `.sig`, `installMode=passive` | |
| P3-T3 End-to-end | ✅ | Mock HTTPS test: download → verify → install → restart → health | T3 v1.0.2 pending |
| P3-T4 GATE-3 B01-B07 | 🔄 | B01-B05/B07 T1 pass | B06 cần Windows 11 |
| P4-T1 Go health + rollback | ✅ | Watchdog + backup + max one rollback | |
| P4-T2 Rust health + rollback adapter | 🔄 | Code/test logic có sẵn | Cần Windows 11 |
| P4-T3 Test nhóm E | ✅ | E01-E03/E05-E08 T1 pass, E04 N/A | |
| P4-T4 GATE-4 | 🔄 | T1 pass | T3 pending |
| P5-T1 Test nhóm A | 🔄 | A01/A02/A06 pass | A03-A05 pending user |
| P5-T2 Test nhóm C/D/F/G | 🔄 | C/D pass, G v1.0.0/v1.0.1 pass | F + v1.0.2 pending |
| P5-T3 Docs | ✅ | 8 docs + README + keys guide | |
| P5-T4 Final report | ⏳ | Chờ kết quả Windows 11 | |
| P5-T5 GATE-5 | ⏳ | Chờ T3 | |

## Quy trình kiểm chứng v1.0.2

1. Chờ cả 2 workflow trên nhánh `release/v1.0.2` xanh.
2. Kiểm tra Release v1.0.2 có đủ Go `.exe/.sig/.sha256/manifest-go.json` và Rust installer `.exe/.sig/.sha256/latest-rust.json`.
3. Trên Windows 11, tải Go/Rust v1.0.0 hoặc v1.0.1 và chạy thử update lên v1.0.2.
4. Ghi version hiển thị, trạng thái updater, log, UAC/SmartScreen, kết quả restart vào ClickUp History Work.
5. Nếu update lỗi, giữ nguyên state/log, sửa code, tăng patch version tiếp theo rồi lặp CI.

## Issues & Lệch hướng

| # | Mô tả | Hành động | Trạng thái |
|---|---|---|---|
| 1 | Wails v3 không có updater service | Wails v2 + tự viết staged install | ✅ |
| 2 | Agent Linux không build Windows | CI + user T3 | ✅ |
| 3 | Tauri 2 không có built-in rollback | Adapter tự viết | ✅ |
| 6 | Go embed dist thiếu file | `.gitkeep` + build frontend trước test | ✅ |
| 7 | Rust thiếu icon.ico | `build.rs` tự sinh icon | ✅ |
| 9 | Private key cũ mất | Demo key trong repo + secret fallback | ✅ |
| 10 | Token không push tag | Nhánh `release/v*` tự tạo release/tag | ✅ |
| 11 | Demo scrypt quá nặng | Re-encrypt nhẹ, public key không đổi | ✅ |
| 12 | Sign/verify khác implementation | Go ký và verify cùng `aead.dev/minisign` | ✅ |
