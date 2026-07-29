# CloudBase Deployment And Operations

Use this reference for `cloudbaserc.json`, public domains, routes, deployment verification, file
inspection, recovery, undeploy, and ENV destruction.

## Contents

- Application configuration
- Local validation
- Domain and deployment order
- Routes and aliases
- Verification and recovery
- Undeploy and destruction

## Application Configuration

The workspace-root `cloudbaserc.json` is the only deployment declaration. A nested configuration is
not a subproject override and managed staging ignores it. Preserve unknown official fields while
reconciling known App and Function fields with the repository.

For a plain no-build static directory:

```json
{
  "$schema": "https://static.cloudbase.net/cli/cloudbaserc.schema.json",
  "version": "2.0",
  "app": {
    "serviceName": "my-site",
    "root": ".",
    "framework": "static",
    "installCommand": "",
    "buildCommand": "",
    "outputDir": ".",
    "deployPath": "/",
    "ignore": [
      "**/.git/**",
      "**/.gitignore",
      "**/.codebanana/**",
      "**/.coding-agent/**",
      "**/.vercel/**",
      "**/vercel.json",
      "user-guides/**"
    ]
  }
}
```

`serviceName` must be 1-40 lowercase letters, digits, or hyphens, start with a letter, and end with a
letter or digit. The output directory must exist and its root must contain `index.html`. Never use a
top-level `static`, `app.static`, a Function that serves frontend files, or a direct Hosting upload.

For React, Vue, Vite, Next.js, Nuxt, Angular, and other framework Apps, normally omit
`app.framework` and let CloudBase detect it. Remove a stale contradictory value. Keep
`framework: "static"` only for authored HTML/CSS/JavaScript with no framework or build step.

Set `app.root` to the directory containing the application package and framework configuration.
Install/build commands and `outputDir` are relative to that root; do not compensate for a wrong
root with duplicated paths. Keep `app.deployPath` as `/`.

A frontend-only Next.js project must use static export and deploy `out`. `.next` is internal build
state, not a static root. If the project needs SSR, API routes, Server Actions, or another server
runtime feature, deploy that behavior as an HTTP Function. Preserve existing static-export settings
when adding browser Auth.

Never write `envId` into the workspace file. The tools inject the current persisted ENV into an
ephemeral configuration on every operation.

## Local Validation

Before calling deployment:

1. Inspect repository layout, package manifests, lockfiles, framework configuration, scripts, and
   expected output.
2. Validate every App and Function package at its owning package root.
3. Use the lockfile's package manager and run the real build, compile, type-check, or test command.
4. Verify the configured output exists and is current; for a frontend App, verify root
   `index.html`.
5. Fix every nonzero result before deployment.

For npm, prefer `npm ci --no-audit --no-fund` when the lockfile is valid; repair dependency metadata
with `npm install --no-audit --no-fund` only when needed. A dry-run resolution check does not replace
the real install and build. Do not pipe command output through `tail`, `grep`, or another command
that can hide the failing exit code.

## Domain And Deployment Order

A public App or HTTP Function requires a reserved slug. Event-only deployments do not.

Choose a meaningful 3-63 character slug containing lowercase ASCII letters, digits, and hyphens,
starting and ending with a letter or digit. Prefer an explicit user choice, then conversation
history, project/package metadata, `app.serviceName`, or an HTTP Function name. Never include a URL,
domain suffix, uppercase letter, underscore, space, PII, account identifier, or secret.

Reserve before any operation that may create the ENV:

```text
cloudbase_domain(
  action="reserve",
  domain_slug="my-app-01",
  conflict_policy="suffix_if_taken"
)
```

Use `exact` for a user-selected slug and present `suggested_slugs` on conflict. Use
`suffix_if_taken` for an Agent-selected slug.

The full domain suffix is server configuration and is intentionally absent from this Skill. Never
read configuration to construct it. `cloudbase_domain` and `cloudbase_deployment` return the only
authoritative complete domain and URL.

After required Database, Auth, schema, and Storage are ready, call:

```text
cloudbase_deployment(action="deploy")
```

The tool deploys every workload declared by `cloudbaserc.json`, waits for provider readiness,
reconciles routes, registers the returned managed origin with Auth, and returns the canonical URL.
It does not delete previously deployed workloads merely because they are no longer declared.
Never substitute a provider-generated URL or a URL inferred from resource names.

## Routes And Aliases

