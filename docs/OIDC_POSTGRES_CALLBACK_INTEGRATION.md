# PostgreSQL-backed OIDC callback integration

The `integration`-tagged callback test uses a disposable, fully migrated PostgreSQL database via `PLOOS_TEST_DATABASE_URL`. It creates a real one-time login transaction, completes the unmounted HTTP callback with controlled token-exchange, ID-token-verification and identity-resolution doubles, verifies the resulting persisted session, rejects callback replay, and revokes the session. No production database is used.

This is **not** a full external OIDC end-to-end test: it does not contact an OIDC provider, validate a real signed token, or exercise Caddy. Provider-backed and proxy-backed end-to-end tests are required before enabling login in production.
