# Organization and unit model (initial)

The first implementation adds a pure-Go domain model for organizations, units, memberships and authorization. It does **not** create API routes, persistence, identities, domain verification or DNS writes.

## Roles

| Role | Scope | Allowed operations |
| --- | --- | --- |
| owner | organization | view, manage organization, manage units, manage allocated DNS |
| admin | organization | same initial operations as owner |
| viewer | organization or unit | view only |
| unit_admin | one unit | view, manage allocated DNS for that exact unit |

No inherited access to child units is granted to unit-scoped memberships. Roles and actions are allowlisted; unknown values fail closed. The authorization caller must resolve the target organization and unit from trusted storage, not accept arbitrary client-supplied identifiers.

## Next steps

- PostgreSQL migrations with foreign keys, unique constraints and parent ownership checks.
- Authentication and server-side authorization at each API handler.
- Transactional prevention of unit cycles and cross-organization parents.
- Subdomain allocations with canonical DNS-name validation, ownership verification and collision prevention.
- Separate permission for NS delegation and sensitive DNS records.
- Security tests for cross-tenant access and concurrent allocation races.

**Important:** This model is a starting point, not a production-ready access-control system.
