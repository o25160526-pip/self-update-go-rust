# keys/ — Minisign key cho update signing

## TL;DR

| File | Nội dung | Bí mật? |
|---|---|---|
| `demo-signing.pub` | Public key (dạng file minisign 2 dòng) | Không |
| `demo-signing.key` | Private key **đã mã hoá** bằng password bên dưới | Demo only — **không dùng cho production** |
| `demo-signing.key.b64` | Base64 của toàn bộ file `.key` — đúng format mà `TAURI_SIGNING_PRIVATE_KEY` cần | Demo only |
| `demo-signing.password` | Password để giải mã private key (`demo-password`) | Demo only |

- **Key ID:** `5756DB5A3509A8C1`
- **Thuật toán:** Minisign / Ed25519, private key mã hoá bằng scrypt (chuẩn minisign C).

Public key này được **pin** ở 3 chỗ, cả 3 phải luôn khớp nhau:

1. `go-demo/updater/keys.go` → `PublicKey`
2. `rust-demo/logic/src/keys.rs` → `PINNED_MINISIGN_PUBKEY`
3. `rust-demo/src-tauri/tauri.conf.json` → `plugins.updater.pubkey`

## Cơ chế fallback trong CI

Hai workflow (`release-go.yml`, `release-rust.yml`) resolve key theo thứ tự:

1. Nếu có secret **`TAURI_SIGNING_PRIVATE_KEY`** → dùng secret đó (+ `TAURI_SIGNING_PRIVATE_KEY_PASSWORD`).
2. Nếu không có → **tự động dùng demo key trong repo này**.

Nghĩa là repo build + ký + release được ngay khi clone về, không cần cấu hình gì. Khi triển khai thật thì chỉ cần thêm secret, không phải sửa workflow.

## Khi triển khai thật (production)

```powershell
# PowerShell, trong thư mục rust-demo
npx tauri signer generate -w "$env:USERPROFILE\.tauri\myapp.key"

# Lấy giá trị để dán vào GitHub Secrets
[Convert]::ToBase64String([IO.File]::ReadAllBytes("$env:USERPROFILE\.tauri\myapp.key"))
```

Rồi:

1. Settings → Secrets and variables → Actions → thêm `TAURI_SIGNING_PRIVATE_KEY` (giá trị base64 ở trên) và `TAURI_SIGNING_PRIVATE_KEY_PASSWORD`.
2. Thay public key mới vào **cả 3** chỗ pin phía trên.
3. Xoá thư mục `keys/` này (hoặc để lại nhưng nhớ: nó không còn được dùng khi secret tồn tại).

> ⚠️ Private key demo này công khai trên GitHub. Bất kỳ ai cũng ký được artifact hợp lệ với nó. Chỉ dùng để demo/test luồng update.
