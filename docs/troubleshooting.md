# Troubleshooting

## Chạy app

**SmartScreen: "Windows protected your PC"**
App chưa có Authenticode code-signing certificate. Minisign chỉ để updater verify, Windows không hiểu nó. Bấm *More info* → *Run anyway*. Muốn hết hẳn phải mua EV code-signing cert và ký bằng `signtool`.

**App mở lên rồi tắt ngay**
Chạy từ PowerShell để thấy log: `.\go-demo.exe --check-only`. Thường do thiếu WebView2 Runtime → tải WebView2 Evergreen Bootstrapper của Microsoft.

**UI trắng trơn**
Frontend chưa build trước khi `wails build`, nên `//go:embed all:frontend/dist` chỉ nhúng file `.gitkeep`. Chạy `npm ci; npm run build` trong `go-demo\frontend` rồi build lại.

**Antivirus xoá file `.exe` mới tải về**
Thêm exclusion cho `%LOCALAPPDATA%\go-demo\update-tmp`. Đây cũng là lý do luồng update verify hash + chữ ký trước khi thay exe.

## Update

**State đứng ở `checking` hoặc `failed`**
Xem vùng log trong app. Nguyên nhân thường gặp:
- Chưa có release nào trên GitHub → manifest 404.
- Không có mạng, hoặc proxy chặn `github.com` và `objects.githubusercontent.com`.
- `manifest-go.json` thiếu field, hoặc URL không phải HTTPS → validate từ chối.

**`invalid Minisign signature`**
Artifact được ký bằng key khác với key pin trong client. Kiểm tra 3 chỗ pin có cùng key ID (xem `keys/README.md`), và xem log bước *Resolve signing key* trong CI đang dùng secret hay demo key.

**`artifact SHA-256 mismatch`**
Manifest và artifact lệch nhau, thường vì upload lại artifact mà quên regenerate manifest. Chạy lại release.

**`doi ten exe hien tai: Access is denied`**
Có tiến trình khác đang giữ file (2 instance cùng chạy), hoặc app nằm trong `C:\Program Files` mà không có quyền ghi. Bản Go là portable — nên để ở thư mục user (Desktop, Downloads, `%LOCALAPPDATA%`).

**Update xong nhưng version vẫn cũ**
Có thể đã rollback. Xem `%LOCALAPPDATA%\go-demo\state.json`: `status: rolled-back`, `rollbackAttempts: 1`. Bản lỗi được giữ ở `*.failed` để debug.

**Reset sạch trạng thái update**
```powershell
Remove-Item -Recurse -Force "$env:LOCALAPPDATA\go-demo"
```

## Build

**`pattern all:frontend/dist: no matching files found`**
Chưa build frontend — xem mục UI trắng trơn ở trên.

**Rust: `link.exe not found` / `error: linker`**
Thiếu Visual Studio Build Tools 2022 workload *Desktop development with C++*.

**Tauri: không tìm thấy `icon.ico`**
`build.rs` tự sinh khi thiếu. Nếu vẫn lỗi, xoá `rust-demo\src-tauri\icons\icon.ico` rồi build lại.

**`npm ci` báo lệch lock file**
Dùng đúng Node 22 / npm 10 như CI.

## CI

**Step `Test` chạy rất lâu**
Demo key cố tình dùng scrypt tham số nhẹ (N=1024, r=8, p=1). Tham số "sensitive" (N=65536, p=16) làm `go test -race` chậm hàng phút khi giải mã key. Nếu thay key thật, nhớ điều này.

**`khong tao duoc release vX.Y.Z`**
Hai workflow tạo release song song; bước publish có retry 5 lần. Nếu vẫn fail, kiểm tra quyền `contents: write`.
