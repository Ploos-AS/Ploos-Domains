# Browser OIDC PostgreSQL callback dependencies

`prepareBrowserOIDCPostgresLoginStores` constructs the browser login transaction writer/consumer, linked-identity resolver, and browser session writer from one supplied PostgreSQL runtime connection. It requires an independently configured signed-ID-token verifier and deliberately leaves code exchange unset until `prepareBrowserOIDCCallbackDependencies` adds the protected-secret exchanger.

The caller must first check PostgreSQL readiness and retain the open restricted-role connection for the lifetime of handlers. This helper does not open a connection, verify migrations, perform requests, mount routes or enable login. Startup still fails closed. No production deployment.
