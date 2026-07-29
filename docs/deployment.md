# Deployment Guide

## Demo deployment

1. Use the public demo key only for validation.
2. Push `release/vMAJOR.MINOR.PATCH`.
3. Wait for both Windows workflows.
4. Confirm the release contains Go and Rust artifacts, manifests, signatures, and checksums.
5. Test the previous release on Windows 11 before calling the rollout complete.

## Production deployment

1. Generate a new Minisign/Tauri signing key outside the repository.
2. Store the private key and password in GitHub Actions Secrets: `TAURI_SIGNING_PRIVATE_KEY` and `TAURI_SIGNING_PRIVATE_KEY_PASSWORD`.
3. Replace the pinned public key in Go, Rust logic, and Tauri config in one coordinated release.
4. Copy `shared/update-server.example.json` to `update-server.json` and set the approved host.
5. Serve manifests and binaries over HTTPS with immutable version paths, correct content types, and no authentication requirement for clients.
6. Publish a canary release, test update and rollback, then promote the same signed artifacts.
7. Keep old manifests and artifacts available for at least one rollback window.

## Host requirements

- HTTPS, stable URLs, byte-for-byte artifact integrity.
- `manifest-go.json` for Go and `latest-rust.json` for Tauri.
- Correct `url`, version, platform, signature, checksum, size, and release notes.
- No redirects to private or expiring URLs.
- Monitoring for 4xx/5xx, download failures, and signature rejection.

## Rollback of a bad release

Do not mutate an existing release. Publish the next patch version with the fix, disable the bad version at the server/catalog layer, and keep the previous known-good artifact available.
