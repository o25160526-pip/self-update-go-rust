# PROMPT AGENT: Demo auto-update qua GitHub Releases bằng Go và Rust

# PROMPT AGENT: Demo auto-update qua GitHub Releases bằng Go và Rust
## 1\. Vai trò
Bạn là một senior desktop engineer kiêm release engineer. Hãy tạo **hai app demo độc lập nhưng có cấu trúc source tương đương** để kiểm chứng end-to-end việc phát hành và tự cập nhật ứng dụng desktop qua GitHub Releases:
*   Bản Go.
*   Bản Rust.

Demo chỉ cần hiển thị giao diện tối giản: `Hello, version X.Y.Z`. Không xây nghiệp vụ PostgreSQL hay API thật trong demo, nhưng phải giữ kiến trúc có thể mở rộng sau này cho local service, API bên ngoài và PostgreSQL 11.

Không được chỉ viết hướng dẫn lý thuyết. Phải tạo source code chạy được, script build/release, cấu hình updater, tài liệu Windows 11 và kịch bản kiểm thử có kết quả quan sát được.
## 2\. Mục tiêu bắt buộc
Khi hoàn thành, người dùng phải có thể:

1. Clone repository trên Windows 11.
2. Build bản Go hoặc Rust ở local.
3. Chạy app và thấy `Hello, version 1.0.0`.
4. Tạo tag release mới, ví dụ `v1.0.1`.
5. Chạy một lệnh hoặc GitHub Actions để build artifact và publish GitHub Release.
6. Mở app cũ `1.0.0`.
7. App tự kiểm tra GitHub Releases.
8. App hiển thị thông báo có bản cập nhật.
9. App tự động tải, xác minh chữ ký, cài đặt và restart mà **không yêu cầu người dùng xác nhận**.
10. Sau khi restart, app hiển thị `Hello, version 1.0.1`.
11. Có kịch bản mô phỏng artifact hỏng, chữ ký sai, app crash sau update và rollback về bản ổn định trước đó.

Mặc định chọn:
*   Kênh phát hành: GitHub Releases.
*   Kênh update: `stable`.
*   Manifest update: static manifest hoặc format chính thức của framework, đặt trong GitHub Release hoặc GitHub Pages/URL tĩnh.
*   Update policy: tự động kiểm tra khi app khởi động và định kỳ.
*   Update mode: download và install tự động, không hỏi xác nhận.
*   Security: HTTPS, SHA-256, chữ ký artifact/manifest, public key pin trong app.
*   Rollback: giữ bản `last-known-good`, chỉ đánh dấu bản mới là active sau health-check thành công.
*   Target build chính: Windows 11 x64. Thiết kế source và CI phải sẵn sàng mở rộng macOS/Linux.
## 3\. Nguyên tắc an toàn
Không được triển khai kiểu app tự ghi đè chính executable đang chạy nếu framework có helper/updater core chính thức.

Ưu tiên dùng core có sẵn:
*   Go: **Wails v3** **`app.Updater`** nếu dùng Wails; nếu demo dùng Go binary thuần thì dùng một thư viện updater ổn định như `creativeprojects/go-selfupdate`, nhưng phải tự bổ sung helper, staged install, restart và rollback rõ ràng.
*   Rust: **Tauri 2 updater plugin** nếu dùng Tauri; không tự viết lại download, verify và install nếu plugin đã hỗ trợ.

Nếu version framework thực tế không hỗ trợ đúng một chức năng, không giả vờ rằng có. Ghi rõ giới hạn, chọn phương án gần nhất và cung cấp adapter nhỏ có test.

Không hard-code private signing key trong source. Private key chỉ nằm ở GitHub Actions Secrets hoặc file local bị gitignore để test.

Không nhúng GitHub Personal Access Token vào app client. Nếu GitHub Release public thì client chỉ dùng URL public. Nếu cần private release, phải thêm update gateway hoặc URL token ngắn hạn, nhưng không dùng trong demo mặc định.
## 4\. Hai app phải có layout source tương đương
Tạo một repository có layout dễ so sánh:

```text
 desktop-auto-update-demo/
 ├─ README.md
 ├─ docs/
 │  ├─ architecture.md
 │  ├─ windows-11-build.md
 │  ├─ release.md
 │  ├─ update-flow.md
 │  ├─ rollback.md
 │  ├─ testing.md
 │  └─ troubleshooting.md
 ├─ shared/
 │  ├─ update-policy.example.json
 │  ├─ release-manifest.example.json
 │  └─ test-cases.md
 ├─ go-demo/
 │  ├─ app/
 │  ├─ frontend/
 │  ├─ updater/
 │  ├─ installer/
 │  ├─ scripts/
 │  ├─ build/
 │  ├─ go.mod
 │  └─ README.md
 ├─ rust-demo/
 │  ├─ src/
 │  ├─ frontend/
 │  ├─ updater/
 │  ├─ installer/
 │  ├─ scripts/
 │  ├─ build/
 │  ├─ Cargo.toml
 │  └─ README.md
 ├─ .github/
 │  └─ workflows/
 │     ├─ release-go.yml
 │     └─ release-rust.yml
 └─ .gitignore
```

