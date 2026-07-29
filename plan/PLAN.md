# PLAN.md — Demo auto-update qua GitHub Releases bằng Go và Rust

> **Trạng thái:** DRAFT — chờ user approval trước khi bắt đầu P0.
> **Mục tiêu theo prompt:** 11 bước end-to-end (§2) + 7 nhóm test A-G (§12) + Definition of Done (§13).
> **Mọi quyết định kỹ thuật** phải trace về `TRUST-SOURCES.md`. Nếu TRUST-SOURCES chưa có → verify trước, rồi cập nhật.

---

## Phân chia Phase & Gate

| Phase | Mục tiêu | Kết quả kiểm chứng chính | Gate (để sang Phase tiếp) |
|---|---|---|---|
| **P0** | Scaffolding + môi trường | Repo tạo, `cargo test`/`go test` pass (Linux), CI chạy được (Windows) | **GATE-0**: CI green trên `windows-latest` |
| **P1** | App hello version + UI tối giản | 2 app hiển thị "Hello, version X.Y.Z" | **GATE-1**: 2 app chạy được trên Windows 11 (user xác nhận) |
| **P2** | Check update + manifest | App phát hiện "up-to-date" hoặc "update-available" | **GATE-2**: App cũ phát hiện bản mới trên GitHub Release |
| **P3** | Download + verify + staged install + restart | App tự download, verify chữ ký, install, restart → version mới | **GATE-3**: Full flow B01-B07 pass |
| **P4** | Rollback + health-check | App crash sau update → rollback về last-known-good | **GATE-4**: Test E01-E08 pass |
| **P5** | Test + docs + báo cáo cuối | 7 nhóm test A-G có kết quả, 7 docs file hoàn chỉnh | **GATE-5**: Definition of Done (§13) thỏa mãn |

### Nguyên tắc Gate
- **Mỗi Gate** phải được user xác nhận (trừ GATE-0 có thể tự verify bằng CI log).
- **Không được** bắt đầu Phase tiếp nếu Gate chưa pass.
- **Nếu phát hiện lệch hướng** (logic sai, API khác docs, framework không hỗ trợ) → dừng, ghi nhận vào `plan/PROGRESS.md` cột "Issues", báo user, chờ quyết định.

---

## Quy ước kiểm chứng (3 Tier)

| Tier | Môi trường | Phạm vi | Ai chạy |
|---|---|---|---|
| T1 | Local Linux (agent) | Logic thuần, test đơn vị, lint, grep secret | Agent tự chạy |
| T2 | GitHub Actions `windows-latest` | Build Windows, test, đóng gói, ký, publish | Agent push code → CI chạy → agent đọc log |
| T3 | Máy Windows 11 thật của user | UAC, SmartScreen, antivirus, cảm nhận UI, tắt máy giữa install | User chạy tay, chụp/paste kết quả |

→ Mỗi bước "Kiểm chứng" trong PLAN sẽ ghi rõ **[T1]**, **[T2]** hoặc **[T3]**.

---

## Quy ước ghi nhận tiến độ

File `plan/PROGRESS.md` sẽ được cập nhật **sau mỗi task hoàn thành**. Mỗi dòng có:

```
| [P0-T1] | ✅ Hoàn thành | 2026-07-29 14:00 | Kết quả: repo tạo, 15 file | Ghi chú: |
| [P0-T2] | 🔄 Đang làm    | 2026-07-29 14:30 | Kết quả: —                  | Ghi chú: đợi CI |
```

**Task sau** phải đọc `PROGRESS.md` trước khi bắt đầu để biết task trước đã hoàn thành kết quả gì.

---

## ═══════════════════════════════════════
## PHASE P0 — SCAFFOLDING & MÔI TRƯỜNG
## ═══════════════════════════════════════

### P0-T1: Tạo cấu trúc repo

