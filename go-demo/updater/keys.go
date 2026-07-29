package updater

// PublicKey là Minisign public key được pin trong client, dạng base64 của TOÀN BỘ
// file .pub (giống format mà Tauri updater yêu cầu ở `plugins.updater.pubkey`).
//
// Key này PHẢI khớp với:
//   - rust-demo/logic/src/keys.rs        (PINNED_MINISIGN_PUBKEY)
//   - rust-demo/src-tauri/tauri.conf.json (plugins.updater.pubkey)
//   - keys/demo-signing.pub
//
// Private key tương ứng: keys/demo-signing.key (password: keys/demo-signing.password).
// Khi triển khai thật: sinh key mới, thay hằng này, và đặt private key vào
// GitHub Secrets TAURI_SIGNING_PRIVATE_KEY (CI tự ưu tiên secret nếu có).
// KHÔNG đổi key mà không có kế hoạch key rotation: client cũ sẽ từ chối update.
const PublicKey = "dW50cnVzdGVkIGNvbW1lbnQ6IG1pbmlzaWduIHB1YmxpYyBrZXk6IDU3NTZEQjVBMzUwOUE4QzEKUldUQnFBazFXdHRXVjgvL3ZSVks0Ky9EaVhWTTQ5aVJpeTBuNDdzZUlFS0p0SlhGYmZVWkYrN0oK"

// PinnedMinisignPublicKey là dạng text (2 dòng) của cùng public key, dùng cho log và docs.
const PinnedMinisignPublicKey = "untrusted comment: minisign public key: 5756DB5A3509A8C1\nRWTBqAk1WttWV8//vRVK4+/DiXVM49iRiy0n47seIEKJtJXFbfUZF+7J"

// PinnedKeyID là key ID (hex) của public key đang pin.
const PinnedKeyID = "5756DB5A3509A8C1"
