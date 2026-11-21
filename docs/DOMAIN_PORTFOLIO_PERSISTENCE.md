# Domain portfolio persistence

Migration `000007` adds a separate `domain_portfolio` table. This deliberately does **not** modify the existing `domains` table, whose DNS verification and subdomain allocation constraints serve a different purpose.

- Portfolio rows are unique per organization and canonical domain; the same externally held domain can appear in different organizations' tracking lists without granting DNS control.
- PostgreSQL queries always include the organization ID. This is **not** a substitute for caller authorization: the service must verify the authenticated user's organization membership and privileges before any read or write.
- `Put` is a low-level administrative persistence primitive, not a public endpoint. Never accept arbitrary organization IDs or verified status from an untrusted caller.
- The repository does not confirm ownership, query a registrar, alter DNS, or initiate a registration.
- Dates and statuses are locally recorded metadata. Before production, add provenance/audit requirements, authorized service methods, pagination and integration tests with a real PostgreSQL database.
- Migration rollback deletes portfolio records; take backups and follow migration procedures before running it on a real database.

No production migration or live server changes are performed by this PR.
