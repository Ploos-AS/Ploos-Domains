# Verification persistence (phase 2)

`CompleteVerification` is a low-level PostgreSQL transaction helper. It atomically consumes one unexpired challenge and marks the matching domain verified, locking the domain and challenge rows to prevent concurrent consumption.

## Not an authenticated API

**Do not expose this function directly to untrusted callers.** A service layer must first:

1. Authenticate the requester and authorize access to the domain's organization.
2. Fetch the stored challenge and digest from trusted database state.
3. Perform a fresh TXT lookup using `verification.Verify`.
4. Call `CompleteVerification` with the server-resolved domain and challenge IDs.
5. Apply rate limits, audit logs, DNS error handling, and replay-safe HTTP responses.

The DNS lookup is deliberately outside the transaction to avoid holding database locks during network I/O. This creates a short interval between DNS verification and commit, which must be addressed in the full threat model (e.g. short challenge lifetimes, trusted resolvers, and re-verification policy).

## Known gaps

- No database driver, connection pool configuration, HTTP endpoints or authentication integration yet.
- The helper does not itself check TXT values or digests; callers must do so.
- The current schema does not enforce a database-side maximum challenge lifetime.
- A verified domain cannot currently be unverified via the normal update path; explicit revocation workflow is future work.
- Domain claim ownership conflicts, hierarchical DNS allocations and organization-wide auditing remain future work.
