# Audited invitation issuance

Use `AuditedInvitationIssuer.Issue` for PostgreSQL-backed invitation creation. It locks the issuer's organization-scoped owner/admin membership and creates both the invitation and its `invitation.issued` audit event in one transaction. If either write fails, neither is committed. Only the SHA-256 token digest is persisted.

The legacy `InvitationManager.Issue` and `PostgresInvitationWriter.Create` are **not audited** and must not be used for production issuance. A later compatibility cleanup should remove or deprecate those entry points.

The caller must resolve `actorID` from an authenticated identity, not request input. No public HTTP endpoint is introduced. Rate limiting, audit-table access controls, append-only enforcement, retention and explicit failure-injection tests remain pending.
