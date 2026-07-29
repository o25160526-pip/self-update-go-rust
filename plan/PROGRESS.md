# PROGRESS.md — Theo dõi tiến độ thực hiện

> Mirror trong ClickUp: History Work — self-update-go-rust.

## GOAL hiện tại

✅ Go và Rust chạy trên Windows 11 ở v1.0.3; Go update và Rust UI/update đã được user xác nhận.

## Trạng thái

| Hạng mục | Trạng thái | Ghi chú |
|---|---|---|
| P0 Build Windows | ✅ | Go/Rust CI xanh |
| P1 Hello version | ✅ | User xác nhận v1.0.3 |
| P2 Check update + manifest | ✅ | Releases + manifest + signatures |
| P3 Download/verify/install/restart | ✅ | Go user-confirmed; Rust signed package validated |
| P4 Health-check/rollback | 🔄 | Manual rollback menu added; Rust installer cache still required for real rollback |
| P5 Docs/deployment/template | ✅ | Docs and AGENTS.md shipped |
| N1 Multi-host update server | 🔄 | Go runtime config shipped; failover/schema next |
| N2 Product template | 🔄 | Guide shipped; bootstrap next |
| N3 Production hardening | ⏳ | Rotate keys, Authenticode, SBOM, protected environment |
| N4 Manual rollback menu | 🔄 | Go real restore; Rust cached-installer path with explicit error when absent |

## Next work

1. Cache the previous signed Rust installer during update so the rollback menu can restore it.
2. Add shared config schema validation and host priority/failover.
3. Generate Rust endpoints from the shared server config.
4. Enable GitHub Template Repository and product bootstrap tooling.
5. Add Windows 11 smoke runner and production signing/SBOM.
