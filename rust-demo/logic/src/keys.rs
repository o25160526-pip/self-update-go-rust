// keys.rs — pinned minisign public key cho update verification.
//
// PHẢI khớp với:
//   - go-demo/updater/keys.go              (PublicKey / PinnedMinisignPublicKey)
//   - src-tauri/tauri.conf.json            (plugins.updater.pubkey)
//   - keys/demo-signing.pub
//
// Private key demo: keys/demo-signing.key (password: keys/demo-signing.password).
// Production: sinh key mới, thay hằng này, đặt private key vào GitHub Secrets
// TAURI_SIGNING_PRIVATE_KEY (CI tự ưu tiên secret nếu có, không có thì dùng demo key).

pub const PINNED_MINISIGN_PUBKEY: &str = "untrusted comment: minisign public key: 5756DB5A3509A8C1\nRWTBqAk1WttWV8//vRVK4+/DiXVM49iRiy0n47seIEKJtJXFbfUZF+7J";

/// Key ID (hex) của public key đang pin.
pub const PINNED_KEY_ID: &str = "5756DB5A3509A8C1";
