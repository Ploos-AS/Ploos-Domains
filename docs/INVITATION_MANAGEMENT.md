# Invitation management foundation

`InvitationManager.Issue` creates an identity-bound, short-lived invitation after checking that a **trusted internal actor ID** has an organization-scoped `owner` or `admin` membership. It returns a one-time 256-bit secret to be delivered securely; only its SHA-256 digest is stored.

This operation does **not** create users, assign memberships, or grant roles. The existing invitation-consumption operation likewise does not grant membership. An audited, separately authorized membership-grant workflow is still required.

## Security boundaries

- Actor identity must be resolved from an authenticated session, never from caller-controlled JSON or headers.
- Invitations are bound to the exact OIDC issuer and subject, not an email address.
- Invites expire after at most 24 hours; this limit is currently enforced by the service, not by the database.
- Issuance must be rate-limited and audited before exposure to customers.
- Token delivery, revocation, listing and administrative HTTP endpoints are not implemented.
- Real PostgreSQL authorization integration tests are follow-up work.
- Public self-registration remains disabled.

This PR introduces no externally accessible endpoints.
