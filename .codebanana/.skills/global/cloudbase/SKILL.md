---
name: cloudbase
description: Use before building or changing an application that needs deployment, a backend/API/server runtime in any language, user authentication, persistent database or shared CRUD state, or object storage, even when the user does not name CloudBase. Also use for CloudBase troubleshooting, cloudbase_* tools, cloudbaserc.json, TCB, or cb-* ENV IDs. Do not use for UI-only or local static work with no deployment or managed resources, or for another provider unless CloudBase integration or migration is requested. Manage CloudBase only through the provided cloudbase_* tools.
---

# CloudBase Application Development

## Managed Boundary

Use only these tools for CloudBase control-plane work:

- `cloudbase_domain`
- `cloudbase_database`
- `cloudbase_auth`
- `cloudbase_storage`
- `cloudbase_deployment`
- `cloudbase_command`

They resolve the current globally unique project ID, its single integrated ENV, provider
credentials, ownership, and readiness. Never pass an ENV ID, region, Tencent credential, database
password, or provider resource owned by another project.

Do not use the official `tcb` CLI, Tencent Cloud SDK, CloudBase MCP, raw management API, console, or
an installed management client to bypass these tools. Use `cloudbase_command` only for its
project-scoped allowlist, such as Function inspection, invocation, policy, and logs. Its `args` are
100% real TCB argument tokens after the `tcb` executable: never invent actions, translate flags, or
omit an explicitly required TCB flag such as `--json`. The server only injects project
ENV, config, and credentials. Commands run from an isolated authentication directory; pass an
explicit absolute workspace path for every local file argument. It is never a deployment, domain,
credential, logging-lifecycle, database-lifecycle, or ENV-lifecycle bypass.

Application code may use the official CloudBase browser or runtime SDK for data-plane access.
Obtain browser initialization only from `cloudbase_auth(action="frontend_config")`. Never expose
Tencent AK/SK, database credentials, privileged SQL, or internal endpoints in application code,
workspace configuration, logs, or model-visible output.

CloudBase owns relational database principals. Never create, alter, rename, or drop users or roles;
change passwords or memberships; or modify credential catalogs. Use existing execution roles,
table/schema grants, and PostgreSQL RLS.

## Load Only Relevant References

Before implementing a matching concern, call the exact command shown. Keep the `references/`
prefix; a bare filename is invalid.

| Concern | Required load |
|---|---|
| Public deployment, `cloudbaserc.json`, domains, routes, status, file inspection, undeploy, or ENV destruction | `skill_manager(action="load", skill_name="cloudbase", file="references/deployment-operations.md")` |
| Any Event or HTTP Function, backend framework, runtime identity, invocation, or Function logs | `skill_manager(action="load", skill_name="cloudbase", file="references/function-runtime.md")` |
| Browser login, Publishable Key, session lifecycle, route guards, or Auth acceptance | `skill_manager(action="load", skill_name="cloudbase", file="references/auth-browser.md")` |
| Database, CRUD, RLS, caller-token handoff, object storage, or trusted data access | `skill_manager(action="load", skill_name="cloudbase", file="references/data-access.md")` |

Load multiple references when the architecture crosses concerns. Do not invent behavior when a
required reference has not loaded.

## Choose The Smallest Architecture

1. Use a CloudBase App for plain static files, SPA bundles, and static exports.
2. Use browser SDK access when Auth plus resource policy fully enforces ordinary owner-scoped CRUD
   or private file access. Do not create a Function merely to proxy this path.
3. Add a Function when the application needs trusted validation, cross-user operations, secrets,
   transactions, framework APIs, SSR, streaming, or provider-verified caller identity.
4. Use an Event Function for SDK invocation, provider-verified caller identity, events, and
   background handlers. Use an HTTP Function for public REST, GraphQL, SSR, SSE, WebSocket, or web
   frameworks such as Express, Next.js server runtime, FastAPI, Flask, or Django.
5. Use Database only for durable structured state and Storage only for object bytes. Keep searchable
   object metadata in Database.

ENV creation does not provision a relational database. Before any mutating Auth, Storage, or
deployment call creates a new ENV, choose PostgreSQL when relational persistence is required and no
compatibility requirement selects MySQL:

```text
cloudbase_database(action="ensure_database", engine="postgresql")
```

