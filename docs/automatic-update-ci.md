# Automatic Update CI

Every push to `main` or `release/v*`, and every pull request to `main`, runs `smoke-update.yml`.

## What is checked automatically

- Go update comparison: `1.0.3` → `1.0.4`.
- Go pending state: `pending=1.0.4`, `updatedFrom=1.0.3`.
- Go rollback state: restores `current=1.0.3`, status `rolled-back`, and limits attempts.
- Go signed artifact integration tests: Minisign verification, SHA-256, download, install, restart hook, health marker, and failure rollback.
- Rust state transition: `1.0.3` → `1.0.4` → rollback to `1.0.3`.
- Rust format, clippy, and logic tests.
- Windows build, packaging, signing, and artifact upload remain enforced by the Go/Rust release workflows on every push.

## Separation of concerns

Smoke tests use local fixtures and do not publish releases. Real GitHub/Azure update sources are queried only by release/tag flows or an explicit Windows acceptance run. This keeps normal development fast while still blocking regressions before merge.

## Failure policy

A failed smoke gate blocks the change. Fix the code or test, push again, and wait for Go, Rust, and smoke jobs to complete before considering the commit green.
