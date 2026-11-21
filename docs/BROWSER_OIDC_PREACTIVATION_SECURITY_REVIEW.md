# Browser OIDC pre-activation security review

## Verified in source (not a production penetration test)

- `buildHandler` keeps browser login disabled by default, rejects invalid feature flag values, and fails closed on explicit opt-in.
- The browser login start and callback handlers are prepared but not registered in `server.NewHandler`.
- Login start requires POST, HTTPS and exact same-origin `Origin`; callback requires GET, HTTPS, single `state` and `code`, and a browser-bound one-time transaction.
- PKCE, signed ID-token verification and nonce checks are implemented in separate components.
- Browser session tokens are random and stored only as digests in PostgreSQL; cookies use Secure, HttpOnly and `__Host-` conventions.
- Startup checks provider Discovery, protected client secret file, database schema and least-privilege role before preparing runtime.

## Required before enabling browser login

1. Confirm live OIDC provider configuration, callback URI registration, signing keys and trusted reverse-proxy IPs in a dedicated test environment.
2. Review CSRF, redirect behavior, proxy header handling, cookie scope and session expiration in an end-to-end browser test.
3. Exercise actual PostgreSQL role privileges and migrations, including failed-login paths and transaction replay.
4. Verify clean runtime ownership and shutdown while HTTP requests are active.
5. Add explicit route mounting behind a separately reviewed activation gate; do not reuse this PR as authorization to enable login.
6. Conduct security review of the complete flow, including audit logging, rate limiting and operational secret rotation.

This review is a source-level checklist. It does not claim that all requirements are satisfied or that browser login is production-ready.
