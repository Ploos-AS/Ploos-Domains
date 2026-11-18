# Organization-scoped audit queries

`AuditQuery.List` provides a database-layer, organization-scoped audit read for organization-level `owner` and `admin` members. Viewers and unit administrators are denied. The caller **must derive actorID from a verified, authenticated OIDC identity**; never trust a client-supplied actor ID.

The method checks the actor's membership in the same read-only transaction as the audit query and locks the membership row against concurrent modification. Results are restricted to the requested organization in SQL and returned newest first, with keyset pagination using `(occurred_at,id)`. Page size is 1–100. It returns only structured event fields, without bearer tokens or invitation secrets.

This is **not an HTTP API**. A production API must authenticate the request, resolve the actor, authorize the requested organization, validate and encode cursors safely, rate-limit access and prevent cross-tenant data disclosure. The database account used for this query needs SELECT on both `memberships` and `audit_events`; the `ploos_audit_reader` role alone does not grant membership SELECT.

The table currently allows privileged owners to mutate or delete rows, and a global audit reader role can see all organizations if granted directly. Do not expose direct database credentials to tenants.
