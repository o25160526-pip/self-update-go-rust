# Template Guide

This repo is the updater platform. A product fork should add business functionality, not duplicate update machinery.

## Fork checklist

1. Copy the repository into a new product repository or use it as a GitHub template.
2. Rename `go-demo` and `rust-demo` product identifiers.
3. Change application name, identifier, icons, window title, and manifest repository.
4. Generate a product signing key and pin only its public key.
5. Configure `update-server.json` and the Rust endpoint overlay.
6. Add business modules under `go-demo/features` and `rust-demo/src-tauri/src/features`.
7. Keep updater tests and add product smoke tests.
8. Create product-specific release notes and Windows 11 validation evidence.

## Extension contract

Business modules may depend on version, updater state, logs, and policy APIs. They must not edit state files, replace executables, bypass signature verification, or call release APIs directly.

## Product-specific work

- UI and domain workflows
- Local storage and migrations
- Authentication and authorization
- Telemetry with privacy review
- Product installer branding
- Product-specific health checks

## Platform-owned work

- Signed update manifests
- Download and verification
- Install/restart/watchdog
- Rollback and last-known-good state
- CI release workflows
- Security and deployment docs
