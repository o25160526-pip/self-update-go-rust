# TRUST-SOURCES.md — Nguồn tin cậy & Version đã khóa

> **Vai trò của file này:** mọi quyết định kỹ thuật trong `PLAN.md` phải trace được về một dòng trong file này.
> Agent **không được** tự suy đoán version hay API. Nếu cần thông tin mới → verify bằng lệnh trong cột "Verify command", rồi cập nhật file này kèm ngày.
>
> **Lần verify gần nhất:** 2026-07-29 (Asia/Bangkok) — thực hiện bằng GitHub API / crates.io API / go.dev API.
> **Hạn dùng:** nếu lần verify > 30 ngày, phải chạy lại `scripts/verify-versions.sh` trước khi code tiếp.

---

## A. Kết quả verify (2026-07-29)

| # | Thành phần | Version thật đã verify | Nguồn (trust source) | Verify command |
|---|---|---|---|---|
| TS-01 | Wails **v2** (stable) | `v2.13.0` — published `2026-07-06T10:44:59Z` | GitHub Releases API | `curl -s https://api.github.com/repos/wailsapp/wails/releases/latest \| grep tag_name` |
| TS-02 | Wails **v3** (alpha) | `v3.0.0-alpha.102` — **Alpha, chưa stable** | GitHub Tags API + README ("v3 \| Alpha") | `curl -s "https://api.github.com/repos/wailsapp/wails/tags?per_page=100" \| grep v3.0.0-alpha` |
| TS-03 | Wails v3 services có sẵn | `dock, fileserver, kvstore, log, notifications, sqlite` — **KHÔNG có updater** | GitHub Contents API tại ref `v3.0.0-alpha.102` | `curl -s "https://api.github.com/repos/wailsapp/wails/contents/v3/pkg/services?ref=v3.0.0-alpha.102" \| grep '"name"'` |
| TS-04 | Tauri (crate) | `2.11.5` (max_stable); GitHub tag `tauri-v2.11.5` @ `2026-07-01T13:56:43Z` | crates.io API + GitHub Releases API | `curl -s -H "User-Agent: research" https://crates.io/api/v1/crates/tauri` |
| TS-05 | `tauri-plugin-updater` | `2.10.1` — updated `2026-04-04T16:49:33Z` | crates.io API | `curl -s -H "User-Agent: research" https://crates.io/api/v1/crates/tauri-plugin-updater` |
| TS-06 | `tauri-apps/tauri-action` | `action-v1.0.0` | GitHub Releases API | `curl -s https://api.github.com/repos/tauri-apps/tauri-action/releases/latest \| grep tag_name` |
| TS-07 | `creativeprojects/go-selfupdate` | `v1.6.0` | GitHub Tags API | `curl -s "https://api.github.com/repos/creativeprojects/go-selfupdate/tags?per_page=5" \| grep '"name"'` |
| TS-08 | go-selfupdate deps chính | `Masterminds/semver/v3 v3.5.0`, `google/go-github/v86 v86.0.0`, `ulikunitz/xz v0.5.15`, `golang.org/x/crypto v0.53.0`; go directive `1.25.12` ⚠️ | `go.mod` @ tag v1.6.0 | `curl -s https://raw.githubusercontent.com/creativeprojects/go-selfupdate/v1.6.0/go.mod` |
| TS-09 | Go toolchain stable | `go1.26.5` (mới nhất), `go1.25.12` (previous) | go.dev download API | `curl -s "https://go.dev/dl/?mode=json" \| python3 -c "import sys,json;print([x['version'] for x in json.load(sys.stdin)])"` |
| TS-10 | `jedisct1/go-minisign` | `0.2.7` (tag KHÔNG có `v` prefix) | GitHub Tags API | `curl -s "https://api.github.com/repos/jedisct1/go-minisign/tags?per_page=3" \| grep '"name"'` |
| TS-11 | `minisign-verify` (Rust, Tauri dùng) | `0.2.5` | crates.io API | `curl -s -H "User-Agent: research" https://crates.io/api/v1/crates/minisign-verify` |
| TS-12 | `minisign` CLI | `0.12` | GitHub Releases API | `curl -s https://api.github.com/repos/jedisct1/minisign/releases/latest \| grep tag_name` |
| TS-13 | Tauri updater manifest format | `latest.json` với keys bắt buộc: `version`, `platforms.<target>.url`, `platforms.<target>.signature`; optional `notes`, `pub_date`. Config: `createUpdaterArtifacts: true`, `pubkey` (**phải là nội dung key, KHÔNG phải path**), `endpoints`, `windows.installMode` | Tauri v2 docs — https://v2.tauri.app/plugin/updater | Đọc lại trang docs |
| TS-14 | Tauri signature scheme | **Minisign (Ed25519)**. Key sinh bằng `tauri signer generate`. Env khi build: `TAURI_SIGNING_PRIVATE_KEY`, `TAURI_SIGNING_PRIVATE_KEY_PASSWORD` | Tauri v2 docs (Signing Updates) | Đọc lại trang docs |
| TS-15 | Tauri capabilities cần thiết | `updater:default`, `process:default`, `process:allow-restart` trong `src-tauri/capabilities/*.json` | Tauri v2 docs / CrabNebula guide | Đọc lại trang docs |

