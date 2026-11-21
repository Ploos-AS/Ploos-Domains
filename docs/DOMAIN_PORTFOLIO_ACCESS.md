# Organization-authorized portfolio reads

`PortfolioReader.ListForUser` accepts a **trusted internal user ID**, an organization ID and a PostgreSQL connection. The SQL checks for an organization-level `owner`, `admin` or `viewer` membership in the same statement as selecting portfolio records. Unit-only membership does not grant whole-organization access.

This is a backend primitive, **not an HTTP endpoint**. The authentication layer must derive the user ID from a validated session/OIDC identity, never from client-controlled request data. The empty result intentionally does not disclose whether the organization exists or the user lacks membership.

The low-level `PortfolioPostgresRepository` remains available for trusted internal processes only; do not wire it directly to customer-facing handlers. No write authorization, domain ownership verification or registrar sync is added here.

Before exposing an API, add PostgreSQL integration tests for owner/admin/viewer, unit-only membership, revoked membership, cross-organization isolation and invalid IDs; define response semantics and pagination. No production database or DNS changes in this PR.