Nếu Wails/Tauri tạo layout khác, giữ ý nghĩa tương đương và bổ sung `docs/source-map.md` để chỉ ra mapping:

| Ý nghĩa | Go | Rust |
| ---| ---| --- |
| Version nguồn duy nhất | `internal/version` hoặc tương đương | `src/version` hoặc tương đương |
| UI hello version | `frontend` | `frontend` |
| Check update | `updater/check` | `updater/check` |
| Download/verify/install | `updater/install` | `updater/install` |
| Restart | `updater/restart` | `updater/restart` |
| Health-check | `updater/health` | `updater/health` |
| Rollback | `updater/rollback` | `updater/rollback` |
| Build script | `scripts/build` | `scripts/build` |
| Release workflow | `release-go.yml` | `release-rust.yml` |

## 5\. Chọn framework và ghi rõ phiên bản
Tại thời điểm bắt đầu, kiểm tra tài liệu chính thức và khóa version:
*   Go framework: ưu tiên Wails v3 nếu updater đã đủ khả năng cho demo.
*   Rust framework: ưu tiên Tauri 2 và `tauri-plugin-updater`.
*   Frontend: dùng HTML/CSS/TypeScript tối giản hoặc frontend có sẵn nhưng không kéo thêm dependency không cần thiết.
*   Go/Rust toolchain: ghi rõ version trong `go.mod`, `go.work`, `Cargo.toml`, `rust-toolchain.toml` hoặc file tương đương.
*   Node package manager: khóa bằng lockfile.
*   GitHub Actions runner: dùng Windows runner cho artifact Windows; có thể thêm matrix sau nhưng không làm phức tạp demo.

Trước khi code, tạo bảng trong `docs/architecture.md`:

| Thành phần | Go | Rust | Lý do chọn |
| ---| ---| ---| --- |
| Desktop framework | ... | ... | ... |
| Updater core | ... | ... | ... |
| Manifest format | ... | ... | ... |
| Signature | ... | ... | ... |
| Installer | ... | ... | ... |
| Restart mechanism | ... | ... | ... |
| Rollback mechanism | ... | ... | ... |

## 6\. Tính năng ứng dụng tối thiểu
Mỗi app phải:
*   Hiển thị tên app, version hiện tại, OS và architecture.
*   Hiển thị trạng thái updater: `checking`, `up-to-date`, `update-available`, `downloading`, `verifying`, `installing`, `restarting`, `failed`, `rolled-back`.
*   Khi có version mới, hiển thị thông báo rõ ràng.
*   Theo policy mặc định, tự động download/install mà không cần nút xác nhận.
*   Có một vùng log dễ nhìn hoặc nút mở log.
*   Hiển thị release notes ngắn từ manifest nếu có.
*   Có chế độ `--offline-test` hoặc mock provider để test mà không cần mạng.
*   Có `--version` và `--print-update-state` để kiểm tra từ terminal.
*   Không lưu secret, token hay private key vào log.

UI không cần nút “Update now” trong flow mặc định. Có thể thêm nút debug để mô phỏng check/update, nhưng không được làm thay đổi hành vi tự động mặc định.
## 7\. Metadata và policy tự động cập nhật
Tạo file mẫu, có schema và giải thích từng trường:

```json
{
  "channel": "stable",
  "autoCheckOnStartup": true,
  "autoCheckIntervalMinutes": 60,
  "autoDownload": true,
  "autoInstall": true,
  "requireUserConfirmation": false,
  "restartAutomatically": true,
  "allowDowngrade": false,
  "rollbackOnStartupFailure": true,
  "healthCheckTimeoutSeconds": 30,
  "maxRollbackAttempts": 1,
  "minSupportedVersion": "1.0.0"
}
```

Nếu framework có manifest riêng, map các trường trên sang format của framework và ghi rõ trường nào do client policy xử lý, trường nào do server manifest xử lý.

Manifest phải hỗ trợ tối thiểu:

```json
{
  "version": "1.0.1",
  "channel": "stable",
  "publishedAt": "2026-07-29T00:00:00Z",
  "releaseNotes": "Fix update flow and restart",
  "minSupportedVersion": "1.0.0",
  "platform": "windows-x86_64",
  "url": "https://github.com/OWNER/REPO/releases/download/v1.0.1/app-windows-x64.exe",
  "sha256": "...",
  "signature": "...",
  "size": 12345678,
  "mandatory": false
}
```