**Chi tiết kỹ thuật:**
- Tạo toàn bộ thư mục theo layout §4 của prompt.
- Tạo `.gitignore` (ít nhất: `*.exe`, `*.msi`, `*.nsis`, `node_modules/`, `dist/`, `build/`, `target/`, `*.key`, `*.key.pub` riêng tư, `.env`).
- Tạo `README.md` skeleton (sẽ điền dần ở P5).
- Tạo `docs/source-map.md` skeleton (mapping Go ↔ Rust).
- Tạo `shared/update-policy.example.json` và `shared/release-manifest.example.json` theo §7.
- Tạo `shared/test-cases.md` skeleton (7 nhóm A-G).
- Tạo `plan/PROGRESS.md` khởi tạo.

**Kiểm chứng [T1]:**
- [ ] `find . -type d` khớp layout §4 (so sánh bằng script).
- [ ] `git status` cho thấy tất cả file mới.
- [ ] `.gitignore` cover `*.key`, `*.env`, `node_modules/`.

---

### P0-T2: Khởi tạo Rust app (Tauri 2)

**Chi tiết kỹ thuật:**
- Chạy `npm create tauri-app@latest rust-demo -- --template vanilla-ts --manager npm`.
  - Nếu `create-tauri-app` không hỗ trợ CLI non-interactive → tạo thủ công bằng `cargo init` + cấu hình `Cargo.toml`, `tauri.conf.json`, `package.json`.
- `Cargo.toml`: khóa `tauri = "2.11.5"`, `tauri-plugin-updater = "2.10.1"`.
- `rust-toolchain.toml`: `channel = "1.94.1"` (match rustc hiện tại).
- `package.json`: khóa Node 22, npm 10.
- Khởi tạo `src-tauri/capabilities/default.json` với `updater:default`, `process:default`, `process:allow-restart`.
- `src-tauri/tauri.conf.json`: thêm `"updater": { "active": true, "endpoints": [], "pubkey": "PLACEHOLDER", "windows": { "installMode": "passive" } }`.
- Tạo `src-tauri/src/version.rs` → `pub const VERSION: &str = "1.0.0";` (single source of truth).
- Tạo `frontend/index.html` skeleton hiển thị "Hello, version X.Y.Z" (đọc từ Tauri API).

**Kiểm chứng [T1]:**
- [ ] `cargo check` pass (trong `rust-demo/`).
- [ ] `grep -r "1.0.0" src-tauri/src/version.rs` → đúng.
- [ ] Không có `*.key` hoặc secret trong source.

**Kiểm chứng [T2]:**
- [ ] CI workflow `release-rust.yml` (chỉ build step, chưa release) → `cargo build --release` pass trên `windows-latest`.

---

### P0-T3: Cài Go + Khởi tạo Go app (Wails v2)

**Chi tiết kỹ thuật:**
- Cài Go: `go install golang.org/dl/go1.25.12@latest && go1.25.12 download && alias go=go1.25.12` (hoặc dùng `go1.26.5` nếu go-selfupdate tương thích).
  - **Verify TS-08**: sau khi cài Go, chạy `go mod download` trong `go-demo/` rồi `go list -m all` để xác nhận `go-selfupdate` resolve đúng `v1.6.0`.
- Cài Wails v2 CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.
- Khởi tạo app: `wails init -n go-demo -t vanilla` (hoặc tạo thủ công nếu CLI không hỗ trợ).
- Di chuyển file theo layout §4: `app/`, `frontend/`, `updater/`, `installer/`, `scripts/`, `build/`.
- `go.mod`: khóa `go 1.25.12`, require `github.com/creativeprojects/go-selfupdate v1.6.0`, `github.com/jedisct1/go-minisign v0.2.7`.
- Tạo `internal/version/version.go` → `const Version = "1.0.0"` (single source of truth).
- Tạo `frontend/index.html` skeleton hiển thị "Hello, version X.Y.Z" (đọc từ Wails runtime).

**Kiểm chứng [T1]:**
- [ ] `go build ./...` pass (trong `go-demo/`).
- [ ] `go list -m all | grep go-selfupdate` → `v1.6.0`.
- [ ] `grep -r "1.0.0" internal/version/version.go` → đúng.
- [ ] Không có secret trong source.

**Kiểm chứng [T2]:**
- [ ] CI workflow `release-go.yml` (chỉ build step) → `go build` pass trên `windows-latest`.

---

### P0-T4: Tạo GitHub Actions workflows (build-only)

**Chi tiết kỹ thuật:**

