# Simulated Caddy proxy OIDC integration test

The integration-tagged test uses a disposable migrated PostgreSQL database and both unmounted OIDC HTTP handlers. It simulates Caddy's trusted loopback connection and forwarded HTTPS headers, exercises login start, verifies state/nonce/PKCE persistence, checks wrong-browser and untrusted-peer rejection, completes callback with test doubles for token exchange and ID-token verification, resolves the issued PostgreSQL session, and rejects callback replay.

**Limitations:** No actual Caddy process, OIDC provider, signed ID token, or production configuration is involved. This does not replace browser-based end-to-end testing, Caddy configuration review, or security review. No routes are mounted or deployed.
