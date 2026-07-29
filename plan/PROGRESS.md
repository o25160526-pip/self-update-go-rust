# PROGRESS.md — Theo dõi tiến độ thực hiện

> **Cập nhật sau mỗi task hoàn thành.** Task sau phải đọc file này trước khi bắt đầu.
> **Quy ước:** ✅ Hoàn thành | 🔄 Đang làm | ⏳ Chưa bắt đầu | ❌ Bị block | ⏭️ Skipped

---

## Tổng quan

| Phase | Trạng thái | GATE | Ghi chú |
|---|---|---|---|
| P0 — Scaffolding & Môi trường | ⏳ | GATE-0 | Chưa bắt đầu |
| P1 — App Hello Version | ⏳ | GATE-1 | |
| P2 — Check Update + Manifest | ⏳ | GATE-2 | |
| P3 — Download + Verify + Install + Restart | ⏳ | GATE-3 | |
| P4 — Rollback + Health-Check | ⏳ | GATE-4 | |
| P5 — Test + Docs + Báo cáo | ⏳ | GATE-5 | |

---

## Chi tiết từng Task

| Task | Trạng thái | Ngày hoàn thành | Kết quả kiểm chứng | Ghi chú / Issues |
|---|---|---|---|---|
| **P0-T1** Tạo cấu trúc repo | ✅ | 2026-07-29 13:12 | Layout 15 thư mục khớp §4; 6 file skeleton tạo; .gitignore cover secrets | |
| **P0-T2** Khởi tạo Rust app (Tauri 2) | 🔄 | | | |
| **P0-T3** Cài Go + Khởi tạo Go app (Wails v2) | ⏳ | | | |
| **P0-T4** Tạo GitHub Actions workflows | ⏳ | | | |
| **P0-T5** GATE-0: CI green | ⏳ | | | |
| **P1-T1** Go app hello version | ⏳ | | | |
| **P1-T2** Rust app hello version | ⏳ | | | |
| **P1-T3** Architecture doc + source-map | ⏳ | | | |
| **P1-T4** GATE-1: 2 app chạy được | ⏳ | | | |
| **P2-T1** Sinh key pair Minisign | ⏳ | | | |
| **P2-T2** Go app check update | ⏳ | | | |
| **P2-T3** Rust app check update | ⏳ | | | |
| **P2-T4** Publish manifest v1.0.0 | ⏳ | | | |
| **P2-T5** GATE-2: App phát hiện bản mới | ⏳ | | | |
| **P3-T1** Go app download + verify | ⏳ | | | |
| **P3-T2** Rust app download + verify + install | ⏳ | | | |
| **P3-T3** Luồng update end-to-end (§8 steps 1-14) | ⏳ | | | |
| **P3-T4** GATE-3: Full flow B01-B07 pass | ⏳ | | | |
| **P4-T1** Go app health-check + rollback | ⏳ | | | |
| **P4-T2** Rust app health-check + rollback | ⏳ | | | |
| **P4-T3** Test nhóm E (restart & rollback) | ⏳ | | | |
| **P4-T4** GATE-4: Rollback hoạt động | ⏳ | | | |
| **P5-T1** Test nhóm A (app cơ bản) | ⏳ | | | |
| **P5-T2** Test nhóm C/D/F/G | ⏳ | | | |
| **P5-T3** Hoàn thiện docs (7 file) | ⏳ | | | |
| **P5-T4** Báo cáo cuối cùng (§14) | ⏳ | | | |
| **P5-T5** GATE-5: Definition of Done | ⏳ | | | |

---

## Issues & Lệch hướng

| # | Ngày phát hiện | Mô tả | Hành động | Trạng thái |
|---|---|---|---|---|
| 1 | 2026-07-29 | Wails v3 **không có updater service** (TS-03), prompt giả định có | Chuyển sang Wails v2 + go-selfupdate, ghi rõ giới hạn trong docs | ✅ Đã quyết định |
| 2 | 2026-07-29 | Agent chạy Linux, không thể build/run Windows app trực tiếp | Dùng CI (T2) + user (T3) để kiểm chứng Windows | ✅ Đã thiết kế |
| 3 | 2026-07-29 | Tauri 2 không có built-in rollback | Tự viết adapter, ghi rõ giới hạn | ✅ Đã thiết kế |
