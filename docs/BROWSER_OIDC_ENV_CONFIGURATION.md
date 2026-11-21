# Browser OIDC environment configuration

This PR adds a strict, operator-controlled environment loader for the future browser login startup wiring. It is not invoked by the running server yet, and `PLOOS_BROWSER_OIDC_LOGIN_ENABLED=true` **still fails closed**.

Required keys (all prefixed with `PLOOS_BROWSER_OIDC_`): `ISSUER`, `AUTHORIZATION_ENDPOINT`, `TOKEN_ENDPOINT`, `CLIENT_ID`, `REDIRECT_URI`, `PUBLIC_HOST`, `TRUSTED_PROXY_IPS` (comma-separated IP addresses), `LOGIN_TTL` and `SESSION_TTL` (Go duration strings).

The loader validates syntax, lifetimes and trusted peer addresses using `OIDCDeploymentConfig.Validate`. Do not put a client secret in these variables or commit it to the repository. Discovery matching, protected secret loading, PostgreSQL checks, handler construction and router registration remain separate steps. No routes or production settings are changed.
