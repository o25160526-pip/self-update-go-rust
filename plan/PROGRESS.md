# PROGRESS.md — Theo dõi tiến độ thực hiện

> **Cập nhật sau mỗi task hoàn thành.** Task sau phải đọc file này trước khi bắt đầu.
> **Mirror trong ClickUp:** doc "History Work — self-update-go-rust" (trang 📊 PROGRESS, 🕒 Work Log, 📥 Yêu cầu từ user).
> **Quy ước:** ✅ Hoàn thành | 🔄 Đang làm | ⏳ Chưa bắt đầu | ❌ Bị block | ⏭️ Skipped

---

## Tổng quan

| Phase | Trạng thái | GATE | Ghi chú |
|---|---|---|---|
| P0 — Scaffolding & Môi trường | ✅ | GATE-0 | CI xanh cả Go + Rust trên windows-latest |
| P1 — App Hello Version | 🔄 | GATE-1 | Code xong; chờ user xác nhận trên Windows 11 [T3] |
| P2 — Check Update + Manifest | 🔄 | GATE-2 | Key + fallback secret xong; chờ release đầu tiên |
| P3 — Download + Verify + Install + Restart | 🔄 | GATE-3 | Go: đủ luồng + test T1; Rust: Tauri updater |
| P4 — Rollback + Health-Check | 🔄 | GATE-4 | Go: watchdog + health-check + rollback có test T1 |
| P5 — Test + Docs + Báo cáo | 🔄 | GATE-5 | 8 docs xong; báo cáo cuối chờ kết quả T3 |

---

## Chi tiết từng Task

| Task | Trạng thái | Ngày | Kết quả kiểm chứng | Ghi chú |
|---|---|---|---|---|
| **P0-T1** Tạo cấu trúc repo | ✅ | 07-29 | Layout khớp §4 | |
| **P0-T2** Khởi tạo Rust app | ✅ | 07-29 | Cargo workspace: logic + src-tauri | |
| **P0-T3** Cài Go + Go app | ✅ | 07-29 | Go 1.25.12 + Wails v2.13.0 | |
| **P0-T4** GitHub Actions workflows | ✅ | 07-29 | actionlint pass | |
| **P0-T5** GATE-0: CI green | ✅ | 07-29 | Go run #5 + Rust run #4 success | |
| **P1-T1** Go app hello version | ✅ | 07-29 | `--version` → 1.0.0 | |
| **P1-T2** Rust app hello version | ✅ | 07-29 | 15 test logic pass | |
| **P1-T3** Architecture doc + source-map | ✅ | 07-29 | 2 file hoàn chỉnh | |
| **P1-T4** GATE-1: 2 app chạy được | ⏳ | | Chờ user chạy trên Windows 11 [T3] | |
| **P2-T1** Key pair Minisign | ✅ | 07-29 | Demo key trong `keys/` (ID 5756DB5A3509A8C1), CI fallback secret → repo key | Issue #9, #11 |
| **P2-T2** Go app check update | ✅ | 07-29 | Validate manifest HTTPS/channel/platform/semver | |
| **P2-T3** Rust app check update | ✅ | 07-29 | Tauri updater + mock offline | |
| **P2-T4** Publish manifest v1.0.0 | 🔄 | | Push nhánh `release/v1.0.0` → CI tự tạo release | Issue #10 |
| **P2-T5** GATE-2: phát hiện bản mới | ⏳ | | | |
| **P3-T1** Go app download + verify | ✅ | 07-29 | size + SHA-256 + Minisign, có test T1 | |
| **P3-T2** Rust app download + verify + install | ✅ | 07-29 | Tauri updater, installMode passive | |
| **P3-T3** Luồng update end-to-end | ✅ | 07-29 | `TestServiceUpdateThanhCong`: ký bằng demo key thật → HTTPS mock → verify → install → restart → healthy | Chờ T3 |
| **P3-T4** GATE-3: B01-B07 | 🔄 | | B01-B05, B07 pass T1; B06 chờ T3 | |
| **P4-T1** Go app health-check + rollback | ✅ | 07-29 | Watchdog + health marker + max 1 rollback, có test T1 | |
| **P4-T2** Rust app health-check + rollback | 🔄 | | Adapter tự viết trong src-tauri/src/updater | |
| **P4-T3** Test nhóm E | ✅ | 07-29 | E01-E03, E05-E08 pass T1; E04 N/A | |
| **P4-T4** GATE-4: Rollback hoạt động | 🔄 | | Pass T1, chờ T3 | |
| **P5-T1** Test nhóm A | 🔄 | | A01, A02, A06 ✅; A03-A05 chờ T3 | |
| **P5-T2** Test nhóm C/D/F/G | 🔄 | | C, D pass T1; F chờ T3; G chờ release | |
| **P5-T3** Hoàn thiện docs | ✅ | 07-29 | architecture, source-map, update-flow, rollback, release, windows-11-build, testing, troubleshooting + README + keys/README | |
| **P5-T4** Báo cáo cuối cùng | ⏳ | | Chờ kết quả T3 của user | |
| **P5-T5** GATE-5: Definition of Done | ⏳ | | | |

---

## Issues & Lệch hướng

| # | Ngày | Mô tả | Hành động | Trạng thái |
|---|---|---|---|---|
| 1 | 07-29 | Wails v3 không có updater service | Wails v2 + tự viết staged install | ✅ |
| 2 | 07-29 | Agent Linux, không build Windows | CI (T2) + user (T3) | ✅ |
| 3 | 07-29 | Tauri 2 không có built-in rollback | Tự viết adapter | ✅ |
| 4 | 07-29 | Token thiếu `workflows` permission | Push qua GitHub App integration | ✅ |
| 5 | 07-29 | go-minisign tag không có v prefix | `aead.dev/minisign v0.3.0` | ✅ |
| 6 | 07-29 | CI Go fail: `all:frontend/dist` không có file | Commit `.gitkeep` + build frontend trước test | ✅ |
| 7 | 07-29 | CI Rust fail: thiếu `icon.ico` | `build.rs` tự sinh icon | ✅ |
| 8 | 07-29 | Log CI không đọc được qua API | In 60-80 dòng cuối ra `::error::` annotation | ✅ |
| 9 | 07-29 | **Private key phiên trước sinh trong sandbox rồi mất** → public key pin thành mồ côi, không ký được gì | Sinh keypair demo mới, **commit vào `keys/`**; CI ưu tiên secret `TAURI_SIGNING_PRIVATE_KEY`, không có thì dùng demo key. Xoá key mồ côi `238E4D81EA7CD689` | ✅ |
| 10 | 07-29 | Integration token **không push được tag** → không trigger được release | Thêm trigger nhánh `release/v*`; workflow tự `gh release create` (tạo luôn tag) | ✅ |
| 11 | 07-29 | Demo key ban đầu dùng scrypt "sensitive" (N=65536, p=16) → `go test -race` trên Windows treo hàng phút khi giải mã key | Re-encrypt đúng keypair đó với tham số nhẹ (N=1024, r=8, p=1), public key không đổi; test cache key bằng `sync.Once`; thêm `-timeout 300s` | ✅ |
| 12 | 07-29 | Ký bản Go bằng `npx tauri signer sign` nhưng verify bằng `aead.dev/minisign` → rủi ro lệch implementation | Bản Go tự ký bằng `cmd/sign-artifact` (cùng lib với verify), bỏ phụ thuộc npm khỏi workflow Go | ✅ |
