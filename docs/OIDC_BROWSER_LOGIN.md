# OIDC browser login — transaction primitives

This PR provides cryptographically independent `state`, `nonce`, and PKCE verifier values and an HTTPS authorization-code URL builder using PKCE S256.

**Not a complete OIDC login:** no HTTP endpoints, server-side transaction storage, browser binding, one-time state validation, code exchange, nonce validation of a signed ID token, cookie issuance, session rotation, CSRF controls, or logout. The existing OIDC ID-token verifier must be extended or composed to validate nonce claims before issuing a browser session. Do not launch login or expose user data using these primitives alone.

Future implementation must use a trusted, fixed issuer and redirect URI, a short-lived single-use transaction, and a backend-managed session cookie. No tokens should be stored in browser localStorage. No production deployment.
