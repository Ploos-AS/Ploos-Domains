# Browser OIDC least-privilege role CI check

After migrations 000008 and 000009, CI runs `db/security/test_browser_oidc_role.sql` against its disposable PostgreSQL database. The SQL creates a temporary NOLOGIN role inside a transaction, grants only the browser login's expected permissions, switches to that role, verifies required grants and absence of schema ownership/superuser/excessive grants, then revokes a required permission and confirms it is no longer available. The transaction is rolled back so no role or grants persist.

This checks PostgreSQL role semantics, **not** yet the Go `checkBrowserOIDCDatabase` function under a restricted connection. That integration test, deployment role provisioning and startup wiring remain future work. Browser login remains disabled; no production changes.
