# CI upgrade and Azure Pipeline

## GitHub Actions

- `actions/checkout@v7`
- `actions/setup-node@v7` with Node 24
- `actions/upload-artifact@v4.6.2`
- GitHub runners now default JavaScript actions toward Node 24; Node 20 warnings are from older action majors.
- Rustfmt runs first, then a check gate, so formatting drift is corrected in the workspace and verified.

## Azure Pipelines

`azure-pipelines.yml` mirrors the build order: checkout, Node 24, Go, Rust, tests, Windows packaging, and separate Go/Rust pipeline artifacts. Azure DevOps Services uses `UseNode@1` and `PublishPipelineArtifact@1`; Azure DevOps Server should use Build Artifacts instead.

## Update source policy

Clients can scan multiple HTTPS manifest sources. GitHub is priority 0 and is checked first for exact rollback targets. If it does not contain the requested version, configured Azure and other sources are scanned. For normal updates, all healthy sources are scanned and the highest valid SemVer is selected.

## Timeout fix

The old Go client used a 30-second HTTP client timeout for the entire artifact download. Large Windows executables hit `context deadline exceeded`. The client now allows 10 minutes and retains context cancellation, while manifest validation and signature verification remain mandatory.
