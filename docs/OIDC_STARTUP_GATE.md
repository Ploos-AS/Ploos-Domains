# OIDC startup validation gate

`PrepareOIDCLoginHandlers` requires structurally valid fixed configuration, non-nil handler dependencies, consistent session lifetime, and successful provider discovery verification **before** returning any login handlers. Discovery errors fail closed. The returned handlers remain unmounted until a future reviewed router integration explicitly registers them.

At startup, pass `OIDCDiscoveryClient{}` with normal TLS certificate verification, supply PostgreSQL-backed transaction/session stores, and wire a real signed ID-token verifier and confidential code exchanger. Do not use a test discovery stub outside tests.

**Remaining work:** router startup wiring and real provider integration, secure secret management, Caddy forwarding configuration, CSRF review and logout. This PR does not activate or deploy login.
