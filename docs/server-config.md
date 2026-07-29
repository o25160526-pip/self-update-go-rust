# Update Server Configuration

The template supports moving the update catalog away from GitHub Releases.

## Go runtime configuration

Copy `shared/update-server.example.json` to `update-server.json` beside the Go executable, or set `GO_DEMO_UPDATE_SERVER_CONFIG` to an absolute path. The active host determines the Go manifest URL. Missing, invalid, or incomplete configuration safely falls back to the built-in GitHub URL.

Example:

```json
{
  "activeHost": "staging",
  "hosts": {
    "staging": {
      "baseUrl": "https://updates.example.com/app/stable",
      "goManifest": "manifest-go.json",
      "rustManifest": "latest-rust.json"
    }
  }
}
```

Only HTTPS hosts should be used. The client still validates manifest and artifact URLs, SHA-256, Minisign signature, platform, channel, and SemVer.

## Rust build configuration

Tauri updater endpoints are compiled into `tauri.conf.json`. For a non-GitHub host, update the endpoint during the release build or maintain a product-specific Tauri config overlay. The public key remains pinned and must match the signing key used by the server.

## Host migration

1. Publish the same signed manifest and artifact on the new host.
2. Test the new host with a staging `update-server.json`.
3. Release a client configured for the new Rust endpoint.
4. Keep the old host online until all supported clients have migrated.
