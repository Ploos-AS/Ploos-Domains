# Authenticated domain portfolio read API

Route: `GET /api/v1/organizations/{organizationID}/portfolio`

Mount with `NewHandlerWithPortfolio`, passing a configured OIDC verifier, trusted issuer+subject-to-user resolver, and organization-authorized `PortfolioReader`. The endpoint is **not mounted** if any dependency is nil. An HTTP request must not supply its own internal user ID.

Response: `{"domains":[...]}` (empty array when no entries). The response has no-store headers; the route is rate-limited per process. Invalid tokens are rejected by the existing OIDC middleware. Unlinked identities return 403. Backend failure returns 503. A user without organization membership receives an empty list, intentionally not revealing whether an organization exists. This endpoint is read-only.

Important: current portfolio SQL returns the full result set; add pagination and query limits before use with large customer portfolios. Add end-to-end authenticated integration tests before production. No registrar calls, DNS changes, domain registration or deployment.
