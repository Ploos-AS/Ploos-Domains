# Browser OIDC PostgreSQL CI coverage

The PostgreSQL CI job previously ran integration-tagged Go tests only in `internal/verification` and `internal/auth`. This unintentionally omitted the OIDC callback, proxy and test-only router PostgreSQL integration tests under `internal/server`, including the test added in PR #91.

The job now includes `./internal/server` with the existing migrated PostgreSQL service and `PLOOS_TEST_DATABASE_URL`. The Go unit-test job is unchanged. This change does not activate browser login or connect to a live OIDC provider.

Before any live-provider testing, use a dedicated test tenant and test client credentials; do not put secrets in CI logs or pull-request workflows from untrusted contributors.
