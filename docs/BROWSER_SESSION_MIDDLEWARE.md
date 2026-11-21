# Browser session middleware

`RequireBrowserSession` reads exactly one `__Host-ploos_session` cookie and resolves it using the server-side store on each request. The internal user ID is exposed only through `BrowserSessionUserID(ctx)`. Missing, duplicate, expired or revoked tokens fail closed; storage failures yield 503. No public route is mounted in this PR.

Integration test: run `go test -tags integration ./internal/server` with `PLOOS_TEST_DATABASE_URL` pointing to an **isolated, fully migrated, disposable PostgreSQL database**. Test creates and deletes its own user and sessions.

Still required before customer login: OIDC callback with PKCE/state/nonce validation, secure session issuance and rotation, logout, CSRF defenses, secure browser origin/TLS, and protected portal routing. Never treat a browser session user ID alone as organization authorization.
