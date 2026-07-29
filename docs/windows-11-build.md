# Build trên Windows 11

> Mặc định các lệnh chạy ở **PowerShell**. Chỗ nào cần Git Bash sẽ ghi rõ.

## Yêu cầu môi trường

| Thành phần | Version | Cài bằng |
|---|---|---|
| Windows | 11 x64 | |
| Go | 1.25.12 | `winget install GoLang.Go` |
| Node | 22 (npm 10) | `winget install OpenJS.NodeJS.LTS` |
| Rust | 1.94.1 + target `x86_64-pc-windows-msvc` | `winget install Rustlang.Rustup` |
| VS Build Tools 2022 | workload *Desktop development with C++* | cần cho linker MSVC |
| WebView2 Runtime | mới nhất | Windows 11 có sẵn |
| Wails CLI | v2.13.0 | `go install github.com/wailsapp/wails/v2/cmd/wails@v2.13.0` |
| gh CLI | mới nhất | chỉ cần khi release bằng tay |

Kiểm tra nhanh:

```powershell
go version; node -v; npm -v; rustc --version; wails doctor
```

## Clone

```powershell
git clone https://github.com/o25160526-pip/self-update-go-rust.git
cd self-update-go-rust
```

## Build bản Go (Wails v2)

```powershell
cd go-demo\frontend
npm ci
npm run build          # BAT BUOC truoc go build: //go:embed all:frontend/dist
cd ..
go test ./... -race
wails build -clean -ldflags "-X go-demo/internal/version.Version=1.0.0"
.\build\bin\go-demo.exe --version
```

Kết quả: `go-demo\build\bin\go-demo.exe` — portable, chạy trực tiếp, không cần cài.

## Build bản Rust (Tauri 2)

```powershell
cd rust-demo
npm ci
cargo fmt --all -- --check
cargo clippy --manifest-path logic\Cargo.toml --all-targets -- -D warnings
cargo test --manifest-path logic\Cargo.toml
npm run tauri build -- --target x86_64-pc-windows-msvc
```

Kết quả: `rust-demo\target\x86_64-pc-windows-msvc\release\bundle\nsis\*-setup.exe` (kèm `.sig` nếu có key ký).

`src-tauri\build.rs` tự sinh `icons\icon.ico` khi thiếu — `tauri-build`, `generate_context!` và NSIS đều bắt buộc có file `.ico`.

## Ký artifact tại máy

```powershell
cd go-demo
go run ./cmd/sign-artifact -key ..\keys\demo-signing.key -password demo-password -in build\bin\go-demo.exe
go run ./cmd/verify-artifact -file build\bin\go-demo.exe -signature build\bin\go-demo.exe.sig
```

Tauri tự ký lúc build nếu set biến môi trường:

```powershell
$env:TAURI_SIGNING_PRIVATE_KEY = (Get-Content ..\keys\demo-signing.key.b64 -Raw).Trim()
$env:TAURI_SIGNING_PRIVATE_KEY_PASSWORD = "demo-password"
```

## Chạy chế độ dev

```powershell
cd go-demo;   wails dev
cd rust-demo; npm run tauri dev
```

## Lỗi hay gặp

Xem [troubleshooting.md](troubleshooting.md).
