# OIDC token verifier

PR #13 adds a concrete OIDC **ID-token** verifier using `github.com/coreos/go-oidc/v3/oidc`.

Configure a fixed HTTPS issuer and expected client ID (audience). The library discovers the provider metadata and JWKS, validates signatures using the configured signing algorithms (RS256 or ES256), validates issuer, audience and expiry, and handles remote key retrieval and refresh.

`NewOIDCVerifier` performs network discovery at startup; use a bounded startup context and fail startup if discovery fails. Do not use token-provided issuer URLs, disable TLS checks, or bypass signature verification.

## Important distinction: ID tokens versus API access tokens

This implementation validates **OIDC ID tokens**, which are normally intended for the client application, **not arbitrary OAuth2 access tokens**. Do not use ID tokens as the final public API bearer-token design without reviewing the authentication flow. For a public resource-server API, validate provider-issued **access tokens** with an appropriate resource audience and token-type checks, or use a secure session/BFF design.

## Limitations

- No HTTP routes are wired to this verifier yet.
- No issuer+subject mapping to internal users yet.
- No login redirect/callback, nonce, PKCE or session management.
- No revocation, audit or rate limiting.
- No live-provider integration test; current tests cover configuration and fail-closed behavior.
- The Go dependency must be resolved and `go.sum` committed before production builds.

The design remains provider-neutral; compatibility with Authelia and Keycloak must be verified against their actual provider configurations.
