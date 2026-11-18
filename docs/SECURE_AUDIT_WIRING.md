# Secure audit API wiring

The audit endpoint is **disabled by default**. Enable it only after deploying migrations, a trusted OIDC provider and a dedicated least-privilege PostgreSQL application login.

Environment variables:

- `PLOOS_AUDIT_API_ENABLED=true` — explicit opt-in; any other value leaves the audit route unmounted.
- `PLOOS_OIDC_ISSUER` — configured HTTPS OIDC issuer (never taken from incoming requests).
- `PLOOS_OIDC_AUDIENCE` — the configured client ID/audience for **OIDC ID tokens**.
- `PLOOS_DATABASE_URL` — PostgreSQL DSN, supplied through deployment secrets; use TLS in production.

Startup fails if audit is enabled but configuration, OIDC discovery, or database connectivity fails. The handler is mounted with cryptographic bearer verification, internal user resolution and organization-scoped owner/admin authorization.

**Important:** `auth.OIDCVerifier` validates **OIDC ID tokens**, not general OAuth2 access tokens. Do not pass arbitrary access tokens. Confirm that this token model is suitable for the client and deployment before exposing the route publicly.

**Database privileges:** Do not connect using the schema owner or migration user. The application login needs SELECT on `user_oidc_identities`, `memberships`, and `audit_events`. Audit SELECT must be exposed only through the trusted API, never direct tenant database access.

Remaining production work: rate limiting, reverse-proxy TLS, secure secret injection, observability, credential rotation, strict audit read role configuration and end-to-end deployment tests.
