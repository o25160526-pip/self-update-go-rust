# Rollback & health-check

## Manual rollback menu

Both apps now expose **Rollback phiên bản trước** in the updater menu.

### Go
The Go action is real and local-only. It selects the newest `*.backup-<version>` created by a successful update, restores it atomically, writes `status=rolled-back`, starts the restored executable, and exits the current process. It never downloads a binary for rollback.

If no backup exists, the menu reports `no previous local version is available`.

### Rust
The Rust action launches `%LOCALAPPDATA%\\rust-demo\\updates\\previous-installer.exe` when that cached installer exists, then exits so NSIS can restore the prior install. If the cache is absent, the menu reports the exact path and does not pretend rollback succeeded. The next hardening step is to cache the previous installer during every Tauri update.

## Automatic rollback

Go uses a watchdog and health marker. Rust uses the state adapter and signed installer flow. Manual rollback requires confirmation and remains local-only.
