# PROGRESS.md — Theo dõi tiến độ thực hiện

> Cập nhật sau mỗi task hoàn thành.
> **Mirror trong ClickUp:** doc "History Work — self-update-go-rust" (trang 📊 PROGRESS + 🕒 Work Log + 📥 Yêu cầu từ user).
> Quy ước: ✅ Hoàn thành | 🔄 Đang làm | ⏳ Chưa bắt đầu | ❌ Bị block | ⏭️ Skipped

## Tổng quan

| Phase | Trạng thái | GATE | Ghi chú |
|---|---|---|---|
| P0 — Scaffolding & Môi trường | ✅ | GATE-0 | Go + Rust build Windows xanh |
| P1 — App Hello Version | 🔄 | GATE-1 | Rust v1.0.2 blank UI được sửa, phát hành v1.0.3 |
| P2 — Check Update + Manifest | ✅ | GATE-2 | v1.0.0/v1.0.1/v1.0.2 đã có release + manifest |
| P3 — Download + Verify + Install + Restart | 🔄 | GATE-3 | Go T3 pass; Rust T3 đang xác nhận |
| P4 — Rollback + Health-Check | 🔄 | GATE-4 | T1 pass; T3 cần test Windows 11 |
| P5 — Test + Docs + Báo cáo | 🔄 | GATE-5 | Docs hoàn chỉnh; chờ T3 |

## Releases

| Release | Trạng thái | Ghi chú |
|---|---|---|
| v1.0.0 | ✅ Published | Baseline |
| v1.0.1 | ✅ Published | Go update đã xác nhận |
| v1.0.2 | ✅ Published | Go update pass; Rust UI lỗi do thiếu `withGlobalTauri` |
| v1.0.3 | 🔄 | Bản sửa Rust UI đang build |

## Phát hiện và xử lý Rust UI

**Triệu chứng:** Rust app hiển thị `Hello, version loading...`, OS `...`, state `checking`, log chỉ có `Updater log`.

**Nguyên nhân:** frontend gọi `window.__TAURI__.core.invoke`, nhưng `tauri.conf.json` chưa bật `app.withGlobalTauri`. Vì vậy bridge không tồn tại, JavaScript dừng trước khi lấy version/state.

**Đã sửa:** bật `withGlobalTauri: true`; frontend kiểm tra bridge rõ ràng, hiển thị lỗi startup thay vì loading vô hạn; bỏ dùng `window.__TAURI__.os` sai cách và hiển thị Windows/amd64 ổn định.

## Còn lại

1. CI v1.0.3 xanh cả Go/Rust.
2. Tải Rust v1.0.3 trên Windows 11, xác nhận version/state/UI.
3. Test Rust update từ v1.0.2 → v1.0.3 và rollback.
4. Cập nhật kết quả T3 vào ClickUp History Work, sau đó chốt GATE-5.