⚠️ **TS-08 cần re-verify:** dòng `go`/`toolchain` trong go.mod bị truncate khi đọc qua curl. Phải xác nhận chính xác ở **P0-T3** sau khi cài Go (`go mod download && go list -m all`). Tạm coi yêu cầu là **Go ≥ 1.25**.

---

## B. Môi trường máy agent (verify 2026-07-29)

| Thành phần | Trạng thái |
|---|---|
| OS | Linux (container, kernel 5.15, x86_64) — **KHÔNG phải Windows 11** |
| Rust | `rustc 1.94.1` / `cargo 1.94.1` ✅ |
| Node / npm | `v22.22.2` / `10.9.7` ✅ |
| Git | `2.34.1` ✅ |
| **Go** | ❌ chưa cài → P0-T3 phải cài |
| **gh** (GitHub CLI) | ❌ chưa cài |
| **pwsh** (PowerShell 7) | ❌ chưa cài → không lint được script `.ps1` tại local |
| **mingw / wine** | ❌ chưa có → **không thể build hoặc chạy Windows .exe tại local** |

**Hệ quả bắt buộc (rất quan trọng cho việc lập kế hoạch):**

Agent **không thể** tự chạy app Windows để kiểm chứng. Do đó mọi bước kiểm chứng phải được gán 1 trong 3 tier:

| Tier | Ai/cái gì kiểm chứng | Dùng cho |
|---|---|---|
| **T1 — Local (Linux)** | Agent chạy trực tiếp: `cargo test`, `go test`, JSON-schema validate, `actionlint`, grep secret, build Linux target | Logic thuần: semver compare, parse manifest, SHA-256, verify minisign, state machine, policy |
| **T2 — CI (GitHub Actions `windows-latest`)** | Agent đẩy code → CI chạy build + test thật trên Windows và trả log | Build Windows x64, `go test`/`cargo test` trên Windows, đóng gói artifact, ký, publish release, test file-lock/staged-install |
| **T3 — User (máy Windows 11 thật)** | Người dùng chạy tay và chụp/paste kết quả | UAC prompt, SmartScreen, antivirus quarantine, tắt máy giữa lúc install, cảm nhận UI |

→ **T2 là công cụ kiểm chứng chính.** Repo đã có remote `github.com/o25160526-pip/self-update-go-rust` nên CI khả dụng. Không được đẩy code lên remote mà chưa có approval của user (xem GATE-1 trong PLAN.md).

---

## C. Kết luận kỹ thuật rút ra từ verify (không được đảo ngược mà không cập nhật file này)

