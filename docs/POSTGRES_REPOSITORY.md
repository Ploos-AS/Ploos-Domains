# PostgreSQL verification repository

`PostgresRepository` implements the verification service's `ChallengeRepository` interface using standard `database/sql`. It loads challenge details from a join against the owning domain, including the trusted organization ID and SHA-256 digest, and delegates completion to the atomic transaction helper.

## Integration boundary

This PR intentionally does not add an HTTP route, driver dependency, authentication or database connection configuration. Those are required to run the adapter. The SQL smoke test validates persistence structure, not the Go transaction end-to-end.

**Security warning:** The service authorization must be based on a trusted authenticated actor. The repository's `LoadPending` does not independently enforce organization scope; the service compares the loaded organization to the authorized organization. The `Complete` transaction locks the domain and challenge but does not recheck organization authorization or TXT digest. Only call it from the verification service after a fresh DNS check.

Before public API rollout: add real-driver integration tests, database-backed authorizer, authenticated request context, rate limits, audit, and concurrency tests.
