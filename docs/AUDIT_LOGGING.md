# Security audit log foundation

This migration introduces `audit_events`, a structured table for organization-scoped security events. It records event type, organization, optional actor and target user IDs, and timestamp.

Allowed events: invitation issued/redeemed and membership granted/revoked/role changed. There are **no free-form payloads**, bearer tokens, raw invitation secrets, IP addresses or OIDC credentials in this schema.

`AuditRecorder.RecordTx` must be invoked **inside the same transaction** as the audited security operation. A failed audit insert must fail the operation; callers must not log after committing. This PR provides the schema and recorder but does **not yet wire it into existing issuance, redemption or grant operations**. Therefore, these operations are **not yet audited**.

## Follow-up requirements

- Integrate audit writes atomically into invitation redemption and membership grant.
- Integrate issuance audit with invitation creation in one transaction.
- Apply retention and access-control policies to the audit table.
- Prevent unauthorized updates/deletes with dedicated database roles or append-only enforcement.
- Build safe, organization-scoped audit viewing with pagination and authorization.
- Never expose audit data across organizations.
