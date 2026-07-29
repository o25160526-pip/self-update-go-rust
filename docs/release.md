# Phát hành release

## Quy tắc version

- Mỗi release dùng SemVer dạng `vMAJOR.MINOR.PATCH`.
- Workflow lấy version từ `release/v<version>` hoặc tag `v<version>`.
- Go build inject version bằng `-ldflags`.
- Rust workflow cập nhật `src-tauri/tauri.conf.json` trước khi build.
- Không sửa version thủ công trong source chỉ để phát hành patch.

## Phát hành bằng nhánh release, khuyến nghị

Ví dụ bản tiếp theo là **v1.0.2**:

```powershell
git checkout main
git pull origin main
git checkout -b release/v1.0.2
git push origin release/v1.0.2
```

Push nhánh sẽ chạy đồng thời `release-go.yml` và `release-rust.yml`. Mỗi workflow:

1. Lấy `1.0.2` từ tên nhánh.
2. Build đúng version trên Windows.
3. Ký bằng GitHub Secret nếu đã cấu hình, nếu chưa thì fallback demo key trong `keys/`.
4. Verify signature trước khi upload.
5. Tạo hoặc mở Release `v1.0.2`.
6. Upload asset riêng của Go hoặc Rust.

## Checklist bắt buộc sau CI

```powershell
# Chạy trên PowerShell, thay URL theo release thật
Invoke-WebRequest https://github.com/o25160526-pip/self-update-go-rust/releases/download/v1.0.2/manifest-go.json
Invoke-WebRequest https://github.com/o25160526-pip/self-update-go-rust/releases/download/v1.0.2/latest-rust.json
```

Release phải có đủ:

| Nhóm | Asset bắt buộc |
|---|---|
| Go | `go-demo-windows-x64.exe`, `.sig`, `.sha256`, `manifest-go.json` |
| Rust | `rust-demo-windows-x64-setup.exe`, `.sig`, `.sha256`, `latest-rust.json` |

## Kiểm thử update thật trên Windows 11

1. Tải Go `v1.0.0` hoặc `v1.0.1` từ [Releases](https://github.com/o25160526-pip/self-update-go-rust/releases).
2. Chạy bản cũ từ thư mục user có quyền ghi, không đặt trong `C:\Program Files`.
3. Để app tự kiểm tra hoặc chạy Go với `--check-only`.
4. Xác nhận log đi qua `update-available → downloading → verifying → installing → restarting`.
5. Sau restart xác nhận version `1.0.2` và state `up-to-date`.
6. Test Rust bằng installer cũ, xác nhận NSIS/UAC, rồi xác nhận bản mới mở được.
7. Ghi kết quả T3 vào ClickUp History Work.

## Ký bằng key production

```powershell
cd rust-demo
npx tauri signer generate -w "$env:USERPROFILE\.tauri\myapp.key"
[Convert]::ToBase64String([IO.File]::ReadAllBytes("$env:USERPROFILE\.tauri\myapp.key"))
```

Thêm GitHub Actions Secrets:

- `TAURI_SIGNING_PRIVATE_KEY`: nội dung base64 của private key.
- `TAURI_SIGNING_PRIVATE_KEY_PASSWORD`: password private key.

Client vẫn phải được build với **public key mới** ở cả Go, Rust logic và `tauri.conf.json`. Không được chỉ đổi secret, vì client cũ sẽ không verify được chữ ký bằng public key cũ.

## Phát hành tag trực tiếp

```powershell
git tag v1.0.2
git push origin v1.0.2
```

Dùng cách này khi credential có quyền push tag. Với integration hiện tại, dùng nhánh `release/v1.0.2` an toàn hơn.

## Rollback release

Không xoá release đã phát hành để sửa lỗi. Tăng patch version, ví dụ `v1.0.3`, rồi phát hành lại. Client sẽ bỏ qua bản thấp hơn và lấy bản mới hơn.
