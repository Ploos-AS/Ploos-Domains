# Manual OIDC JWKS readiness smoke test

`TestBrowserOIDCLiveJWKSReadiness` checks the public JWKS endpoint advertised by a deliberately configured test OIDC provider. It requires the same explicit `PLOOS_TEST_LIVE_OIDC_ENABLED=true` gate as the existing live-provider smoke tests.

Set `PLOOS_TEST_LIVE_OIDC_ISSUER` to the dedicated test provider's HTTPS root issuer, then run:

```sh
go test -tags=integration -run '^TestBrowserOIDCLiveJWKSReadiness$' -count=1 ./internal/server
```

The test fetches public Discovery metadata and JWKS, refuses HTTP redirects, bounds response sizes, and checks for at least one advertised signing key. It uses no client secret, browser session, PostgreSQL connection or authorization code. **It does not verify a signature or establish end-to-end login readiness.** Production route mounting remains disabled. Never use untrusted issuer URLs: the explicit test gate intentionally permits outbound network requests to the configured provider.
