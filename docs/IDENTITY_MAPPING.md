# OIDC identity mapping

An authenticated OIDC identity is identified by the pair **(issuer, subject)**, not by subject alone. The new `user_oidc_identities` table maps that pair to an internal `users.id`.

`IdentityRepository.ResolveUserID` performs **lookup only**. Unknown identities are rejected; there is no implicit user creation or account linking.

## Legacy migration

The existing `users.external_subject` field is retained. Its values have no reliable issuer provenance, so this migration deliberately **does not** auto-link them. Administrators must verify identity ownership before adding links. Removing or relaxing the legacy field and implementing safe provisioning is a separate migration.

## Security requirements

- Only pass `Identity` values produced by a correctly configured, cryptographically verified OIDC provider.
- Never trust issuer or subject supplied in HTTP request bodies or headers.
- Account linking must require explicit, authenticated approval and be audited.
- Multiple issuers can use the same subject without collision.
- The database uniqueness constraint permits only one subject per issuer per internal user. Supporting multiple identities from one issuer requires a reviewed schema change.
- This PR does not add public routes or identity provisioning.

The initial SQL issuer check is a format guard, not a complete URL parser. Production configuration validation belongs in the OIDC verifier.
