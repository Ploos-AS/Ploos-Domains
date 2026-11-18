# OIDC user provisioning

The `UserProvisioner` creates an internal user and links the verified OIDC `(issuer,subject)` in one PostgreSQL transaction. A concurrent request for the same identity resolves to the same internal user. The operation creates **no organization membership, role, or domain permission**.

## Critical trust boundary

`Provision` does **not** verify an OIDC token. Only call it after a cryptographically verified identity and an explicit enrollment-policy decision. In particular, the service must not blindly auto-provision arbitrary bearer tokens or allow an untrusted client to select issuer/subject.

The current `users.external_subject` schema requires a unique non-null legacy value. New users receive an opaque `oidc-provisioned:<UUID>` placeholder. This is transitional; migrate legacy identity handling separately.

## Outstanding work

- Integration tests with PostgreSQL for concurrency and rollback.
- Enrollment policy (invite-only, approved organization or explicitly open signup).
- Auditing, rate limiting, lifecycle management and identity linking controls.
- Authenticated API access using appropriately scoped access tokens or secure sessions.

No HTTP routes are exposed by this PR.
