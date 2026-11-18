# Authenticated audit HTTP handler

`server.NewAuditHandler` provides a read-only JSON handler for `GET /api/v1/organizations/{organizationID}/audit-events`.

**Authentication is mandatory:** The handler requires an identity placed in request context by `auth.Middleware`, configured with a cryptographic OIDC verifier. It resolves the actor's internal user ID from the verified issuer+subject. It never trusts a user ID supplied by the client.

Authorization remains in `AuditQuery.List`: organization-scoped owner/admin only, tenant-filtered SQL. Supports `limit` (1–100, default 50) and opaque base64url JSON `cursor`. Responses set `Cache-Control: no-store`. No write routes.

**Deployment:** This PR does **not** register the handler in `server.NewHandler` or configure OIDC/database environment variables. It must be explicitly mounted with verified middleware and real database dependencies in a separate deployment PR. This prevents accidentally exposing a partially configured endpoint. Cursors are opaque but not cryptographically signed; they are validated and authorization is rechecked on each request.

Before public exposure: configure trusted OIDC, DB role separation, request rate limiting, request size/timeouts and tenant isolation tests.
