# Ploos Domains ↔ Ploos Deploy integration

## Current status: contract only — NOT DEPLOYABLE

Ploos Deploy currently validates a small JSON manifest (`name`, `image`, `port`, `health_path`), but its guarded deploy implementation is limited to a disposable Balder test app named `deploy-test`, port `18083`. It **cannot deploy Ploos Domains to forseti yet**. Do not run `ploos-deploy deploy --apply` against this example or bypass the guard.

`ploos-deploy.example.json` is a **schema-compatible placeholder** for a future integration. The all-`a` digest is intentionally fake and must be replaced with an actual immutable digest from the successful GHCR publishing workflow. The proposed loopback host port `18088` is not reserved; verify it is unused before deployment.

## Proposed target

- Public origin: `https://domains.ploos.no` (DNS/TLS not configured by this repository).
- Host: `forseti`, after production-readiness checks and explicit deployment approval.
- Container image: `ghcr.io/ploos-as/ploos-domains@sha256:<verified-digest>`.
- Internal HTTP: container port `8080`; host loopback port to be chosen after inventory.
- Health endpoint: `GET /healthz`; health must be checked through the actual deployed container.
- Database: separate PostgreSQL service, migrations applied by a separate privileged account; audit API remains disabled until least-privilege OIDC integration is tested.

## Blockers before first deployment

1. **Ploos Deploy:** add safe, tested support for named apps on forseti, configurable host-to-container port mapping, immutable digest pull, health verification, rollback and state isolation. Keep the disposable Balder guard intact.
2. **Container runtime:** confirm the Ploos Domains binary listens on the configured port, and Compose/Ploos Deploy use compatible environment and networking.
3. **Secrets:** provide protected runtime injection for database and OIDC configuration without committing credentials.
4. **GHCR:** confirm publication succeeded, image is pullable by the deployment identity, and digest is correct.
5. **Routing:** configure DNS and Caddy HTTPS for `domains.ploos.no`, with only a loopback upstream.
6. **Safety:** verify migrations, DB roles, tenant isolation, monitoring, backups, restore tests and rollback procedure.
7. **Acceptance:** run plan/render first, then a controlled staging deployment, smoke tests, and explicit approval before any production apply.

## Manual validation without deployment

The example can be checked against Ploos Deploy's manifest loader using its `plan`/`render` commands in an isolated development environment. This is **not** a guarantee that the resulting Compose file is suitable for Ploos Domains, and must never be interpreted as permission to deploy.

See [Ploos Deploy](https://github.com/Ploos-AS/Ploos-Deploy) for the implementation and its current limitations.
