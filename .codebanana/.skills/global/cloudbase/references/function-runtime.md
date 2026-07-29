# CloudBase Function Runtime

Use this reference before writing, deploying, or debugging any CloudBase Function. Keep source and
configuration in the workspace; perform provider operations only through managed `cloudbase_*`
tools.

## Contents

- Choose the Function model
- Deployment configuration
- Event Functions
- HTTP Functions
- Data identity and credentials
- Routing and authorization
- Test and diagnose

## Choose The Function Model

Do not add a Function when browser Auth plus policy safely implements owner CRUD or private file
access.

| Requirement | Type | Runtime contract |
|---|---|---|
| Authenticated SDK operation, provider-verified caller, event, or background handler | Event | `event, context`; no listening port |
| Public REST, GraphQL, SSR, framework API, SSE, or WebSocket | HTTP | Web server on `0.0.0.0:9000` |

Event Functions process an invocation handler. HTTP Functions use native HTTP semantics and may
process concurrent requests in one instance. HTTP gateway exposure does not convert an Event
Function handler into an HTTP Function.

## Deployment Configuration

Declare Functions in the workspace-root `cloudbaserc.json` and deploy all declared workloads with
`cloudbase_deployment(action="deploy")`:

```json
{
  "$schema": "https://static.cloudbase.net/cli/cloudbaserc.schema.json",
  "version": "2.0",
  "functionRoot": "functions",
  "functions": [
    {
      "name": "api",
      "type": "HTTP",
      "runtime": "Nodejs20.19",
      "timeout": 60,
      "memorySize": 512,
      "envVariables": {}
    },
    {
      "name": "worker",
      "type": "Event",
      "runtime": "Nodejs20.19",
      "handler": "index.main",
      "timeout": 60,
      "memorySize": 256,
      "envVariables": {}
    }
  ]
}
```

- Never hardcode `envId`, `TCB_ENV`, region, or a provider credential. The server injects the
  authoritative ENV into deployment and runtime.
- Select an explicitly supported Function type and runtime before first deployment. Changing either
  on an existing Function requires deliberate delete and recreation; do not infer it from a retry.
- A configured `dir` is workspace-relative. Without it, source is
  `<functionRoot>/<function-name>`. Source must remain inside the workspace.
- `envVariables` is replaced as a complete set, not incrementally merged. Keep it intentional and
  non-secret.
- Package dependencies for the provider's Linux runtime. Native macOS wheels and binaries are not
  portable.
- Validate package names, versions, lockfiles, and local builds before deployment. A provider
  dependency error does not prove `node_modules` is missing.

For generated Node.js Functions that need the server SDK, use the managed tested exact dependency:

```json
{
  "dependencies": {
    "@cloudbase/node-sdk": "4.0.3"
  }
}
```

Do not invent a newer major or use an unverified range. Verify any deliberate upgrade against the
registry and a real local install/build first.

## Event Functions

An Event Function exports a handler and never starts a server:

```js
exports.main = async (event, context) => {
  return { ok: true, input: event };
};
```

Keep `handler` aligned with the source export, validate event data, return JSON-serializable values,
and make retryable side effects idempotent.

For an authenticated browser operation, invoke through the initialized Web SDK:

```js
const result = await app.callFunction({
  name: "update-profile",
  data: { displayName },
});
```

Read the provider-verified caller from the runtime instead of accepting a user ID:

```js
const tcb = require("@cloudbase/node-sdk");
const app = tcb.init({ env: tcb.SYMBOL_DEFAULT_ENV });

exports.main = async (event) => {
  const { uid } = app.auth().getUserInfo();
  if (!uid) throw new Error("Authentication required");
  return { uid };
};
```

Require a signed-in, non-anonymous caller in the ENV-wide Function policy and still enforce
application ownership and roles. Management-authorized invocation through `cloudbase_command` does
not prove browser authorization.

## HTTP Functions

An HTTP Function must:

- Include executable `scf_bootstrap` at the Function source root with LF endings, a valid shebang,
  and the runtime's absolute executable.
- Listen on `0.0.0.0:9000`, never only localhost or another port.
- Write temporary data only under `/tmp`; local source storage is not durable.
- Package Linux-compatible runtime dependencies and build framework output before packaging.
- Provide a cheap health endpoint, bounded request bodies, explicit status codes, and structured
  error handling.
- Initialize reusable clients and pools outside request handlers, but keep mutable per-request and
  per-user state local because requests are concurrent.

Example Python 3.10 bootstrap:

```bash
#!/bin/bash
export PYTHONPATH="${PWD}/third_party:${PWD}:${PYTHONPATH}"
exec /var/lang/python310/bin/python3.10 -m uvicorn app:app --host 0.0.0.0 --port 9000
```

With native Node.js HTTP, choose CommonJS or ESM consistently, parse `req.url`, and bound the request
stream before parsing JSON. Native requests do not provide Express helpers.

For Next.js server runtime, package production standalone output and start that packaged server;
the Function platform does not build an arbitrary source tree. For SSE, return
`text/event-stream`, disable caching, stop on disconnect, and account for timeout and concurrency.
For WebSocket, declare `protocolType: "WS"` with an idle timeout and verify a real secure
connection through the managed URL before claiming support.

