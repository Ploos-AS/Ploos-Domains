# Verification service boundary

This change introduces an application-level coordinator that connects three explicitly injected dependencies:

1. `DomainAuthorizer` — authenticates the actor elsewhere and checks their organization/domain permission.
2. `ChallengeRepository` — loads a challenge from trusted persistence and completes it transactionally.
3. `TXTResolver` — checks the TXT record at the expected DNS name.

The service refuses requests with missing identities or dependencies, denies cross-organization challenge mismatches, and never calls completion unless authorization and DNS proof both succeed.

## Important: not yet an HTTP API

This PR intentionally does **not** expose a public HTTP route, because authentication, organization scoping, database adapters, rate limiting, and end-to-end integration tests must be wired first. The service interfaces are the foundation for that API.

The existing `CompleteVerification` transaction helper must be used behind a repository adapter that loads challenge ownership from the database and rechecks it under locks. Do not trust organization IDs, domain IDs or challenge metadata supplied by clients. Avoid logging raw tokens.

Next: PostgreSQL repository adapter and integration tests; authenticated API handlers only after those checks are complete.
