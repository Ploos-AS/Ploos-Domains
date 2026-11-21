# Browser OIDC confidential code exchange configuration

`loadBrowserOIDCCodeExchanger` combines a validated `OIDCDeploymentConfig` with a client secret loaded from `PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE`. It constructs the existing `server.OIDCCodeExchange` using the fixed token endpoint, client ID and redirect URI. The secret is never printed in error messages.

This helper is **not yet invoked by buildHandler** and does not make an HTTP request, mount routes, or enable browser login. The next integration step must bind the exchanger to callback dependencies, with verified OIDC token handling and dedicated runtime database connections. Keep browser login disabled until all components are wired and tested.
