# AGENTS.md

## Project purpose
This repository is a reusable Windows desktop auto-update template. Product teams should add business modules without rewriting updater, signing, release, rollback, or health-check code.

## Required workflow
1. Read `plan/PROGRESS.md` and `docs/template-guide.md` before changes.
2. Keep Go and Rust behavior aligned through `docs/source-map.md`.
3. Run Go tests, Rust logic tests, formatting, and secret scanning before pushing.
4. Use SemVer releases and never reuse a published version.
5. Never commit production private keys. Demo keys are for tests only.
6. Update `docs/`, `plan/PROGRESS.md`, and ClickUp History Work in the same work session.
7. For updater changes, test: manifest validation, hash/signature rejection, restart, rollback, offline startup, and Windows packaging.

## Architecture rules
- Go remains portable Wails v2 with staged install and watchdog rollback.
- Rust remains Tauri 2 with the official updater plugin and an explicit rollback adapter.
- Update hosts are configured through `shared/update-server.example.json`; Go may load a local `update-server.json` at runtime. Rust release builds receive the configured endpoint through the workflow/configuration step.
- Clients pin the public verification key. Changing a production key requires a coordinated client release.

## Definition of done
A change is done only when code, tests, docs, process status, release notes, and Windows CI are updated and green.