Tạo 2 file workflow chỉ build (chưa release), dùng để verify CI pipeline:

**`.github/workflows/release-go.yml`:**
```yaml
name: Build Go Demo
on:
  push:
    tags: ["v*"]
  workflow_dispatch:
    inputs:
      version:
        description: "Version (e.g. 1.0.0)"
        required: true
jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25.12"
      - uses: actions/setup-node@v4
        with:
          node-version: "22"
      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
      - name: Build
        working-directory: go-demo
        run: wails build -clean -ldflags "-X internal/version.Version=${{ github.event.inputs.version || '1.0.0' }}"
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: go-demo-windows-x64
          path: go-demo/build/bin/
```

**`.github/workflows/release-rust.yml`:**
```yaml
name: Build Rust Demo
on:
  push:
    tags: ["v*"]
  workflow_dispatch:
    inputs:
      version:
        description: "Version (e.g. 1.0.0)"
        required: true
jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@stable
        with:
          toolchain: "1.94.1"
      - uses: actions/setup-node@v4
        with:
          node-version: "22"
      - name: Install frontend deps
        working-directory: rust-demo
        run: npm install
      - name: Build Tauri
        uses: tauri-apps/tauri-action@v0
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tagName: ${{ github.event.inputs.version || '1.0.0' }}
          releaseName: "Rust Demo v__VERSION__"
          releaseBody: "See assets"
          releaseDraft: true
          prerelease: false
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: rust-demo-windows-x64
          path: rust-demo/src-tauri/target/release/bundle/
```

**Kiểm chứng [T1]:**
- [ ] `actionlint .github/workflows/release-go.yml` pass (hoặc validate YAML syntax).
- [ ] `actionlint .github/workflows/release-rust.yml` pass.
- [ ] Không có secret hardcode trong workflow.

**Kiểm chứng [T2]:**
- [ ] Trigger `workflow_dispatch` → build pass trên `windows-latest` cho cả 2 workflow.

---

### P0-T5: GATE-0 — CI green

**Điều kiện pass:**
- [ ] Workflow `release-go.yml` → build thành công, artifact upload.
- [ ] Workflow `release-rust.yml` → build thành công, artifact upload.
- [ ] Không có secret trong log CI.

**Nếu fail:** debug → fix → re-push → chờ CI lại. Không sang P1.

---

## ═══════════════════════════════════════
## PHASE P1 — APP HELLO VERSION + UI TỐI GIẢN
## ═══════════════════════════════════════

### P1-T1: Go app — hiển thị hello version

**Chi tiết kỹ thuật:**
- `frontend/index.html`: hiển thị "Hello, version X.Y.Z" + OS + architecture.
- `app/main.go`: bind version từ `internal/version.Version` vào Wails runtime.
- `app/main.go`: khởi tạo trạng thái updater = `checking` (sẽ triển khai ở P2).
- `frontend/`: hiển thị trạng thái updater (danh sách state §6), vùng log, release notes.
- `--version` flag: in version ra stdout và exit.
- `--print-update-state` flag: in trạng thái updater hiện tại.
- `--offline-test` flag: dùng mock provider (sẽ triển khai ở P2).

**Kiểm chứng [T1]:**
- [ ] `go build -ldflags "-X internal/version.Version=1.0.0" ./...` → build thành công.
- [ ] `./go-demo --version` → `1.0.0` (trên Linux, nếu cross-compile được).

**Kiểm chứng [T2]:**
- [ ] CI build → artifact `.exe` → user chạy → thấy "Hello, version 1.0.0".

**Kiểm chứng [T3]:**
- [ ] User chạy app trên Windows 11 → thấy UI đúng.

---

### P1-T2: Rust app — hiển thị hello version

**Chi tiết kỹ thuật:**
- `frontend/index.html`: hiển thị "Hello, version X.Y.Z" + OS + architecture.
- `src-tauri/src/main.rs`: expose version qua Tauri command `get_version`.
- `src-tauri/src/main.rs`: khởi tạo trạng thái updater = `checking`.
- `frontend/`: hiển thị trạng thái updater, vùng log, release notes.
- `--version` flag: in version ra stdout và exit.
- `--print-update-state` flag: in trạng thái updater hiện tại.
- `--offline-test` flag: dùng mock provider.

