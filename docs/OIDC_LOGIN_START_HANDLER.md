# OIDC login-start handler (not mounted)

`NewOIDCLoginStartHandler` composes random state/nonce/PKCE, independent browser binding, the PostgreSQL login transaction writer, a secure host-only binding cookie, and a 303 redirect to a trusted HTTPS OIDC authorization endpoint. It accepts POST only, requires HTTPS and an exact same-origin Origin header, and fails closed if persistence fails.

**Not mounted or production-ready.** Direct TLS detection will need a trusted reverse-proxy strategy before using this handler behind Caddy; never trust arbitrary forwarded headers. The configured callback must be registered with the OIDC provider. Login callback, code exchange, nonce verification, identity mapping, session issuance, logout and end-to-end tests are still required. No production deployment.