Không để client tin `url`, version hoặc signature nếu chưa verify manifest/artifact theo cơ chế ký của framework.
## 8\. Luồng update bắt buộc
Implement và mô tả chính xác flow sau:

1. App đọc version build từ một nguồn duy nhất.
2. App đọc policy local.
3. App gọi GitHub Releases hoặc manifest URL qua HTTPS.
4. App lọc đúng channel, OS và CPU.
5. App so sánh semver, chống downgrade.
6. Nếu không có bản mới, hiển thị `Up to date`.
7. Nếu có bản mới, hiển thị thông báo.
8. Tự động tải artifact vào thư mục tạm riêng.
9. Kiểm tra kích thước và SHA-256.
10. Xác minh chữ ký bằng public key pin trong app/framework config.
11. Ghi state `pending`.
12. Flush dữ liệu, đóng local service nếu có, đóng app chính.
13. Helper/launcher/installer thực hiện swap hoặc install.
14. Restart app với metadata `updated-from`.
15. App mới thực hiện startup health-check.
16. Chỉ ghi `last-known-good` sau khi health-check thành công.
17. Nếu timeout/crash, helper hoặc launcher rollback về bản trước.
18. Hiển thị kết quả update sau restart.

Tuyệt đối không đánh dấu update thành công chỉ vì download xong.
## 9\. Release pipeline từ source đến GitHub Release
Viết `docs/release.md` để người mới làm theo từng bước:
### Local release
Cung cấp lệnh cụ thể, chạy được trên PowerShell Windows 11:

```powershell
. scripts\version.ps1 1.0.0
. scripts\build.ps1 -App go -Version 1.0.0
. scripts\build.ps1 -App rust -Version 1.0.0
. scripts\package.ps1 -App go -Version 1.0.0
. scripts\package.ps1 -App rust -Version 1.0.0
```

Không bắt buộc giữ đúng tên lệnh nếu framework cần khác, nhưng phải cung cấp lệnh thật và nhất quán giữa hai bản.
### Git tag và GitHub Actions
Mô tả:

