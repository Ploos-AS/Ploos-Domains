# Organization membership grants

`MembershipManager.Grant` is a database-backed, transactional service for granting **new organization-level** memberships. It accepts only `viewer` and `admin` as target roles:

- An organization `owner` can grant `viewer` or `admin`.
- An organization `admin` can grant `viewer` only.
- Unit administrators, viewers and non-members cannot grant either.
- Nobody can use this method to create an `owner`, change an existing membership, or grant to themselves.

The service locks the granting actor's organization-level membership for the duration of the transaction. Existing memberships are not overwritten.

## Important integration boundaries

The actor ID must come from authenticated identity resolution, never a client-supplied parameter. This service is **not yet wired to invitation redemption**: an invitation token alone cannot authorize membership. A future audited onboarding service must bind the redeemed invitation to its organization and verified identity before calling an appropriately authorized grant workflow. There are no public routes in this PR.

Owner bootstrap/transfer, membership revocation, unit-scoped roles, rate limiting and audit logging require separate reviews.
