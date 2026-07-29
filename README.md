# self-update-go-rust

Demo **auto-update qua GitHub Releases** cho app desktop Windows, làm song song 2 bản:

| | Go | Rust |
|---|---|---|
| Framework | Wails v2.13.0 | Tauri 2.11.5 |
| Updater | `updater/` tự viết (check + verify + staged install + rollback) | `tauri-plugin-updater` 2.10.1 |
| Artifact | `go-demo-windows-x64.exe` (portable) | `rust-demo-windows-x64-setup.exe` (NSIS) |
| Manifest | `manifest-go.json` | `latest-rust.json` |
| Chữ ký | Minisign / Ed25519 (key ID `5756DB5A3509A8C1`) | cùng key |

## Chạy thử trên Windows 11 (không cần build)

1. Vào [Releases](https://github.com/o25160526-pip/self-update-go-rust/releases) và tải:
   - **Go:** `go-demo-windows-x64.exe` — chạy trực tiếp, không cần cài.
   - **Rust:** `rust-demo-windows-x64-setup.exe` — cài đặt qua NSIS.
2. Mở app → thấy `Hello, version X.Y.Z` + trạng thái updater + log.
3. Khi có release mới hơn, app **tự** phát hiện → tải → verify SHA-256 và chữ ký → cài → khởi động lại với version mới. Không cần bấm gì.

SmartScreen sẽ cảnh báo vì app chưa có code-signing certificate → *More info* → *Run anyway*. Xem [docs/troubleshooting.md](docs/troubleshooting.md).

### Tự kiểm tra file tải về (tuỳ chọn)

```powershell
(Get-FileHash .\go-demo-windows-x64.exe -Algorithm SHA256).Hash.ToLower()
Get-Content .\go-demo-windows-x64.exe.sha256

cd go-demo
go run ./cmd/verify-artifact -manifest ..\manifest-go.json -file ..\go-demo-windows-x64.exe
```

## CLI flags

| Flag | Tác dụng |
|---|---|
| `--version` | In version rồi thoát |
| `--print-update-state` | In trạng thái updater đã lưu |
| `--offline-test` | Dùng mock provider, không gọi mạng |
| `--check-only` | (Go) chạy trọn luồng update ở chế độ CLI rồi thoát |

## Build từ source

Chi tiết: [docs/windows-11-build.md](docs/windows-11-build.md). Tóm tắt:

```powershell
# Go
cd go-demo\frontend; npm ci; npm run build; cd ..
wails build -clean -ldflags "-X go-demo/internal/version.Version=1.0.0"

# Rust
cd rust-demo; npm ci; npm run tauri build
```

## Phát hành

CI tự build + ký + tạo release khi push nhánh `release/v<version>` hoặc tag `v<version>`. Xem [docs/release.md](docs/release.md).

## Tài liệu

| File | Nội dung |
|---|---|
| [docs/architecture.md](docs/architecture.md) | Framework, mapping manifest, khác biệt Go ↔ Rust |
| [docs/source-map.md](docs/source-map.md) | Bản đồ source Go ↔ Rust |
| [docs/update-flow.md](docs/update-flow.md) | Luồng update + state machine |
| [docs/rollback.md](docs/rollback.md) | Health-check + rollback |
| [docs/release.md](docs/release.md) | Cách phát hành |
| [docs/windows-11-build.md](docs/windows-11-build.md) | Build trên Windows 11 |
| [docs/testing.md](docs/testing.md) | 7 nhóm test A-G |
| [docs/troubleshooting.md](docs/troubleshooting.md) | Lỗi thường gặp |
| [keys/README.md](keys/README.md) | Key ký + fallback secret |
| [plan/PLAN.md](plan/PLAN.md) · [plan/PROGRESS.md](plan/PROGRESS.md) | Kế hoạch và tiến độ |

## Bảo mật (đọc trước khi dùng thật)

- Private key trong `keys/` là **key demo công khai**, chỉ để test. Production: sinh key mới rồi đặt vào GitHub Secrets `TAURI_SIGNING_PRIVATE_KEY` — CI tự ưu tiên secret nếu có.
- Chưa có Authenticode code-signing nên Windows SmartScreen vẫn cảnh báo.
- Client không cần GitHub token vì release là public.