**Kiểm chứng [T1]:**
- [ ] `cargo build` pass.
- [ ] `./rust-demo --version` → `1.0.0` (trên Linux, nếu cross-compile được).

**Kiểm chứng [T2]:**
- [ ] CI build → artifact → user chạy → thấy "Hello, version 1.0.0".

**Kiểm chứng [T3]:**
- [ ] User chạy app trên Windows 11 → thấy UI đúng.

---

### P1-T3: Architecture doc + source-map

**Chi tiết kỹ thuật:**
- Viết `docs/architecture.md`:
  - Bảng framework & version đã khóa (trace → TRUST-SOURCES.md).
  - Bảng mapping trường manifest (prompt §7 ↔ Tauri native format).
  - Bảng cơ chế install/rollback khác nhau giữa Go và Rust (TRUST-SOURCES C-5).
  - Diagram: luồng update §8 (dùng Mermaid).
- Viết `docs/source-map.md`: bảng mapping Go ↔ Rust theo §4.

**Kiểm chứng [T1]:**
- [ ] Tất cả version trong doc khớp TRUST-SOURCES.md.
- [ ] Mermaid diagram render được.

---

### P1-T4: GATE-1 — 2 app chạy được

**Điều kiện pass:**
- [ ] User chạy Go app trên Windows 11 → thấy "Hello, version 1.0.0". [T3]
- [ ] User chạy Rust app trên Windows 11 → thấy "Hello, version 1.0.0". [T3]
- [ ] `docs/architecture.md` và `docs/source-map.md` hoàn chỉnh.

---

## ═══════════════════════════════════════
## PHASE P2 — CHECK UPDATE + MANIFEST
## ═══════════════════════════════════════

### P2-T1: Sinh key pair Minisign (cho test)

**Chi tiết kỹ thuật:**
- Cài `minisign` CLI (hoặc dùng Go library `jedisct1/go-minisign` để sinh key).
- Sinh key pair: `minisign -G -p .keys/demo.pub -s .keys/demo.key` (local only, gitignore).
- **Public key** → nhúng vào cả Go app và Rust app (hardcode const hoặc config file).
- **Private key** → lưu vào GitHub Actions Secrets `MINISIGN_PRIVATE_KEY` (và `MINISIGN_PRIVATE_KEY_PASSWORD` nếu có).
- Ghi hướng dẫn sinh key vào `docs/architecture.md` (phần §10).

**Kiểm chứng [T1]:**
- [ ] `minisign -V -p .keys/demo.pub -x <signature> -m <file>` → verify thành công.
- [ ] `.keys/demo.key` có trong `.gitignore`.
- [ ] `grep -r "demo.key" --include="*.go" --include="*.rs" --include="*.json"` → không tìm thấy private key.

---

### P2-T2: Go app — check update

**Chi tiết kỹ thuật:**
- Tạo `updater/check.go`:
  - Dùng `go-selfupdate` để query GitHub Releases API.
  - Lọc channel `stable`, OS `windows`, arch `x86_64`.
  - So sánh semver (dùng `Masterminds/semver`), chống downgrade.
  - Nếu không có bản mới → trạng thái `up-to-date`.
  - Nếu có bản mới → trạng thái `update-available`, hiển thị thông báo + release notes.
- Tạo `updater/policy.go`: đọc `update-policy.json` (nếu có), fallback default theo §7.
- Tạo `updater/manifest.go`: parse `manifest-go.json` (format §7), validate schema.
- Tạo `updater/mock.go`: mock provider cho `--offline-test` mode.

**Kiểm chứng [T1]:**
- [ ] `go test ./updater/... -v` → pass (test semver compare, filter channel/OS/arch, downgrade reject).
- [ ] Mock provider hoạt động: chạy app với `--offline-test` → trạng thái `up-to-date` hoặc `update-available` tùy mock data.

---

### P2-T3: Rust app — check update

**Chi tiết kỹ thuật:**
- Cấu hình `tauri-plugin-updater` trong `src-tauri/tauri.conf.json`:
  - `endpoints`: URL trỏ đến `latest-rust.json` trên GitHub Release.
  - `pubkey`: nội dung public key (KHÔNG phải path).
