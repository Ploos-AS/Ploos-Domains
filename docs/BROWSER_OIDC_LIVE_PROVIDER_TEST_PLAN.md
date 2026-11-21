# Live OIDC provider smoke test (manual, opt-in)

The integration-tagged test `TestBrowserOIDCLiveProviderDiscovery` checks **only live provider Discovery metadata** against explicitly supplied expected issuer, authorization endpoint, token endpoint and HTTPS JWKS URI. It does not authenticate a user, verify signed tokens, exchange authorization codes, access PostgreSQL or mount login routes.

## Requirements

- Use a dedicated nonproduction OIDC tenant/client and a trusted network.
- Configure the test variables `PLOOS_TEST_LIVE_OIDC_ISSUER`, `PLOOS_TEST_LIVE_OIDC_AUTHORIZATION_ENDPOINT`, `PLOOS_TEST_LIVE_OIDC_TOKEN_ENDPOINT`, `PLOOS_TEST_LIVE_OIDC_CLIENT_ID`, `PLOOS_TEST_LIVE_OIDC_REDIRECT_URI`, and `PLOOS_TEST_LIVE_OIDC_PUBLIC_HOST`.
- Set `PLOOS_TEST_LIVE_OIDC_ENABLED=true` **only** for an intentional manual test. By default the test skips.
- Run `go test -tags=integration -run '^TestBrowserOIDCLiveProviderDiscovery$' -count=1 ./internal/server`.
- No client secret is required or read. Do not place provider credentials in PR comments, shell history, CI variables on untrusted pull requests, or logs.

## Remaining gates

1. Establish a dedicated isolated browser-test environment and verify trusted proxy boundaries.
2. Exercise real authorization code + PKCE and signed ID-token verification with a disposable test identity.
3. Verify browser cookies, replay rejection, logout/revocation, error handling and audit/rate limiting.
4. Review operational key/secret rotation and production route activation separately.

This PR changes neither the normal CI workflow nor the application router. Browser login stays disabled.
