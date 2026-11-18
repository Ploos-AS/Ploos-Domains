# Ploos Domains deployment foundation

**Recommended service hostname:** `domains.ploos.no`. Reserve `api.domains.ploos.no` for a later separate API hostname if necessary. Initially, the portal and API can share one origin.

The Docker Compose file binds to **127.0.0.1:8088**, so a trusted reverse proxy (e.g. Caddy) must terminate HTTPS and proxy to that loopback address. Do not publish port 8088 to the internet. Set `PLOOS_DOMAINS_IMAGE` to a published, pinned container image digest. This PR does not publish images or change DNS.

`PLOOS_AUDIT_API_ENABLED` defaults to false. When enabled, configure a trusted OIDC issuer, correct ID-token audience, TLS-protected PostgreSQL DSN and a non-owner, least-privilege database account. The application performs a startup privilege check.

The Compose file demonstrates deployment only; it is **not** a production secret-management solution. Avoid putting database credentials in tracked files, environment examples, shell history or deployment logs. Use a protected environment/secret mechanism on the host. Migrations must be run separately by a privileged account.

Before rollout: verify DNS, HTTPS, trusted proxy headers, database migrations, least-privilege grants, OIDC token semantics, backups, monitoring, log retention and rollback. Integration with Ploos Deploy and image publication will be handled separately.
