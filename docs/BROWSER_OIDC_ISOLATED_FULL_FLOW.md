# Isolated browser OIDC full-flow test

The test-only router exercises POST login start, authorization state and nonce, browser binding cookie, GET callback, PKCE verifier handoff, nonce verification, session issuance and one-time transaction replay rejection. It also checks that missing browser binding does not consume a transaction.

All identity-provider, transaction and session dependencies are **in-memory fakes**. This is a contract regression test, not proof of cryptographic token verification or PostgreSQL atomicity; those need their own tests. No production route mounting, feature activation or deployment is included.