Inspect and preserve an existing engine. Never switch an existing project between MySQL and
PostgreSQL by inference. Never fall back from a failed Function to CloudBase Run; the managed
platform does not expose CloudBase Run or Custom Image HTTP Function deployment.

## Workload Contract

The workspace-root `cloudbaserc.json` is the only workload source of truth:

- A top-level `app` deploys one CloudBase App.
- A non-empty top-level `functions` array deploys all declared Functions.
- Both may coexist. Neither means deployment is invalid.
- A frontend-only site always uses `app`, never a Function or direct Hosting upload.
- Deployment never treats an undeclared existing workload as deletion consent.

Before deployment, inspect the repository, package manifest, framework configuration, build
scripts, and output. Create or reconcile the root configuration rather than trusting stale fields.
Preserve unknown official fields. `app.root` contains the application package; build commands and
`outputDir` are relative to it. Add a stable valid `app.serviceName`.

Omit `app.framework` for framework projects so CloudBase can detect it. Use `framework: "static"`
only for plain authored HTML/CSS/JavaScript without a framework or build step. A Next.js static
export deploys `out`, never `.next`; server features require an HTTP Function.

The workspace configuration never contains `envId`. Managed tools inject the authoritative ENV
through an ephemeral configuration, including after ENV replacement. Validate dependencies and run
the real local build or compile command without pipelines that can hide a failing exit code.

## Domain And Deployment Workflow

Every public App or HTTP Function needs a reservation before infrastructure provisioning. Pass only
a meaningful 3-63 character `domain_slug`: lowercase ASCII letters, digits, and hyphens, starting
and ending with a letter or digit. Never pass a URL, suffix, uppercase text, underscore, space, PII,
account identifier, or secret.

Use `conflict_policy="exact"` for a user-selected slug and show returned suggestions on conflict.
Use `conflict_policy="suffix_if_taken"` for an Agent-selected slug derived from the conversation,
project metadata, `app.serviceName`, or an HTTP Function name. Event-only deployments have no
public domain.

The domain suffix is Agent Server configuration. Never read, infer, append, persist, or hardcode it.
Only managed tool results supply complete domains and URLs. Treat the returned managed URL as
canonical for Auth origins, health checks, persistence, and user-facing output; never construct a
URL from naming conventions or return a provider-generated URL as the final application URL.

Standard sequence:

1. Inspect the repository and load all relevant references.
2. Reserve a slug for a public workload.
3. Ensure required Database, Auth policy, schema, and Storage.
4. Reconcile `cloudbaserc.json` and prove every package locally.
5. Call `cloudbase_deployment(action="deploy")`; omit `routes` for managed defaults.
6. Inspect actual state, test every returned public URL, and exercise one real business operation.

Default routes are `/` to the App, or `/` to the only HTTP Function when there is no App. With an
App, each HTTP Function uses `/<function-name>`. Multiple HTTP Functions without an App require a
complete custom route map with exactly one `/`. Event Functions have no public route.

Domain rename, aliases, and primary promotion do not redeploy workloads. Use `cloudbase_domain`
operations and the deployment reference's recovery rules.

## Operations And Safety

Use `cloudbase_deployment(action="status")` for provider workload and endpoint state. Use
`cloudbase_deployment(action="inspect_files")` for the latest App version manifest; do not infer it
from Static Hosting listings or delete individual Hosting files.

Log recording is off by default and is never enabled by deploy or log queries. Inspect it with
`cloudbase_deployment(action="logging_status")`. Only after an explicit user request and consent to
cost, use `cloudbase_deployment(action="enable_logging", confirm=true)`. Stop new log recording with
`cloudbase_deployment(action="disable_logging", confirm=true)` after explicit confirmation. Query
recorded logs with exact `logs search ... --json` argv through `cloudbase_command`, not `fn log`;
some CloudBase ENVs reject the `GetFunctionLogs` API even with a current CLI.

Treat “undeploy the environment” as ambiguous:

- After explicit confirmation, `cloudbase_deployment(action="undeploy", confirm=true)` removes
  workloads and public routes but preserves the ENV, data, Auth, Storage, key, and reservations.
- Only after explicit confirmation of permanent data loss,
  `cloudbase_deployment(action="destroy_environment", confirm=true)` destroys the current project's
  ENV and its resources.

Never route either action through `cloudbase_command`. Repeated operations must reuse the current
project's persisted ENV mapping.