- `src-tauri/src/updater.rs`:
  - Khởi tạo updater plugin trong `setup()`.
  - Gọi `app.updater_builder().check().await` → trả về `Update` hoặc `UpToDate`.
  - Map trạng thái Tauri → trạng thái hiển thị §6.
  - Hiển thị thông báo + release notes.
- `src-tauri/src/policy.rs`: đọc policy local (nếu có), fallback default.
- `src-tauri/src/mock.rs`: mock provider cho `--offline-test`.

**Kiểm chứng [T1]:**
- [ ] `cargo test` pass (test updater init, check logic).
- [ ] Mock provider hoạt động.

---

### P2-T4: Publish manifest lên GitHub Release (v1.0.0)

**Chi tiết kỹ thuật:**
- Cập nhật 2 workflow để hỗ trợ full release (thêm sign, manifest, upload).
- Tạo tag `v1.0.0` → push → CI build + sign + publish.
- CI tạo `latest-rust.json` (native Tauri format) + `manifest-go.json` (format §7).
- CI upload cả 2 manifest + artifact + checksum + signature lên GitHub Release.

**Kiểm chứng [T2]:**
- [ ] Tag `v1.0.0` → CI tạo release thành công.
- [ ] Release có artifact `.exe`, `.msi`, `latest-rust.json`, `manifest-go.json`, checksum, signature.
- [ ] Manifest trỏ đúng artifact URL.

---

### P2-T5: GATE-2 — App phát hiện bản mới

**Điều kiện pass:**
- [ ] Go app `1.0.0` → check GitHub Releases → trạng thái `up-to-date` (chỉ có v1.0.0). [T3]
- [ ] Publish `v1.0.1` (build thủ công hoặc CI).
- [ ] Go app `1.0.0` → check → trạng thái `update-available`, hiển thị `1.0.1`. [T3]
- [ ] Rust app `1.0.0` → check → trạng thái `update-available`, hiển thị `1.0.1`. [T3]

---

## ═══════════════════════════════════════
## PHASE P3 — DOWNLOAD + VERIFY + INSTALL + RESTART
## ═══════════════════════════════════════

### P3-T1: Go app — download + verify

**Chi tiết kỹ thuật:**
- `updater/download.go`:
  - Download artifact vào `%LOCALAPPDATA%/go-demo/update-tmp/` (thư mục tạm riêng).
  - Kiểm tra file size khớp manifest.
  - Tính SHA-256 và so sánh với manifest → không khớp → reject.
  - Verify Minisign signature bằng `go-minisign` + public key pin → không khớp → reject.
  - Trạng thái: `downloading` → `verifying`.
- `updater/install.go`:
  - **Staged install pattern** (tự viết, vì Wails v2 không có updater core):
    1. Mỗi version nằm trong thư mục riêng: `%LOCALAPPDATA%/go-demo/versions/1.0.0/`.
    2. Download + verify xong → extract/copy vào `%LOCALAPPDATA%/go-demo/versions/1.0.1/`.
    3. Ghi state `pending` vào `%LOCALAPPDATA%/go-demo/state.json`.
    4. Tạo launcher script (`.bat` hoặc helper `.exe`) thực hiện swap pointer.
  - Pointer file: `%LOCALAPPDATA%/go-demo/current` → chứa đường dẫn đến version đang chạy.
  - Trạng thái: `installing`.
- `updater/restart.go`:
  - Flush dữ liệu, đóng local service (nếu có).
  - Chạy launcher: đổi pointer → `%LOCALAPPDATA%/go-demo/versions/1.0.1/`.
  - Restart app bằng `exec.Command` hoặc launcher helper.
  - Trạng thái: `restarting`.

**Kiểm chứng [T1]:**
- [ ] `go test ./updater/... -v` → pass (test download mock, SHA-256 mismatch, signature mismatch).
- [ ] Staged install logic: `state.json` đúng trạng thái `pending` sau khi verify.

**Kiểm chứng [T3]:**
- [ ] App tự download → verify → install → restart → hiển thị `1.0.1`.

---

### P3-T2: Rust app — download + verify + install + restart