Omit `routes` for these defaults:

- App only: `/` targets the App.
- One HTTP Function without an App: `/` targets that Function.
- App plus HTTP Functions: `/` targets the App; each Function uses `/<function-name>`.
- Event Function: no public route.

Multiple HTTP Functions without an App require an explicit complete map with exactly one `/` route.
A custom map contains every HTTP Function; the App, when present, still owns `/`.

Domain-only operations never redeploy workloads:

- Rename with `cloudbase_domain(action="rename", domain_slug=..., conflict_policy=...)`.
- Add another public name with `add_alias`; aliases mirror the primary route set.
- Promote an active route-matched alias with `set_primary`.
- Remove one alias with `remove_alias` and explicit confirmation.

One project has one canonical primary and may have multiple globally unique aliases when its
CloudBase plan permits. If `add_alias` returns `DOMAIN_QUOTA_EXCEEDED`, preserve the
`quota_blocked` reservation. Do not evict another domain. Ask the user to remove the blocked alias,
remove a chosen active alias, or upgrade the plan and retry.

A rename persists target routes before removing the old binding. If inspection reports
`operation_state="rename_pending"`, retry the same rename; do not start a different rename, poll
with fixed sleeps, edit workload names, or redeploy.

## Verification And Recovery

Use `cloudbase_deployment(action="status")` for actual provider workload and endpoint state. A
successful upload or build is insufficient: verify active workloads, every returned HTTPS URL, and
one real business operation.

For an App, the managed tool verifies the provider App before binding the domain and confirms a new
root artifact. If the App is not serving, repair the source or root configuration and redeploy the
same App. Do not undeploy first, switch static files to a Function, or diagnose the custom domain as
the source failure.

Use `cloudbase_deployment(action="inspect_files")` when asked which files the latest App version
deployed:

- `uploaded_files` is a bounded final-version preview.
- `source_files` is a bounded preview after staging and exclusions.
- Counts and `*_truncated` flags indicate omitted entries.
- `manifest_path` points to the complete workspace-relative structured result under
  `.codebanana/cloudbase/deployment-manifests`.

Do not infer an App manifest from Static Hosting listings; they can include stale objects and system
resources. Do not call raw provider APIs or delete individual Hosting files.

Use the project-scoped `cloudbase_command` allowlist for bounded Function inspection, invocation,
policy, and log queries. Its `args` are exact TCB tokens after the executable; include real flags
such as `--json` and never invent command actions. The server adds only the authoritative ENV,
temporary config, and credentials. The process runs from an isolated authentication directory, so
every local file argument must be an explicit absolute path inside the workspace. Do not pass
`tcb`, ENV ID, region, or credentials.

Logging lifecycle belongs to `cloudbase_deployment`, not `cloudbase_command`. Recording is disabled
by default and deploy never enables it. Inspect with `action="logging_status"`; only after explicit
user consent enable with `action="enable_logging", confirm=true`, and stop new recording with
`action="disable_logging", confirm=true`. Query Function logs with exact `logs search ... --json`
argv. Do not choose `fn log` by default: its `GetFunctionLogs` API is unsupported by some ENVs even
when the installed CLI is current. Diagnose a Function failure from source layout, runtime,
`scf_bootstrap`, dependency packaging, configuration, provider request ID, and bounded logs. Do
not replace it with CloudBase Run or another workload.

Managed tools poll official state with bounded timeouts. Do not add fixed sleeps or fabricate a
fallback URL. Retry only the documented resumable operation and retain returned request IDs.

## Undeploy And Destruction

After explicit confirmation:

```text
cloudbase_deployment(action="undeploy", confirm=true)
```

This removes workloads and public routes while preserving the integrated ENV and project mapping,
Database, Storage, Auth, Publishable Key, and domain reservations. Inspect status when removal is
partial. Release aliases individually before releasing the primary when the user explicitly wants
domain ownership removed.

Destroy only after the user explicitly confirms permanent loss of every workload, Database, object,
Auth configuration, key, and other ENV resource:

```text
cloudbase_deployment(action="destroy_environment", confirm=true)
```

Never infer this consent from undeploy, cleanup, database deletion, or domain release. The tool
accepts no ENV ID and deletes only the current project's authoritative ENV. It preserves the
immutable project alias and domain reservations for a replacement ENV.

Never use `cloudbase_command`, an official CLI, SDK, API, MCP, or console procedure for undeploy or
ENV destruction.
