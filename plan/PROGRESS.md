# PROGRESS.md — Theo dõi tiến độ thực hiện

> Mirror trong ClickUp: History Work — self-update-go-rust.

## GOAL hiện tại

✅ Go và Rust chạy trên Windows 11 ở v1.0.3; Go update đã được user xác nhận; Rust UI/update đã được user xác nhận.

## Trạng thái

| Hạng mục | Trạng thái | Ghi chú |
|---|---|---|
| P0 Build Windows | ✅ | Go/Rust CI xanh |
| P1 Hello version | ✅ | User xác nhận v1.0.3 trên Windows 11 |
| P2 Check update + manifest | ✅ | v1.0.0 → v1.0.1 → v1.0.2/v1.0.3 releases |
| P3 Download/verify/install/restart | ✅ | Go update user xác nhận; Rust package signed |
| P4 Health-check/rollback | 🔄 | Logic T1 pass, Windows rollback automation còn là roadmap |
| P5 Docs/deployment/template | ✅ | Root AGENTS.md + deployment/server/template/roadmap docs |
| N1 Multi-host update server | 🔄 | Go runtime config shipped; Rust endpoint overlay planned |
| N2 Product template | 🔄 | Guide shipped; GitHub Template Repository setting remains manual |
| N3 Production hardening | ⏳ | Rotate keys, Authenticode, SBOM, protected environment |

## Releases

- v1.0.0: ✅ published baseline
- v1.0.1: ✅ published, Go update path validated
- v1.0.2: ✅ published, Rust UI bug found
- v1.0.3: ✅ published, Rust UI fixed and user validated

## Next plan

1. Add shared schema validation and host priority/failover for Go and Rust.
2. Add Rust endpoint generation from `update-server.json` during release builds.
3. Enable GitHub Template Repository and add product bootstrap tooling.
4. Rotate demo key before any production use; add Authenticode and SBOM.
5. Add Windows 11 smoke runner for install, UAC, update, restart, and rollback.