**Chi tiết kỹ thuật:**
- `tauri-plugin-updater` tự handle download + verify (Minisign) + install.
- Cấu hình `installMode: "passive"` (tự install, không hỏi user).
- Sau khi download + verify xong → plugin gọi `update.download_and_install()`.
- Plugin tự restart app bằng `app.restart()`.
- Lưu bản cũ vào cache: `%LOCALAPPDATA%/rust-demo/updates/` (Tauri tự quản lý).

**Kiểm chứng [T1]:**
- [ ] `cargo test` pass (test updater flow với mock).

**Kiểm chứng [T3]:**
- [ ] App tự download → verify → install → restart → hiển thị `1.0.1`.

---

### P3-T3: Luồng update end-to-end (§8 steps 1-14)

**Chi tiết kỹ thuật:**
- Đảm bảo cả 2 app implement đầy đủ 14 bước đầu của §8:
  1. App đọc version build từ nguồn duy nhất. ✅ (P1)
  2. App đọc policy local. ✅ (P2)
  3. App gọi GitHub Releases/manifest URL qua HTTPS. ✅ (P2)
  4. Lọc đúng channel, OS, CPU. ✅ (P2)
  5. So sánh semver, chống downgrade. ✅ (P2)
  6. Nếu không có bản mới → `Up to date`. ✅ (P2)
  7. Nếu có bản mới → thông báo. ✅ (P2)
  8. Tự động tải vào thư mục tạm riêng. ✅ (P3)
  9. Kiểm tra size + SHA-256. ✅ (P3)
  10. Xác minh chữ ký Minisign. ✅ (P3)
  11. Ghi state `pending`. ✅ (P3)
  12. Flush dữ liệu, đóng app. ✅ (P3)
  13. Helper/launcher swap/install. ✅ (P3)
  14. Restart app với metadata `updated-from`. ✅ (P3)

**Kiểm chứng [T3]:**
- [ ] Chạy app `1.0.0` → publish `v1.0.1` → app tự update → sau restart thấy `1.0.1`. (Test B01-B07)

---

### P3-T4: GATE-3 — Full flow B01-B07 pass

**Điều kiện pass:**
- [ ] B01: `1.0.0`, không có release mới → `up-to-date`. [T3]
- [ ] B02: `1.0.0`, có `1.0.1` → `update-available`. [T3]
- [ ] B03: Auto download không cần click. [T3]
- [ ] B04: Hash đúng, signature đúng, install thành công. [T3]
- [ ] B05: App tự restart. [T3]
- [ ] B06: Sau restart hiển thị `1.0.1`. [T3]
- [ ] B07: State ghi `last-known-good=1.0.1`. [T3]

---

## ═══════════════════════════════════════
## PHASE P4 — ROLLBACK + HEALTH-CHECK
## ═══════════════════════════════════════

### P4-T1: Go app — health-check + rollback

**Chi tiết kỹ thuật:**
- `updater/health.go`:
  - Sau restart, app mới thực hiện startup health-check:
    - Kiểm tra UI render được.
    - Kiểm tra version mới khớp manifest.
    - Timeout: `healthCheckTimeoutSeconds` (mặc định 30s).
  - Nếu thành công → ghi `last-known-good = <version>` vào state.
  - Nếu thất bại (timeout/crash) → trigger rollback.
- `updater/rollback.go`:
  - Đổi pointer về `last-known-good` version.
  - Restart app.
  - `maxRollbackAttempts: 1` → không rollback lặp vô hạn.
  - Không xóa user data/config.
  - Trạng thái: `rolled-back`.

**Kiểm chứng [T1]:**
- [ ] `go test ./updater/... -v` → test health-check timeout, rollback logic, max attempts.

**Kiểm chứng [T3]:**
- [ ] App crash sau update → rollback về bản cũ → vẫn chạy được.

---

### P4-T2: Rust app — health-check + rollback

**Chi tiết kỹ thuật:**
- Tauri 2 updater plugin không có built-in health-check → phải tự viết adapter:
  - `src-tauri/src/updater/health.rs`:
    - Sau restart, kiểm tra app khởi động thành công.
    - Nếu thành công → ghi `last-known-good` vào local state.
  - `src-tauri/src/updater/rollback.rs`:
    - Nếu health-check thất bại → chạy lại installer của version trước (đã cache).
    - `maxRollbackAttempts: 1`.
    - Không xóa user data/config.
  - **Giới hạn ghi rõ:** Tauri 2 không có built-in rollback → tự viết adapter, có test.

