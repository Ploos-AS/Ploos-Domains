# OIDC login transaction storage

Migration 000009 adds a PostgreSQL table for short-lived OIDC login transactions. `BrowserLoginStore.Create` stores SHA-256 digests of state and browser-binding secrets, and the nonce and PKCE verifier needed by the callback. `Consume` atomically marks a matching, unexpired transaction used and returns its nonce and verifier once.

A future login-start handler must create a random independent browser-binding secret and set a Secure/HttpOnly/SameSite cookie; the callback must require the matching cookie, consume the transaction, clear the binding cookie, exchange the authorization code using PKCE, and validate signed ID-token nonce/issuer/audience before issuing a browser session. No such endpoints exist yet.

The PKCE verifier and nonce are sensitive; protect the database and backups and limit application DB privileges. Add cleanup of expired/consumed rows. Integration test uses a disposable migrated PostgreSQL database via `PLOOS_TEST_DATABASE_URL`. No production migrations run here.
