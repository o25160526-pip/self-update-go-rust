# Roadmap: Next Steps

## N1: Multi-host update catalogs

Add a shared schema validator, host priority/failover, configurable channels, and a Rust endpoint overlay generated from the same server config. Acceptance: staging and production hosts can be switched without source edits; invalid hosts are rejected.

## N2: Product template packaging

Enable GitHub Template Repository, add a product bootstrap script, rename placeholders, and generate a product checklist. Acceptance: a new product can be created by changing identity/config and adding business modules only.

## N3: Stronger production security

Rotate demo keys out of production, add Authenticode signing, provenance/SBOM, protected environments, staged rollout, and release approval. Acceptance: no demo private key is used in production and release evidence is auditable.

## N4: Operational observability

Add structured local logs, opt-in update telemetry, server health dashboards, and failure alerts without collecting user data by default.

## N5: Windows acceptance automation

Add a self-hosted Windows 11 smoke runner for install, UAC, WebView2, update, restart, and rollback. Acceptance: GATE-1/3/4 can be evidenced automatically, not only by manual screenshots.