```powershell
git checkout main
git pull
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions phải:
*   Trigger khi push tag `v*`.
*   Checkout source.
*   Cài toolchain đã khóa version.
*   Build release mode.
*   Tạo artifact Windows 11 x64.
*   Tạo checksum SHA-256.
*   Ký artifact hoặc tạo metadata signature bằng secret của GitHub Actions.
*   Tạo/update manifest.
*   Publish GitHub Release.
*   Upload artifact, checksum, signature và manifest.
*   Không in private key hoặc secret ra log.
*   Fail rõ ràng nếu thiếu secret ký.
*   Có bước kiểm tra artifact sau upload.

Tạo hai workflow tương đương:
*   `.github/workflows/release-go.yml`
*   `.github/workflows/release-rust.yml`

Nếu demo dùng cùng repository và cùng tag, cho phép chọn workflow bằng `workflow_dispatch` hoặc dùng tên artifact khác nhau để không ghi đè.
## 10\. Chữ ký và secret
Chọn cơ chế chữ ký chính thức của framework nếu có.

Tài liệu phải mô tả:
*   Cách tạo key pair trên máy developer chỉ để test.
*   Cách lưu private key trong GitHub Actions Secrets.
*   Cách cấu hình public key trong Go app và Rust app.
*   Cách rotate key trong tương lai.
*   Cách phát hiện signature sai.
*   Cách xử lý khi mất private key.
*   Vì sao SHA-256 không thay thế chữ ký.

Không commit:
*   Private key.
*   GitHub token.
*   Password PostgreSQL.
*   Token API.
*   File release thật có secret.
## 11\. Windows 11 build docs bắt buộc
Tạo `docs/windows-11-build.md` với hướng dẫn copy/paste được:
*   Windows 11 x64.
*   PowerShell 5.1 và PowerShell 7 nếu khác nhau.
*   Git.
*   Go version.
*   Rustup, stable toolchain và Visual Studio Build Tools.
*   C++ workload cần thiết.
*   WebView2 nếu framework yêu cầu.
*   Node.js và package manager nếu frontend cần.
*   GitHub CLI tùy chọn.
*   Kiểm tra PATH.
*   Clone repository.
*   Install dependencies.
*   Build debug.
*   Build release.
*   Chạy app.
*   Xem version.
*   Xem log.
*   Gỡ app và xóa state update.
*   Các lỗi thường gặp: linker, WebView2, quyền thư mục, file đang bị khóa, UAC, antivirus quarantine, signature mismatch.

Mỗi lệnh phải ghi rõ đang chạy ở PowerShell hay Git Bash. Ưu tiên PowerShell.
## 12\. Kịch bản kiểm thử
Tạo `docs/testing.md` và test runner nếu có thể. Mỗi test phải có: mã test, setup, command, expected result, log cần kiểm tra, cleanup.
### Nhóm A: app cơ bản
*   A01: clean build Go.
*   A02: clean build Rust.
*   A03: app hiển thị đúng version.
*   A04: `--version` trả đúng version.
*   A05: app khởi động offline.
*   A06: app không ghi secret vào log.
### Nhóm B: update bình thường
*   B01: chạy `1.0.0`, không có release mới, app báo up-to-date.
*   B02: chạy `1.0.0`, có `1.0.1`, app báo update available.
*   B03: auto download không cần click.
*   B04: hash đúng, signature đúng, install thành công.
*   B05: app tự restart.
*   B06: sau restart hiển thị `1.0.1`.
*   B07: state ghi `last-known-good=1.0.1`.
### Nhóm C: lỗi mạng
*   C01: DNS fail.
*   C02: timeout.
*   C03: HTTP 404.
*   C04: HTTP 500.
*   C05: download bị ngắt giữa chừng.
*   C06: resume hoặc retry đúng policy.
*   C07: app vẫn chạy version cũ, không bị hỏng.
### Nhóm D: bảo mật
*   D01: SHA-256 sai, update bị từ chối.
*   D02: signature sai, update bị từ chối.
*   D03: manifest bị sửa, update bị từ chối.
*   D04: version thấp hơn hiện tại, bị từ chối.
*   D05: artifact không đúng OS/architecture, bị từ chối.
*   D06: private key không xuất hiện trong artifact/log.
*   D07: client không cần GitHub token với public release.
### Nhóm E: restart và rollback
*   E01: app đang giữ file, helper chờ app thoát.
*   E02: app crash trước khi health-check.
*   E03: app mới treo quá timeout.
*   E04: local service không khởi động.
*   E05: rollback về `last-known-good`.
*   E06: không rollback lặp vô hạn.
*   E07: máy tắt giữa lúc install, khởi động lại vẫn phục hồi được.
*   E08: rollback không xóa nhầm user data/config.
### Nhóm F: quyền Windows
*   F01: cài per-user, update không cần admin.
*   F02: cài machine-wide, UAC chỉ xuất hiện ở bước cần thiết.
*   F03: user không có quyền ghi thư mục, lỗi rõ ràng.
*   F04: artifact được tải vào thư mục temp có permission phù hợp.
*   F05: code signing/SmartScreen behavior được ghi nhận.
### Nhóm G: phát hành
*   G01: tag `v1.0.0` tạo release.
*   G02: tag `v1.0.1` tạo artifact mới.
*   G03: artifact Go và Rust không ghi đè nhau.
*   G04: checksum upload đúng.
*   G05: manifest trỏ đúng artifact.
*   G06: workflow fail khi thiếu signing secret.
*   G07: release có release notes.
## 13\. Definition of Done
Chỉ coi hoàn thành khi:
*   Hai app Go và Rust đều build được trên Windows 11 x64.
*   Source layout và naming tương đương, có source-map.
*   App hello version chạy được.
*   Có GitHub Actions cho release.
*   Có artifact và manifest/signature rõ ràng.
*   Client cũ tự phát hiện release mới.
*   Client tự download, verify, install và restart không cần xác nhận.
*   Có state machine và log dễ kiểm tra.
*   Có last-known-good và rollback khi startup health-check thất bại.
*   Có test cho hash sai, signature sai, network fail, crash và quyền Windows.
*   Có docs từ clone source đến GitHub Release và client update thành công.
*   Không có secret trong source, artifact hoặc log.
*   Các giới hạn framework được ghi rõ, không dùng claim marketing thay cho test thực tế.
## 14\. Báo cáo cuối cùng agent phải xuất
Cuối nhiệm vụ, xuất báo cáo ngắn gồm:

1. Cây thư mục thực tế.
2. Framework và version đã khóa.
3. Các lệnh build local cho Go và Rust.
4. Các lệnh tạo tag và release.
5. Tên artifact đã upload.
6. URL GitHub Release dùng để test.
7. Cách client phát hiện update.
8. Cách restart và rollback hoạt động.
9. Kết quả từng nhóm test A-G.
10. Các điểm chưa thể tự động hóa và lý do.
11. Các việc cần làm trước production: code signing/notarization, private update server, staged rollout, telemetry, key rotation và backup.

Đừng trả lời bằng pseudo-code nếu có thể viết code thật. Nếu thiếu credential GitHub, vẫn phải hoàn thành source, workflow, scripts và local test; ghi chính xác bước người dùng cần thực hiện để publish release sau đó.