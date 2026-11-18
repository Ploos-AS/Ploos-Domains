# DNS TXT domain verification

The initial implementation provides cryptographic challenge generation, DNS TXT matching and a database table for pending challenges. **It does not expose an API, automatically mark domains verified or write DNS records.**

## Intended workflow

1. An authenticated organization administrator registers a domain claim.
2. Server creates a 32-byte random token (maximum 24-hour lifetime) and stores only its SHA-256 digest.
3. Show the TXT record name `_ploos-verify.example.org` and value `ploos-domain-verification=<random hex>`.
4. After the customer publishes the record, the server queries authoritative/public DNS through an appropriately configured resolver and compares the TXT value with the digest.
5. In one transaction, lock the challenge and domain, check expiration, ensure it has not been consumed, mark the challenge consumed and set `domains.verified_at`. These transactional steps are **not yet implemented**.
6. Require re-verification after a change in control and implement expiry/revocation policy.

## Security considerations

- A TXT record proves the ability to modify a specific DNS name at verification time; it does **not** prove legal ownership of the registered domain.
- Never grant access merely because a user knows a domain name. Enforce authenticated organization scope, domain-claim collision rules and appropriate approval.
- Verify against a trusted resolver, guard against stale DNS caches, and rate-limit checks.
- Never log raw tokens. Use one-time consumption and reject expired tokens.
- Avoid treating an arbitrary child-zone TXT record as authorization over a parent zone.
- Only allow domain allocation after verified state is committed.
- Domain names currently use ASCII labels; IDNA validation is future work.

No live DNS writes are performed.
