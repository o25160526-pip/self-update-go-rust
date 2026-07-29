# CloudBase Data Access

Use this reference for relational Database, browser CRUD, RLS, browser-to-Function token handoff,
Storage, and trusted data operations. Load `references/auth-browser.md` for browser login and session
lifecycle, and `references/function-runtime.md` before implementing a Function.

## Contents

- Provisioning and credentials
- Choose the data path
- MySQL owner CRUD
- PostgreSQL owner CRUD
- Browser-to-HTTP-Function handoff
- Object Storage
- Trusted operations and acceptance

## Provisioning And Credentials

ENV creation alone does not provision MySQL or PostgreSQL. The baseline FlexDB resource shown by ENV
inspection is a document database and does not identify the relational engine.

Before any mutating Auth, Storage, or deployment operation creates a new ENV, choose PostgreSQL when
relational persistence is required and no explicit compatibility requirement selects MySQL:

```text
cloudbase_database(action="ensure_database", engine="postgresql")
```

Before schema or policy work, call
`cloudbase_database(action="inspect_database")` and preserve its authoritative engine. An existing
ENV cannot switch relational modes in place.

CloudBase owns all database accounts and passwords. Never create, alter, rename, or drop users or
roles; change passwords or membership; or modify credential catalogs. The database tool's `role`
selects an existing execution role only. PostgreSQL access uses SQL GRANT and RLS through
`cloudbase_database(action="set_access", sql=..., role="postgres")`. MySQL access uses its
table-policy parameters.

Neither browser-direct nor caller-scoped Function access uses Tencent AK/SK or database passwords:

| Path | Identity |
|---|---|
| Browser-direct | Publishable Key initializes the SDK; user session supplies caller identity |
| Function-mediated owner access | Browser sends its access token; Function forwards it unchanged |
| Trusted Function operation | Managed runtime identity; never a browser credential |

## Choose The Data Path

Prefer browser-direct access when Auth plus policy fully enforces ordinary owner CRUD. Add a
Function only for backend validation, orchestration, secrets, cross-user work, transactions, or a
required backend framework.

For browser-direct access:

- MySQL uses `app.mysql()` and `PRIVATE` table policy with `_openid`.
- PostgreSQL uses `app.rdb()`, SQL GRANT, and RLS with `auth.uid()`.
- Do not use document-database `collection()` APIs for either relational client.

For Function-mediated owner access, keep policy as the authorization boundary. The Function must
not replace the user token with admin authority, decode it, infer identity from token presence, or
accept an owner ID supplied by the browser.

## MySQL Owner CRUD

Include an indexed `_openid VARCHAR(64)` ownership column and apply:

```text
cloudbase_database(
  action="set_access",
  table_name="notes",
  permission="private"
)
```

CloudBase fills `_openid` from the authenticated session; browser inserts must not provide it.
MySQL table permission values do not apply to PostgreSQL.

## PostgreSQL Owner CRUD

Use `auth.uid()` for the CloudBase user ID. `current_user` is the database role name, not the
application user. Prefer database-owned identity defaults and omit `owner_id` from browser inserts.

A complete owner-only table policy includes schema/table/sequence grants and four operation
policies:

```sql
CREATE TABLE public.notes (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  content TEXT NOT NULL DEFAULT '',
  archived BOOLEAN NOT NULL DEFAULT false,
  owner_id VARCHAR(64) NOT NULL DEFAULT auth.uid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX notes_owner_id_idx ON public.notes(owner_id);

GRANT USAGE ON SCHEMA public TO authenticated;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.notes TO authenticated;
GRANT USAGE, SELECT ON SEQUENCE public.notes_id_seq TO authenticated;

ALTER TABLE public.notes ENABLE ROW LEVEL SECURITY;

CREATE POLICY notes_select_own ON public.notes
  FOR SELECT TO authenticated
  USING (owner_id = auth.uid());

CREATE POLICY notes_insert_own ON public.notes
  FOR INSERT TO authenticated
  WITH CHECK (owner_id = auth.uid());

CREATE POLICY notes_update_own ON public.notes
  FOR UPDATE TO authenticated
  USING (owner_id = auth.uid())
  WITH CHECK (owner_id = auth.uid());

CREATE POLICY notes_delete_own ON public.notes
  FOR DELETE TO authenticated
  USING (owner_id = auth.uid());
```

