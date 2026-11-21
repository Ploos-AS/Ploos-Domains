# Browser OIDC startup discovery gate

On explicit `PLOOS_BROWSER_OIDC_LOGIN_ENABLED=true`, startup validates the operator-controlled OIDC configuration, then uses `OIDCDiscoveryClient` to retrieve fixed-issuer metadata with normal TLS verification, a 10-second context timeout, response size limits and redirect rejection. Discovery must match the configured issuer, authorization and token endpoints before PostgreSQL readiness is checked. Failure aborts startup without mounting any login routes.

Tests use a deterministic validator stub for successful and mismatched discovery. Live provider testing remains a separate deployment gate. Even if discovery and database checks succeed, startup still fails closed with `errBrowserLoginNotConfigured` until protected client secret handling and complete handler wiring are ready. Disabled-by-default startup remains unchanged; no production deployment.
