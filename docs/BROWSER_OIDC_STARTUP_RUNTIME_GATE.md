# Browser OIDC startup runtime gate

When `PLOOS_BROWSER_OIDC_LOGIN_ENABLED=true`, startup first validates the configured issuer, protected client secret, provider Discovery and dedicated PostgreSQL role. After those checks pass, it invokes `initializeBrowserOIDCRuntime` to construct the signed-ID-token verifier, database-backed stores and **unmounted** login handlers. It then closes the runtime's database pool and returns `errBrowserLoginNotConfigured`.

This intentionally repeats some readiness checks during initial rollout. No handlers are registered, no browser login is enabled, and no production deployment is performed. Failures are fail-closed; the disabled/default path is unchanged. A later change should consolidate checks and explicitly review route mounting and lifecycle ownership before enabling login.
