# Subdomain allocation (phase 1)

This PR introduces **inventory and reservation**, not live DNS delegation.

Example: an organization that has verified control of `seniornett.no` can reserve `sirdal.seniornett.no` for its Sirdal unit.

## Invariants

- A domain belongs to exactly one organization.
- Allocations refer to a unit and domain in the same organization.
- A domain must be marked verified before allocation; **the verification workflow is not implemented yet**.
- Each allocated FQDN is unique globally in this database.
- Allocations must be strictly beneath their parent domain.
- Domain identity and verified state cannot be silently altered in place.
- Application-side canonicalization accepts ASCII labels only; Unicode IDNs require future IDNA processing.

## Security gaps before production

- No authenticated API or domain verification workflow exists.
- The database's basic hostname checks are less strict than Go's canonicalizer; future writes must use the canonicalizer and consider adding equivalent database enforcement.
- Overlapping parent/child allocations are **not** yet prohibited. For example, `sirdal.example.org` and `kurs.sirdal.example.org` can both be allocated. This requires an explicit delegation policy before granting DNS write access.
- DNS record-level authorization, NS delegation and zone transfer are not implemented.
- Longer unit hierarchy cycles remain a separate unresolved constraint.

No actual DNS records are created or modified by this PR.