**Kiểm chứng [T1]:**
- [ ] `cargo test` → test health-check, rollback logic.

**Kiểm chứng [T3]:**
- [ ] App crash sau update → rollback → vẫn chạy được.

---

### P4-T3: Test nhóm E (restart & rollback)

**Kiểm chứng [T3]:**
- [ ] E01: App đang giữ file → helper chờ app thoát.
- [ ] E02: App crash trước health-check → rollback.
- [ ] E03: App mới treo quá timeout → rollback.
- [ ] E04: (Không có local service trong demo → ghi chú "N/A, mở rộng sau").
- [ ] E05: Rollback về `last-known-good`.
- [ ] E06: Không rollback lặp vô hạn.
- [ ] E07: Máy tắt giữa lúc install → khởi động lại vẫn phục hồi.
- [ ] E08: Rollback không xóa user data/config.

---

### P4-T4: GATE-4 — Rollback hoạt động

**Điều kiện pass:**
- [ ] E01-E08 (có thể trừ E04) pass.
- [ ] Cả 2 app đều có health-check + rollback.

---

## ═══════════════════════════════════════
## PHASE P5 — TEST + DOCS + BÁO CÁO CUỐI
## ═══════════════════════════════════════

### P5-T1: Test nhóm A (app cơ bản)

**Kiểm chứng [T1/T2/T3]:**
- [ ] A01: Clean build Go. [T2]
- [ ] A02: Clean build Rust. [T2]
- [ ] A03: App hiển thị đúng version. [T3]
- [ ] A04: `--version` trả đúng version. [T3]
- [ ] A05: App khởi động offline. [T3]
- [ ] A06: App không ghi secret vào log. [T1] — `grep -r "secret\|key\|token" *.log`

---

### P5-T2: Test nhóm C (lỗi mạng) + D (bảo mật) + F (quyền Windows) + G (phát hành)

**Kiểm chứng — chủ yếu [T1] (mock) + [T3] (thực):**

**Nhóm C (mạng):**
- [ ] C01-C07: Dùng mock provider để simulate DNS fail, timeout, HTTP 404/500, download ngắt. [T1]

**Nhóm D (bảo mật):**
- [ ] D01: SHA-256 sai → update bị từ chối. [T1] — test unit.
- [ ] D02: Signature sai → update bị từ chối. [T1] — test unit.
- [ ] D03: Manifest bị sửa → update bị từ chối. [T1] — test unit.
- [ ] D04: Version thấp hơn → bị từ chối. [T1] — test unit.
- [ ] D05: Artifact sai OS/arch → bị từ chối. [T1] — test unit.
- [ ] D06: Private key không xuất hiện trong artifact/log. [T1] — `grep -r`.
- [ ] D07: Client không cần GitHub token với public release. [T1] — kiểm tra code không gửi token.

**Nhóm F (quyền Windows):**
- [ ] F01-F05: Chạy trên Windows 11 thật. [T3]

**Nhóm G (phát hành):**
- [ ] G01: Tag `v1.0.0` → tạo release. [T2]
- [ ] G02: Tag `v1.0.1` → artifact mới. [T2]
- [ ] G03: Artifact Go và Rust không ghi đè nhau. [T2]
- [ ] G04: Checksum upload đúng. [T2]
- [ ] G05: Manifest trỏ đúng artifact. [T2]
- [ ] G06: Workflow fail khi thiếu signing secret. [T2]
- [ ] G07: Release có release notes. [T2]

---

### P5-T3: Hoàn thiện docs

