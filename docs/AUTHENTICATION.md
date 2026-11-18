# Authentication architecture (OIDC foundation)

This PR introduces the HTTP authentication boundary, not a complete OIDC implementation.

The middleware accepts one `Authorization: Bearer <token>` header, delegates verification to an injected `Verifier`, rejects invalid/empty claims, and passes a verified `Identity` through the request context. No protected API routes are enabled yet.

## Required production verifier

A future concrete OIDC verifier must validate:

- JWT signature against trusted JWKS and a fixed allowlist of algorithms;
- configured issuer (`iss`) and audience (`aud`);
- expiry (`exp`), not-before (`nbf`) and clock skew;
- key rotation, HTTP timeouts, and provider discovery constraints.

Never trust unsigned, decoded-only, or client-provided identity claims.

## Internal users

The schema currently has `users.external_subject` as a globally unique string. For multiple OIDC issuers, the internal mapping must use a unique **(issuer, subject)** pair, rather than subject alone. Add the migration and identity lookup before exposing protected endpoints.

## Next milestones

1. Implement an OIDC verifier with a maintained library and fixed provider configuration.
2. Add issuer+subject user mapping in PostgreSQL.
3. Connect verified internal users to the existing `PostgresAuthorizer`.
4. Implement authenticated verification routes, rate limiting and audit.
5. Test authorization and JWT rejection end-to-end.

The module is provider-neutral; Authelia, Keycloak and other standards-compliant OIDC providers are intended integration targets, not yet tested integrations.
