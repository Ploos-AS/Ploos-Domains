# Browser OIDC Go readiness check with restricted PostgreSQL role

In disposable PostgreSQL CI, create the dedicated login role using `db/security/setup_browser_oidc_ci_role.sql` after migrations 000008–000009. Then run `go test -tags=integration ./cmd/ploos-domains` with both `PLOOS_TEST_BROWSER_OIDC_DATABASE_URL` (restricted role) and `PLOOS_TEST_DATABASE_URL` (CI database owner).

The test calls the actual Go `checkBrowserOIDCDatabase` function and verifies that the restricted role is accepted. It then revokes the required `INSERT` permission on `browser_sessions` and confirms the Go check fails closed. The fixed CI password is **test-only** and must never be used in deployment. This role is created only in ephemeral CI PostgreSQL; no production database is modified.

Browser login is still disabled and startup wiring is not changed.
