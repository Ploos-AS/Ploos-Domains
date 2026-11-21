# OIDC callback composition (unmounted)

The unmounted `NewOIDCCallbackHandler` requires a single state/code pair, browser-binding cookie, atomically consumed login transaction, PKCE code exchanger, signed ID-token verifier with nonce, trusted identity-to-user resolver and server-side session writer. It issues a secure host-only session cookie only after every step succeeds.

**Not production-ready:** the OAuth code exchanger implementation is not included, and the callback is not mounted. A trusted TLS reverse-proxy origin strategy is required for Caddy; the current direct-TLS guard rejects proxied HTTP. Redirect destination and cookie clearing/error behavior require integration review. End-to-end tests with a real OIDC provider and PostgreSQL are still needed. The caller must not supply identity, expected nonce or user ID from HTTP input. No production changes.
