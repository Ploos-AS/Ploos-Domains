# Atomic invitation redemption and membership

`InvitationMembership.Redeem` combines one-time invitation consumption and creation of an organization-level **viewer** membership in a single PostgreSQL transaction.

- Invitation is bound to exact verified OIDC issuer and subject.
- That identity must already be linked to the authenticated internal user.
- The invitation must be unconsumed and unexpired.
- An existing organization membership is never overwritten; redemption fails and the invitation remains unused.
- A new membership is always `viewer`; invitation holders cannot request `admin` or `owner`.
- Concurrent attempts for the same invitation are serialized by a row lock.

**Trust boundary:** The method does not authenticate requests or validate OIDC signatures. Its caller must use trusted verified identity and resolved internal user ID. It does not provide public endpoints. This PR does not change the existing low-level `ConsumeInvitation` method; callers implementing onboarding should use the atomic method instead.

Invitation issuance, identity verification and public API wiring remain separate. Audit logging and rate limiting must be implemented before exposing this flow.
