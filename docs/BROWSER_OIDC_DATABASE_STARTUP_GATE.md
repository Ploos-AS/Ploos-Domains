# Browser OIDC database startup gate

When `PLOOS_BROWSER_OIDC_LOGIN_ENABLED=true`, startup validates the browser OIDC configuration, then requires a separate `PLOOS_BROWSER_OIDC_DATABASE_URL`. The application opens a bounded PostgreSQL connection, checks connectivity and runs `checkBrowserOIDCDatabase` against the configured runtime role. The connection is closed on success and failure. The check rejects missing privileges, superusers and table owners.

**Important:** even after a successful database check, startup still fails with `errBrowserLoginNotConfigured`. Login routes are not mounted. Provider discovery, secret handling and complete wiring remain prerequisites. When login is disabled, no browser OIDC database connection is attempted. No production deployment.
