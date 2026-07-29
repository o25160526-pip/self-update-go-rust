# CloudBase Browser Auth

Use this reference for browser initialization, login methods, sessions, route guards, and Auth
acceptance. Load `references/data-access.md` as well when authenticated users access Database,
Storage, or an HTTP Function.

## Contents

- Credential model
- Configure Auth
- Email and Phone OTP
- Username and password
- Session and UI lifecycle
- Acceptance

## Credential Model

| Value | Purpose | Allowed exposure |
|---|---|---|
| ENV ID | Selects the project ENV | Browser-safe |
| Publishable Key (`publish_key`) | Initializes the browser data plane | Browser-safe |
| User access token | Identifies a signed-in user | Current browser session only |
| Managed runtime identity | Gives a Function trusted resource access | Function runtime only |
| Tencent AK/SK | Manages provider resources | Agent Server only |

A Publishable Key identifies the ENV client; it is not a user credential and grants no trusted
authority by itself. Auth session plus Database or Storage policy determines access. Obtain the
exact ENV, region, key, and initialization contract only from:

```text
cloudbase_auth(action="frontend_config")
```

CloudBase remains the key's source of truth. Do not reconstruct it, persist another copy, log it, or
use it as an HTTP Function Bearer token. Never put Tencent credentials in browser or Function code.

Use an `@cloudbase/js-sdk` version satisfying `>=3.6.3 <4.0.0`, the compatibility range validated
for this managed environment. Version `4.0.0` has an incompatible Auth API for these flows. Treat a
version change as a deliberate compatibility test, not an automatic upgrade.

In Next.js, initialize the Web SDK in a Client Component or client-only module, never in a Server
Component or Route Handler. For static export, preserve `output: "export"` and the deployable `out`
directory while making the smallest client-boundary or bundler correction.

## Configure Auth

Choose the account model before writing UI:

- Open registration may use `shouldCreateUser: true` when policy isolates each user's data.
- Owner-only access must disable self-registration or enforce an explicit owner allowlist in
  trusted backend logic. Successful registration alone does not establish ownership.

Then:

1. Inspect current state with `cloudbase_auth(action="inspect")`.
2. Enable only required methods with `cloudbase_auth(action="configure_login", ...)`.
3. Disable anonymous login for user-facing private applications.
4. Get browser configuration with `frontend_config`; use the returned initialization code.
5. Add only explicit local development origins with `allow_origins`. Managed deployment registers
   the complete managed origin returned by the domain/deployment tools.

Use `email_otp=true` for Email OTP, `phone_otp=true` for Phone OTP, and
`username_password=true` for username/password login. Phone OTP is available only in Shanghai-region
ENVs. Do not write a login flow until `inspect` confirms its provider is enabled.

## Email And Phone OTP

CloudBase returns the OTP verifier from the initial call. There is no standalone
`auth.verifyOtp()` step:

```js
const sent = await auth.signInWithOtp({
  email,
  options: { shouldCreateUser: true },
});
if (sent.error) throw sent.error;

const verified = await sent.data.verifyOtp({ token: code });
if (verified.error) throw verified.error;

const current = await auth.getSession();
if (current.error) throw current.error;
const session = current.data?.session;
if (!session?.user || session.user.is_anonymous || session.user.isAnonymous) {
  throw new Error("Authentication required");
}
```

Phone OTP uses the same returned-verifier flow:

```js
const sent = await auth.signInWithOtp({
  phone: "+8613800138000",
  options: { shouldCreateUser: true },
});
const verified = await sent.data.verifyOtp({ token: code });
```

Keep the country code in phone numbers. Preserve the verifier until the user submits the code; do
not replace it with a guessed endpoint or a different SDK's OTP API.

## Username And Password

After confirming `username_password` is enabled, login with:

```js
const result = await auth.signInWithPassword({ username, password });
if (result.error) throw result.error;
```

A plain identifier such as `admin` or `editor` is a username, not an email; keep its form input as
`type="text"`. Direct browser `auth.signUp({ username, password })` is conditional on the installed
SDK and provider capability. Verify support before implementing it. If unsupported, report the
capability gap or use an available verified registration provider; never introduce a management
credential or secret key into browser code.

## Session And UI Lifecycle

Use `auth.getSession()` for route guards. Do not use deprecated `getLoginState()`, `getUser()`, a
UID-like value, or possession of the Publishable Key as proof of login.

The required lifecycle is:

1. Restore state with `getSession()` at application startup.
2. Render no protected data until a real, non-anonymous session exists.
3. Complete login through the configured provider.
4. Load protected data only after authentication.
5. On `signOut()`, clear the session-dependent cache and rendered protected state.

A client route guard controls rendering, not authorization. Never compile private records, secrets,
or serialized protected payloads into static HTML or public JavaScript. Resource policy or trusted
backend logic remains the authorization boundary.

When calling a protected HTTP Function, read a fresh `session.access_token` immediately before the
request. Follow the forwarding and refresh contract in `references/data-access.md`. Never copy
access or refresh tokens into application storage, source, build configuration, logs, cookies
managed outside the SDK, or a backend database.

## Acceptance

Before reporting Auth complete:

1. Inspect built HTML and JavaScript for protected business data and secrets.
2. Open the managed URL returned by deployment in a fresh private session.
3. Confirm protected content is absent before login.
4. Complete the configured login and one real protected read or write.
5. Sign out and confirm protected UI and cached state disappear.
6. For per-user data, verify a second user cannot access the first user's records or objects.

The browser lifecycle in `tools/cloudbase/examples/fullstack-notes` is a validated reference for
Email OTP, session restoration, signed-out rendering, and owner-scoped CRUD.
