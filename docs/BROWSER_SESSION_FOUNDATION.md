# Browser session foundation

This change introduces **primitives only** for future backend-managed browser sessions:

- Cryptographically random opaque session tokens (32 random bytes).
- SHA-256 digest for server-side storage; never persist the raw cookie token.
- `__Host-ploos_session` cookie with Secure, HttpOnly, Path=/, no Domain, SameSite=Lax.
- Session TTL between 1 minute and 24 hours; explicit cookie clearing.

**Not implemented:** session persistence, lookup, expiration enforcement, revocation, CSRF protection for state-changing routes, OIDC authorization-code callback, PKCE, state/nonce verification, login/logout endpoints, organization selection, or UI wiring. These must be completed before treating the portal as authenticated.

Security considerations: enforce HTTPS and trusted reverse-proxy configuration; avoid logging raw tokens; hash comparison should use a safe server-side lookup; rotate session identifiers after login; add idle timeout and revocation as needed. SameSite is defense in depth, not a substitute for CSRF controls.

No deployment or production configuration changes.