**Danh sách 7 docs file:**
1. `docs/architecture.md` — đã có skeleton (P1-T3), điền đầy đủ.
2. `docs/windows-11-build.md` — hướng dẫn theo §11, mỗi lệnh ghi rõ PowerShell hay Git Bash.
3. `docs/release.md` — hướng dẫn theo §9, local release + Git tag + CI.
4. `docs/update-flow.md` — mô tả luồng §8, state machine, diagram.
5. `docs/rollback.md` — mô tả cơ chế rollback, health-check, recovery.
6. `docs/testing.md` — 7 nhóm test A-G, mỗi test có mã, setup, command, expected, log, cleanup.
7. `docs/troubleshooting.md` — lỗi thường gặp: linker, WebView2, UAC, antivirus, file lock, signature mismatch.

**Kiểm chứng [T1]:**
- [ ] Mỗi file có đầy đủ nội dung theo prompt §11/§9/§8.
- [ ] Tất cả version trong docs khớp TRUST-SOURCES.md.
- [ ] Mỗi lệnh PowerShell có ghi rõ đang chạy ở PowerShell hay Git Bash.

---

### P5-T4: Báo cáo cuối cùng (§14)

**Nội dung:**
1. Cây thư mục thực tế.
2. Framework và version đã khóa (trace → TRUST-SOURCES.md).
3. Lệnh build local cho Go và Rust.
4. Lệnh tạo tag và release.
5. Tên artifact đã upload.
6. URL GitHub Release dùng để test.
7. Cách client phát hiện update.
8. Cách restart và rollback hoạt động.
9. Kết quả từng nhóm test A-G.
10. Điểm chưa tự động hóa + lý do.
11. Việc cần làm trước production.

**Kiểm chứng [T1]:**
- [ ] Báo cáo đầy đủ 11 mục.
- [ ] Kết quả test A-G có ghi rõ pass/fail/NA.

---

### P5-T5: GATE-5 — Definition of Done (§13)

**Checklist (tất cả phải ✅):**
- [ ] Hai app Go và Rust đều build được trên Windows 11 x64. [T2]
- [ ] Source layout và naming tương đương, có source-map. [T1]
- [ ] App hello version chạy được. [T3]
- [ ] Có GitHub Actions cho release. [T2]
- [ ] Có artifact và manifest/signature rõ ràng. [T2]
- [ ] Client cũ tự phát hiện release mới. [T3]
- [ ] Client tự download, verify, install và restart không cần xác nhận. [T3]
- [ ] Có state machine và log dễ kiểm tra. [T1]
- [ ] Có last-known-good và rollback khi startup health-check thất bại. [T3]
- [ ] Có test cho hash sai, signature sai, network fail, crash và quyền Windows. [T1/T3]
- [ ] Có docs từ clone source đến GitHub Release và client update thành công. [T1]
- [ ] Không có secret trong source, artifact hoặc log. [T1]
- [ ] Các giới hạn framework được ghi rõ, không dùng claim marketing thay cho test thực tế. [T1]

---

## ═══════════════════════════════════════
## PHỤ LỤC
## ═══════════════════════════════════════

### A. Kiểm tra lệch hướng

Sau mỗi Phase, agent phải tự kiểm tra:

| # | Câu hỏi | Nếu "Không" → hành động |
|---|---|---|
| L1 | Tất cả version trong code khớp TRUST-SOURCES.md? | Cập nhật code hoặc TRUST-SOURCES (nếu verify mới). |
| L2 | Layout source khớp §4 + source-map.md? | Điều chỉnh. |
| L3 | Luồng update khớp §8 (18 bước)? | Bổ sung bước thiếu. |
| L4 | Không có secret trong source/log/artifact? | `grep -r` → xóa. |
| L5 | Giới hạn framework được ghi rõ? | Thêm vào docs + PROGRESS.md. |
| L6 | Test A-G có kết quả quan sát được? | Bổ sung expected result. |

### B. Ước tính effort

| Phase | Số task | Ước tính phiên làm việc | Ghi chú |
|---|---|---|---|
| P0 | 5 | 1-2 phiên | Scaffolding + CI |
| P1 | 4 | 1-2 phiên | Hello version + docs |
| P2 | 5 | 2-3 phiên | Check update + manifest + release |
| P3 | 4 | 2-3 phiên | Download + verify + install + restart |
| P4 | 4 | 2-3 phiên | Rollback + health-check |
| P5 | 5 | 2-3 phiên | Test + docs + báo cáo |
| **Tổng** | **27** | **10-16 phiên** | |
