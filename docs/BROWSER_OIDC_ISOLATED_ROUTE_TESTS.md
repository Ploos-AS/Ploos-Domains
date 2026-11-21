# Isolated browser OIDC route tests

PR #89 adds a **test-only** `http.ServeMux` in `internal/server/oidc_isolated_routes_test.go` that registers the prepared login-start and callback handlers under `POST /login/start` and `GET /login/callback`. This permits method, origin, host and failure-path checks through a router rather than calling the handlers directly.

The harness exists only in Go test code. `server.NewHandler` and `buildHandler` remain unchanged; no production router mounts these routes, and explicit OIDC opt-in continues to fail closed. Tests use mock dependencies and do not exercise a live identity provider or PostgreSQL. End-to-end integration, security review and an independently approved activation gate remain prerequisites for any real deployment.