### C-1. Bản Rust: dùng Tauri 2 + plugin updater chính thức ✅
`tauri 2.11.5` + `tauri-plugin-updater 2.10.1` là stable và có sẵn download/verify/install. Theo §3 của prompt ("không tự viết lại nếu plugin đã hỗ trợ") → **dùng plugin, không tự viết**.

### C-2. Bản Go: **KHÔNG có updater core chính thức** ⚠️ — giới hạn phải ghi rõ
Prompt gợi ý "Wails v3 `app.Updater`". Verify TS-02/TS-03 cho thấy:
- Wails v3 vẫn **Alpha** (alpha.102), API có thể breaking bất kỳ lúc nào.
- Trong `v3/pkg/services` **không tồn tại** updater service.

→ Không được giả vờ là có (yêu cầu §3 của prompt). **Quyết định:**
- Desktop framework Go = **Wails v2.13.0** (stable, TS-01).
- Updater core Go = **`creativeprojects/go-selfupdate v1.6.0`** (TS-07) cho phần *check + download release từ GitHub*.
- Phần **staged install + restart + rollback** go-selfupdate không cover → tự viết adapter `updater/install|restart|rollback` **có unit test**, đúng như prompt cho phép.

### C-3. Thống nhất chữ ký: Minisign cho cả hai bản ✅
Tauri dùng Minisign/Ed25519 (TS-14). Go verify được **cùng format** bằng `go-minisign 0.2.7` (TS-10).
→ **Một cặp key duy nhất, một tool `minisign` 0.12 (TS-12), một format signature** cho cả Go và Rust. Đây là điểm mấu chốt tạo ra tính "tương đương" giữa hai bản mà prompt §4 yêu cầu.

### C-4. Manifest: publish **2 file** cho mỗi release
Vì format bắt buộc của Tauri (TS-13) ≠ format prompt §7 yêu cầu:

| File | Ai đọc | Format |
|---|---|---|
| `latest-rust.json` | `tauri-plugin-updater` | **Native Tauri**: `{version, notes, pub_date, platforms.<target>.{signature,url}}` |
| `manifest-go.json` | Go app (updater tự viết) | **Superset theo prompt §7**: `{version, channel, publishedAt, releaseNotes, minSupportedVersion, platform, url, sha256, signature, size, mandatory}` |

→ Bảng mapping trường giữa 2 format là **deliverable bắt buộc** trong `docs/architecture.md` (P1-T4).

### C-5. Cơ chế install trên Windows khác nhau → phải khai báo trong `docs/source-map.md`

| | Go (Wails v2) | Rust (Tauri 2) |
|---|---|---|
| Install | **Launcher + versioned dirs + pointer file** (tự viết) | **NSIS installer** do Tauri bundler tạo, `installMode: passive` |
| Tránh tự ghi đè exe đang chạy | Bằng thiết kế: mỗi version 1 thư mục riêng, chỉ đổi con trỏ | Bằng NSIS: app thoát rồi installer mới chạy |
| Rollback | Đổi pointer về `last-known-good` (rất nhanh, atomic) | Chạy lại installer của version trước (đã cache sẵn) |

---

## D. Nguồn tài liệu chính thức (dùng khi cần tra cứu sâu)

| Chủ đề | URL |
|---|---|
| Tauri v2 updater plugin | https://v2.tauri.app/plugin/updater |
| Tauri v2 capabilities/permissions | https://v2.tauri.app/security/capabilities |
| tauri-action (release CI) | https://github.com/tauri-apps/tauri-action |
| Wails v2 docs | https://wails.io/docs/introduction |
| Wails repo (kiểm tra v3 status) | https://github.com/wailsapp/wails |
| go-selfupdate | https://github.com/creativeprojects/go-selfupdate |
| Minisign spec/tool | https://jedisct1.github.io/minisign/ |
| WebView2 runtime (Windows) | https://developer.microsoft.com/microsoft-edge/webview2/ |
| GitHub Releases REST API | https://docs.github.com/rest/releases/releases |

**Quy tắc:** khi copy code mẫu từ docs → ghi URL + ngày truy cập vào comment ngay trên đoạn code đó.
