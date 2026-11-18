# Invitation-based enrollment foundation

This PR introduces an invitation table and a one-time consumption operation. Each invitation is bound to an **exact OIDC issuer+subject**, an organization, an expiration, and a SHA-256 digest of a 256-bit random bearer token.

A trusted service must authenticate the user, verify the OIDC identity, resolve or provision the internal user under an approved enrollment policy, and then consume the invitation. Consumption checks that the authenticated identity is already linked to the internal user.

**No organization membership or role is granted by consuming an invitation.** Membership assignment and invitation issuance are deliberately separate, privileged and audited future workflows. Therefore this PR alone does not complete user onboarding.

## Security boundaries

- Do not expose the raw token in logs; store only the digest.
- Invite issuance must require organization-admin authorization and enforce TTL/rate limits.
- The invite must target the exact issuer and subject; email alone is not an identity proof.
- This operation must never accept a caller-controlled internal user ID without prior authenticated identity resolution.
- Do not enable public self-registration by default.
- Integration tests and audit logging are required before public API rollout.
