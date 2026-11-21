# Browser binding cookie for OIDC login

This PR adds an independent random browser-binding secret, issued as a host-only `__Host-` cookie with Secure, HttpOnly, SameSite=Lax, Path=/ and a 1–10 minute lifetime. Callback parsing rejects absent or duplicate cookies. The cookie must be cleared after callback handling.

**No login or callback endpoints are mounted.** The future login-start handler must generate the binding, store its SHA-256 digest with the state transaction, then redirect. The callback must require a matching binding cookie and atomically consume state before code exchange, validate PKCE and signed ID-token nonce, resolve the internal user and issue a fresh session. HTTPS, CSRF review, and integration tests are mandatory before production.
