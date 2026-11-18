# Database-backed domain verification authorization

`PostgresAuthorizer` implements the `DomainAuthorizer` interface used by the verification service.

Only **organization-scoped** `owner` and `admin` memberships may verify a domain. A unit-scoped administrator, a viewer, a member of another organization, or an unknown user is denied. The SQL query joins the target domain and the membership through the same organization ID.

## Trust boundary

`actorID` must be obtained from a verified authentication session and resolved to an internal user UUID. Never accept actor identity from a client-controlled header or JSON field.

This change does not add HTTP endpoints, authentication middleware, rate limiting or audit events. The service still needs a properly configured database connection and end-to-end tests before it can be exposed publicly.

Integration tests use a disposable PostgreSQL database with migrations 000001–000003 applied.
