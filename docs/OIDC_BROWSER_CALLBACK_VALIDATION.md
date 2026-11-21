# Browser OIDC ID-token validation

`OIDCVerifier.VerifyBrowserIDToken` uses the existing configured issuer/JWKS verifier to validate the signed ID token, issuer, audience, expiry and allowed algorithms. It additionally rejects missing or mismatched nonce claims using constant-time comparison and requires a nonempty subject.

The expected nonce must come **only** from an atomically consumed, browser-bound `BrowserLoginStore` transaction, never a query parameter or arbitrary client input.

This PR contains a verification primitive and unit tests, **not** a working callback. Remaining steps: validate the callback's state and binding cookie, exchange code with PKCE against trusted token endpoint, verify ID token with the consumed nonce, map the OIDC identity to a trusted internal user, create and set the server-side session cookie, clear transaction cookie, handle errors safely, and add end-to-end integration tests. No production endpoint or deployment.
