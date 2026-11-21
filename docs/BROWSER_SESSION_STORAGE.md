# PostgreSQL browser session storage

Migration `000008` creates `browser_sessions` with a SHA-256 token digest primary key, internal user foreign key, expiration and revocation timestamps. A session store exposes Create, Resolve and Revoke primitives. Resolve checks revocation and expiration using PostgreSQL time; Revoke targets one token only.

**Not an authentication flow.** A trusted OIDC callback must validate state, nonce, PKCE, issuer, audience and token signature before creating a session. Callers must never accept an arbitrary user ID from the browser. Raw tokens must not be logged or stored. Session cookies must only be set over HTTPS.

Remaining: session middleware, CSRF protection, login/logout endpoints, session rotation, cleanup policy, integration tests with migrated PostgreSQL and UI wiring. Production rollout requires security review. No production migrations executed in this PR.
