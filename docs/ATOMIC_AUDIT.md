# Atomic audit integration

Organization membership grants and invitation-to-viewer redemptions now write audit events **inside the same PostgreSQL transaction** as their security-sensitive changes. If the audit insert fails, the transaction rolls back and neither the membership nor the invitation consumption commits.

- `MembershipManager.Grant` writes `membership.granted` with the authenticated granting actor and target.
- `InvitationMembership.Redeem` writes `invitation.redeemed` with the authenticated redeeming user.
- No tokens, invitation secrets, or free-form metadata are written to audit events.
- Integration tests verify successful event persistence and that failed operations do not create extra events.

**Not yet covered:** `InvitationManager.Issue` currently persists invitations through a separate writer interface without a transaction-aware audit hook. It must be refactored to record `invitation.issued` atomically in a follow-up PR. Low-level `InvitationStore.ConsumeInvitation` also remains unaudited and should not be used for onboarding. There are no public endpoints in this PR.

Audit table access restrictions, append-only enforcement, retention and dedicated failure-injection tests remain prerequisites for production rollout.
