# Phát hành release

## Cách 1 — push nhánh `release/v<version>` (khuyến nghị)

```powershell
git checkout main; git pull
git checkout -b release/v1.0.1
git push origin release/v1.0.1
```

Cả 2 workflow sẽ:

1. Lấy version từ tên nhánh (`release/v1.0.1` → `1.0.1`).
2. Build Windows: Go dùng `-ldflags -X ...Version=1.0.1`; Rust ghi `1.0.1` vào `tauri.conf.json` trước khi build.
3. Ký artifact (secret nếu có, không thì demo key trong `keys/`).
4. `gh release create v1.0.1` (tạo luôn tag) rồi upload asset. Hai workflow chạy song song nên bước tạo release có retry: ai tới trước thì tạo, người sau chỉ upload thêm.

> Dùng nhánh vì GitHub App token không push được tag. Nếu bạn tự push tag thì workflow cũng nhận và chạy y như vậy.

## Cách 2 — push tag

```powershell
git tag v1.0.1
git push origin v1.0.1
```

## Asset của mỗi release

| File | Workflow |
|---|---|
| `go-demo-windows-x64.exe` + `.sig` + `.sha256` | release-go |
| `manifest-go.json` | release-go |
| `rust-demo-windows-x64-setup.exe` + `.sig` + `.sha256` | release-rust |
| `latest-rust.json` | release-rust |

Hai bản **không ghi đè nhau** vì tên file khác nhau.

## Client tìm bản mới thế nào

- Go: `https://github.com/<repo>/releases/latest/download/manifest-go.json`
- Rust: `https://github.com/<repo>/releases/latest/download/latest-rust.json`

`releases/latest` luôn trỏ release non-draft, non-prerelease mới nhất → chỉ cần publish release mới là client cũ thấy ngay, không phải đổi code.

## Ký bằng key thật

```powershell
cd rust-demo
npx tauri signer generate -w "$env:USERPROFILE\.tauri\myapp.key"
[Convert]::ToBase64String([IO.File]::ReadAllBytes("$env:USERPROFILE\.tauri\myapp.key"))
```

Thêm 2 secret: `TAURI_SIGNING_PRIVATE_KEY` (chuỗi base64 ở trên) và `TAURI_SIGNING_PRIVATE_KEY_PASSWORD`. Rồi thay public key mới vào cả 3 chỗ pin (xem [keys/README.md](../keys/README.md)). Không có secret thì CI tự dùng demo key nên build vẫn xanh.

## Release local (không qua CI)

```powershell
cd go-demo
wails build -clean -ldflags "-X go-demo/internal/version.Version=1.0.1"
Copy-Item build\bin\*.exe build\bin\go-demo-windows-x64.exe -Force
go run ./cmd/sign-artifact -key ..\keys\demo-signing.key -password demo-password -in build\bin\go-demo-windows-x64.exe
.\scripts\generate-manifest.ps1 -Version 1.0.1 -Artifact build\bin\go-demo-windows-x64.exe -Signature build\bin\go-demo-windows-x64.exe.sig -Repository o25160526-pip/self-update-go-rust
gh release create v1.0.1 --title "Release v1.0.1" --notes "Manual release"
gh release upload v1.0.1 build\bin\go-demo-windows-x64.exe build\bin\go-demo-windows-x64.exe.sig manifest-go.json
```