RLS with no policy denies access. A policy without GRANT still fails. Sequence-backed inserts need
sequence grants. UPDATE normally needs both `USING` and `WITH CHECK`.

Create the browser client only after schema and policy are ready:

```js
const db = app.rdb();

const result = await db
  .from("notes")
  .select("id,title,content,updated_at")
  .eq("archived", false)
  .order("updated_at", { ascending: false })
  .range(0, 49);
if (result.error) throw result.error;
```

Use `.from()`, `.select()`, `.insert()`, `.update().eq()`, `.delete().eq()`, `.eq()`, `.order()`,
and `.range()` for PostgreSQL. Do not substitute NoSQL or legacy builder methods such as
`.collection()`, `.where()`, or `.orderBy()`.

If the provider returns `FailedOperation.PGConnectError`, it failed before SQL evaluation. Retry
once; if the same code persists, stop and report the provider request ID instead of rewriting schema
or authentication policy.

## Browser-To-HTTP-Function Handoff

Use this path only when backend logic adds value. Obtain a current non-anonymous session immediately
before each protected request:

```js
async function callApi(path, options = {}) {
  const current = await auth.getSession();
  if (current.error) throw current.error;
  const session = current.data?.session;
  if (
    !session?.user ||
    session.user.is_anonymous ||
    session.user.isAnonymous ||
    !session.access_token
  ) {
    throw new Error("Authentication required");
  }

  return fetch(path, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${session.access_token}`,
    },
  });
}
```

The HTTP Function forwards the received `Authorization` header unchanged to the CloudBase
relational data API. It must not decode, persist, or log the token, derive a user ID from it, or
treat a non-empty header as verified authorization. Preserve downstream `401` and `403` semantics
without exposing provider internals.

On `401`, call `auth.refreshSession()` at most once, read the new session with `getSession()`, and
retry once. If refresh fails, clear protected state and require login. Never fall back to a
Publishable Key, anonymous identity, service user, Tencent credential, or database credential.

If Function code itself must know a provider-verified caller before data access, use an
authenticated Event Function instead. The FastAPI caller-token pattern and second-user isolation
test are exercised by `tools/cloudbase/examples/fastapi-function-api`.

## Object Storage

| Database mode | Management | Browser API |
|---|---|---|
| Integrated MySQL/default | No `bucket_id` | `app.uploadFile()` and `app.getTempFileURL()` |
| PostgreSQL Storage | Named bucket required | `app.storage.from(bucket).upload(key, file)` |

For PostgreSQL, create a private bucket with
`cloudbase_storage(action="ensure_bucket", ...)`. Configure `storage.objects` access with
PostgreSQL RLS through `cloudbase_database(action="set_access", sql=..., role="postgres")`;
`cloudbase_storage(action="set_access")` is not the PostgreSQL policy path. Scope policy to the
intended bucket and caller using `auth.uid()` and an owner field or user-specific key prefix.

Use signed URLs for private PostgreSQL objects. For integrated Storage, resolve the returned
`fileID` with `getTempFileURL()` instead of constructing a URL. Propagate upload failure before
writing Database metadata.

## Trusted Operations And Acceptance

Put cross-user reads, aggregation, moderation, secrets, and policy-inexpressible authorization in a
Function. Prefer an authenticated Event Function for provider-verified browser identity and read
caller information through the runtime SDK. An HTTP Function may remain caller-scoped by forwarding
the opaque access token, but the header is not locally verified identity.

Before reporting data access complete:

1. Test create, read, update, and delete through the actual user-facing path.
2. Verify signed-out requests cannot read or mutate protected data.
3. Verify a second user cannot access the first user's rows or objects.
4. Verify invalid or expired caller tokens fail and no token appears in logs.
5. Confirm browser bundles contain no Tencent credential, database password, service authority, or
   protected data.
