# Audit log database security

Apply `db/security/audit_roles.sql` **as a privileged migration/schema owner**, after migration 000006. It creates two PostgreSQL **NOLOGIN group roles**:

| Role | Allowed on `audit_events` |
|---|---|
| `ploos_audit_writer` | INSERT only |
| `ploos_audit_reader` | SELECT only |

Neither group can UPDATE or DELETE audit rows. The default PUBLIC table privileges are revoked. The schema owner, superusers, and users inheriting their privileges can still modify the table: **do not use the migration/schema owner for the application**.

These roles are building blocks, not production-ready tenant isolation. The writer must have required privileges on other application tables through separately configured least-privilege roles. The reader role is powerful: it can see *all organizations*. Do not grant it to a customer-facing account. Implement a trusted organization-authorized query service or row-level security before exposing audit reads.

PostgreSQL row ownership and foreign-key cascades require further attention. The current audit schema uses `ON DELETE RESTRICT` for organizations, so audit history can prevent organization deletion. Define retention, archival and deletion policy before production.

Do not add application logins to either role until credentials, deployment configuration, network restrictions and monitoring have been reviewed. A future migration can add tamper-evident chains or triggers, but privilege separation is the primary protection here.