Prefer same-origin frontend and API routes. When separate origins are required, handle `OPTIONS` and
allow only known origins, methods, and headers. Never combine credentialed requests with wildcard
origin.

## Data Identity And Credentials

Choose identity before choosing a client:

- Caller-scoped relational CRUD: receive the signed-in user's access token and forward the
  `Authorization` header unchanged to the CloudBase relational data API. Database policy remains
  the authorization boundary.
- Provider-verified user logic: use an authenticated Event Function.
- Trusted cross-user work: use managed runtime identity and enforce application authorization.

Never place Tencent AK/SK, `DATABASE_URL`, database username/password, Publishable Key, or service
user credentials in Function source or `cloudbaserc.json`. Do not sign in a service user from
Function code.

Node.js runtime identity can initialize its own ENV without static credentials:

```js
const tcb = require("@cloudbase/node-sdk");
const app = tcb.init({ env: tcb.SYMBOL_DEFAULT_ENV });
```

Inspect the project engine before selecting a data client:

- Document database: `app.database()`
- CloudBase MySQL: `app.mysql()`
- CloudBase PostgreSQL: `app.rdb()`

Do not mix these APIs. Reuse initialized clients across invocations and apply ownership checks
before trusted operations.

Python, Java, Go, and other runtimes may use a documented CloudBase HTTP API. A caller-scoped HTTP
Function reads `TCB_ENV`, requires the incoming Bearer token, and forwards it without parsing:

```python
env_id = os.environ["TCB_ENV"]
authorization = request.headers.get("authorization")
response = await client.get(
    f"https://{env_id}.api.tcloudbasegateway.com/v1/rdb/rest/notes",
    headers={"Authorization": authorization},
)
```

Never guess a runtime token shape or introduce a database password when an identity path cannot be
verified. Load `references/data-access.md` for the complete handoff and policy contract.

## Routing And Authorization

Managed defaults are:

- App only: `/` targets the App.
- One HTTP Function without an App: `/` targets the Function.
- App plus HTTP Functions: `/` targets the App and each Function uses `/<function-name>`.
- Event Function: no public route.

Omit deployment `routes` for these defaults. Multiple HTTP Functions without an App require a
complete explicit map with one `/` route.

Use only the complete route and URL returned by managed tools. Never append a server-configured
domain suffix or construct a provider URL.

Managed custom-domain HTTP routes do not provide CloudBase identity authentication at the Function
gateway. Treat them as public. Token forwarding leaves authorization to the data API; a Function
must not hand-parse the token or trust its presence. Use an Event Function when Function code needs
provider-verified caller identity.

Inspect the ENV-wide Function policy with:

```text
cloudbase_command(args=["policy", "get", "--json"])
```

A policy update affects the entire ENV and requires explicit user confirmation.

## Test And Diagnose

1. Run local tests and start an HTTP Function on port `9000` for representative API requests.
2. Reserve a slug before infrastructure only when an App or HTTP Function creates a public route.
3. Deploy with `cloudbase_deployment(action="deploy")`.
4. Inspect reconciliation with `cloudbase_deployment(action="status")`.
5. Test health and one real operation through every returned public URL.
6. Invoke an Event Function with
   `cloudbase_command(args=["fn", "invoke", "worker", "--params", "{\"key\":\"value\"}", "--json"])`.
7. For authenticated Event behavior, also invoke from a real signed-in Web SDK session and verify
   signed-out or anonymous callers are denied.
8. Inspect with `cloudbase_command(args=["fn", "list", "--json"])` or
   `cloudbase_command(args=["fn", "detail", "worker", "--json"])`.

Log recording is disabled by default and deployment never enables it. First call
`cloudbase_deployment(action="logging_status")`. Only after the user explicitly requests logging
and accepts its cost, call `cloudbase_deployment(action="enable_logging", confirm=true)`. Search
newest records with exact TCB argv:

```text
cloudbase_command(args=["logs", "search", "-q", "function_name:\"worker\"", "-t", "1h", "--sort", "desc", "--limit", "20", "--json"])
```

Use `logs search` for Function log content. `--sort desc --limit N` returns the newest `N` matching
records. Do not choose `fn log` by default: it calls `GetFunctionLogs`, which some CloudBase ENVs
reject with an `InvalidParameter` error claiming the developer tools must be upgraded, even when
the CLI is current. If that error occurs, retry the equivalent `logs search` query; do not infer
expired authentication, request browser login, or keep upgrading the CLI. Stop new recording only
after explicit confirmation with
`cloudbase_deployment(action="disable_logging", confirm=true)`; historical logs remain searchable.

`cloudbase_command` is 100% TCB argv passthrough within its allowlist. Pass only exact tokens after
`tcb`; the server adds only authoritative ENV, config, and credentials. Commands run from an
isolated authentication directory, so local file arguments must be absolute workspace paths. Do
not pass `tcb`, ENV ID, region, or credentials. Never invent command actions or use it for Function
deployment, deletion, code updates, or logging lifecycle.

Diagnose failures from source layout, runtime, bootstrap, dependencies, package build, configuration,
provider state, request IDs, and bounded logs. Never fall back to CloudBase Run or Custom Image HTTP
Function deployment; neither is supported by the current managed tools.
