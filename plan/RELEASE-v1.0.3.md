# v1.0.3 release validation

This patch fixes the Rust frontend startup bridge by enabling `app.withGlobalTauri=true` and adding explicit bridge diagnostics. The release must pass both Windows workflows before publication.

## T3 checklist

- Rust shows `Hello, version 1.0.3`.
- OS/Arch are populated.
- State leaves `checking`.
- Updater log is visible.
- Rust v1.0.2 can update to v1.0.3.
