# desktop-auto-update-demo

Demo end-to-end: phát hành và tự cập nhật ứng dụng desktop qua GitHub Releases bằng **Go (Wails v2)** và **Rust (Tauri 2)**.

## Mục tiêu

Xem [`PROMPT AGENT_ Demo auto-update qua GitHub Releases bằng Go và Rust-20260729031923.md`](../PROMPT%20AGENT_%20Demo%20auto-update%20qua%20GitHub%20Releases%20bằng%20Go%20và%20Rust-20260729031923.md) và [`plan/PLAN.md`](plan/PLAN.md).

## Trạng thái

| Phase | Trạng thái |
|---|---|
| P0 — Scaffolding | 🔄 |
| P1 — Hello Version | ⏳ |
| P2 — Check Update | ⏳ |
| P3 — Download + Install | ⏳ |
| P4 — Rollback | ⏳ |
| P5 — Test + Docs | ⏳ |

## Layout

```
desktop-auto-update-demo/
├─ go-demo/          # Go + Wails v2
├─ rust-demo/        # Rust + Tauri 2
├─ shared/           # manifest + policy mẫu
├─ docs/             # tài liệu
├─ .github/workflows/
└─ plan/             # PLAN + TRUST-SOURCES + PROGRESS
```

## Quick start

Xem [`docs/windows-11-build.md`](docs/windows-11-build.md) (sẽ hoàn thiện ở P5).
