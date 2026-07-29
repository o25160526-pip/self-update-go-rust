# PROGRESS.md — Theo dõi tiến độ thực hiện

> **Cập nhật sau mỗi task hoàn thành.** Task sau phải đọc file này trước khi bắt đầu.
> **Quy ước:** ✅ Hoàn thành | 🔄 Đang làm | ⏳ Chưa bắt đầu | ❌ Bị block | ⏭️ Skipped

---

## Tổng quan

| Phase | Trạng thái | GATE | Ghi chú |
|---|---|---|---|
| P0 — Scaffolding & Môi trường | 🔄 | GATE-0 | Workflows đã có; đang kiểm chứng Windows CI |
| P1 — App Hello Version | 🔄 | GATE-1 | Code/T1 pass; chờ artifact CI và xác nhận Windows T3 |
| P2 — Check Update + Manifest | 🔄 | GATE-2 | Core Go/Rust + key pin đã triển khai; chờ release |
| P3 — Download + Verify + Install + Restart | 🔄 | GATE-3 | Go verify core + Tauri official updater đã tích hợp |
| P4 — Rollback + Health-Check | 🔄 | GATE-4 | State/health core có test; rollback Windows T3 chưa xác nhận |
| P5 — Test + Docs + Báo cáo | ⏳ | GATE-5 | |

---

## Chi tiết từng Task

| Task | Trạng thái | Ngày hoàn thành | Kết quả kiểm chứng | Ghi chú / Issues |
|---|---|---|---|---|
| **P0-T1** Tạo cấu trúc repo | ✅ | 2026-07-29 13:12 | Layout 15 thư mục khớp §4 | |
| **P0-T2** Khởi tạo Rust app | ✅ | 2026-07-29 13:20 | Cargo workspace: logic (12 test pass) + src-tauri | |
| **P0-T3** Cài Go + Go app | ✅ | 2026-07-29 13:20 | Go 1.25.12 + Wails v2.13.0; 11 test pass | |
| **P0-T4** GitHub Actions workflows | ✅ | 2026-07-29 14:30 | actionlint pass; push-main + tag release workflows | workflow_dispatch API bị integration token từ chối, dùng push trigger |
| **P0-T5** GATE-0: CI green | 🔄 | | Chờ kết quả Windows CI sau push | |
| **P1-T1** Go app hello version | ✅ | 2026-07-29 13:30 | build+vet pass; 11 test pass; --version=1.0.0; --print-update-state=checking | |
| **P1-T2** Rust app hello version | ✅ | 2026-07-29 13:30 | 12 logic test pass; frontend+main.rs+version.rs done | |
| **P1-T3** Architecture doc + source-map | ✅ | 2026-07-29 13:35 | docs/architecture.md + docs/source-map.md hoàn chỉnh | |
| **P1-T4** GATE-1: 2 app chạy được | ⏳ | | Chờ user chạy trên Windows 11 [T3] | |
| **P2-T1** Sinh key pair Minisign | ✅ | 2026-07-29 14:25 | Private key local gitignored; public key pin chung Go/Rust | Cần upload GitHub Secret trước release |
| **P2-T2** Go app check update | ✅ | 2026-07-29 14:30 | Manifest HTTPS/channel/platform/semver validation; race tests pass | |
| **P2-T3** Rust app check update | ✅ | 2026-07-29 14:30 | Tauri updater integration + offline mock; 15 logic tests/clippy pass | |
| **P2-T4** Publish manifest v1.0.0 | ⏳ | | | |
| **P2-T5** GATE-2: App phát hiện bản mới | ⏳ | | | |
| **P3-T1** Go app download + verify | ⏳ | | | |
| **P3-T2** Rust app download + verify + install | ⏳ | | | |
| **P3-T3** Luồng update end-to-end | ⏳ | | | |
| **P3-T4** GATE-3: Full flow B01-B07 | ⏳ | | | |
| **P4-T1** Go app health-check + rollback | ⏳ | | | |
| **P4-T2** Rust app health-check + rollback | ⏳ | | | |
| **P4-T3** Test nhóm E | ⏳ | | | |
| **P4-T4** GATE-4: Rollback hoạt động | ⏳ | | | |
| **P5-T1** Test nhóm A | ⏳ | | | |
| **P5-T2** Test nhóm C/D/F/G | ⏳ | | | |
| **P5-T3** Hoàn thiện docs (7 file) | ⏳ | | | |
| **P5-T4** Báo cáo cuối cùng | ⏳ | | | |
| **P5-T5** GATE-5: Definition of Done | ⏳ | | | |

---

## Issues & Lệch hướng

| # | Ngày | Mô tả | Hành động | Trạng thái |
|---|---|---|---|---|
| 1 | 07-29 | Wails v3 không có updater service | Wails v2 + go-selfupdate | ✅ |
| 2 | 07-29 | Agent Linux, không build Windows | CI (T2) + user (T3) | ✅ |
| 3 | 07-29 | Tauri 2 không có built-in rollback | Tự viết adapter | ✅ |
| 4 | 07-29 | Token thiếu `workflows` permission | **Cần user add 2 file workflow thủ công** | ❌ BLOCKED |
| 5 | 07-29 | go-minisign tag không có v prefix | aead.dev/minisign v0.3.0 | ✅ |
