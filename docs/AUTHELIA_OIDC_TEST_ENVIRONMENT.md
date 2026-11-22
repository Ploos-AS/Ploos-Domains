# Authelia OIDC test environment — preparation checklist

This document plans a **separate nonproduction** identity-provider environment. It does not deploy Authelia, provision users, enable login routes, or change the Ploos Domains production configuration.

## Isolation and identity

- Run Authelia on a dedicated test host or isolated network segment, with its own HTTPS hostname and TLS certificate.
- Create a dedicated OIDC client for Ploos Domains, separate from all other Ploos applications. Register the exact HTTPS callback URI used by the test environment.
- Create a disposable test identity with no production access and no shared user database.
- Keep Authelia signing keys, session secrets and client credentials out of the repository. Protect files with restrictive permissions, maintain a rotation procedure, and avoid exposing credentials in logs.
- Restrict access to test administration and use a dedicated database and storage configuration. Document backup and teardown procedures.
- Review the deployed Authelia version's official OIDC configuration syntax before generating a runnable configuration; this checklist deliberately does not guess version-specific settings.

## Network and callback boundary

- Publish only the test provider and callback host over HTTPS; do not trust arbitrary forwarded headers.
- Set `PLOOS_TEST_LIVE_OIDC_ISSUER` to the exact public issuer from Discovery, without a trailing slash.
- Set `PLOOS_TEST_LIVE_OIDC_AUTHORIZATION_ENDPOINT`, `PLOOS_TEST_LIVE_OIDC_TOKEN_ENDPOINT`, `PLOOS_TEST_LIVE_OIDC_CLIENT_ID`, `PLOOS_TEST_LIVE_OIDC_REDIRECT_URI`, and `PLOOS_TEST_LIVE_OIDC_PUBLIC_HOST` to match the dedicated client and public HTTPS test host.
- Ensure the callback host is not a production Ploos Domains hostname. Avoid redirect aliases or wildcard callback URIs.
- Validate trusted reverse-proxy source IPs and exact forwarded host/proto values.

## Verification sequence

1. Run the existing live Discovery, verifier-initialization and JWKS smoke tests **manually** with `PLOOS_TEST_LIVE_OIDC_ENABLED=true` in a controlled environment.
2. Verify that the JWKS signing algorithm and key rotation work with the application verifier. The current smoke tests do not validate a signed ID token.
3. Build a separately reviewed nonproduction browser route activation gate; current `buildHandler` intentionally fails closed even with browser login opt-in.
4. Run a real authorization-code + PKCE browser flow using the disposable test identity, checking state, nonce, browser binding and secure session cookies.
5. Exercise expired/replayed state, mismatched nonce, bad signature, wrong issuer/audience, blocked redirects, session expiry/revocation and proxy spoofing.
6. Confirm monitoring, rate limiting, audit events, secret rotation and clean teardown before considering any broader activation.

## Safety boundary

Do not set `PLOOS_BROWSER_OIDC_LOGIN_ENABLED=true` on a production host as a workaround: the current startup gate intentionally refuses activation. No production route mounting is authorized by this document. A passing smoke test is **not** a completed end-to-end login test.
