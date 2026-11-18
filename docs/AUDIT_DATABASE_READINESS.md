# Audit API database readiness

When `PLOOS_AUDIT_API_ENABLED=true`, startup now checks the database login's effective privileges before mounting the route. It rejects:

- PostgreSQL superusers and users that are members of the audit table owner role.
- Users without SELECT on `audit_events`, `memberships`, or `user_oidc_identities`.
- Users with UPDATE, DELETE or TRUNCATE on `audit_events`.

Provision a separate, non-owner PostgreSQL login for the audit API. Apply migrations and `db/security/audit_roles.sql` using a separate privileged migration account. Grant the application login membership in `ploos_audit_reader` plus SELECT on `memberships` and `user_oidc_identities`. Never grant direct database access to tenants. CI exercises a minimal role with these grants.

This is a **read-side safety check**, not a complete production readiness certification. The application login must not have CREATE or schema ownership privileges, dangerous role-admin grants, or other escalation paths; enforce those in deployment. A compromised database owner/superuser can still change audit history. OIDC ID-token suitability, reverse-proxy TLS, secrets management, backups, audit retention, monitoring and distributed rate limits remain deployment prerequisites.

The current check applies only to the opt-in audit API database connection. Other services that write audit events require separately designed INSERT-only credentials.
