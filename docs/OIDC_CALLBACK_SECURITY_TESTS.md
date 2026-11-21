# OIDC callback failure-path tests

This PR adds isolated HTTP handler tests covering rejected state, failed code exchange, invalid ID token/nonce, unresolved identity, failed session persistence, duplicate callback parameters, unexpected parameters, and empty state/code. Tests assert that failure paths do not set a browser session cookie or redirect and that post-binding-validation failures clear the one-time binding cookie.

**Scope:** mock-backed handler tests only. They do not prove PostgreSQL atomicity, real OIDC token verification, concurrent callback replay protection, or complete login success behind Caddy. The PostgreSQL transaction lifecycle has separate integration-tagged tests. Full browser-to-OIDC-to-PostgreSQL integration testing, CSRF review and deployment gating are still required. No routes are mounted and no production changes are made.
