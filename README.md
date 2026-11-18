# Ploos Domains

**Open-source domain, DNS and organization management.**

Ploos Domains is an international, bilingual (English and Norwegian) platform developed by **Ploos AS**. Its distinctive feature is delegated management of subdomains for organizations with chapters, branches and teams.

## Product areas

- **Domains:** inventory, expiry and registrar/reseller adapters (later).
- **DNS:** provider-agnostic record management, dry runs, audit and rollback planning.
- **Organizations:** organization units, membership and least-privilege permissions, subdomain ownership and delegation.
- **Integrations:** Voqryn monitoring, Ploos Deploy provisioning, Ploform and other services (future).

A domain can remain with its existing registrar and DNS provider. Domain registrations and live authoritative DNS hosting are **not** required for the first MVP.

## First MVP

1. Organization and local-chapter hierarchy.
2. Domain and DNS-zone inventory; verified domain control.
3. Scoped access: organization owner, organization admin, chapter admin and viewer.
4. Chapter subdomain allocation with collision prevention.
5. Provider-independent DNS change proposals (dry-run first), approval and audit trail.
6. Bilingual API-facing error codes and UI-ready localization structure.

**Security boundary:** an administrator for `sirdal.example.org` must never acquire permissions for `example.org`, `bergen.example.org` or their records. No DNS changes are applied before explicit authorization, domain verification and provider capability checks.

## Technology direction

- Go backend, PostgreSQL persistence, OpenAPI contract.
- Modular monolith initially; adapters for DNS providers and domain registrars.
- Authoritative DNS (e.g., PowerDNS) is optional infrastructure, not an MVP requirement.
- No dependence on paid GitHub features for CI; use standard GitHub Actions and Go tooling.
- Self-hostable core under the MIT license; commercial managed hosting/support can be offered separately.

## Documentation

- [Architecture and security](docs/ARCHITECTURE.md)
- [Development roadmap](docs/ROADMAP.md)
- [Contributing](CONTRIBUTING.md)

> Status: early planning/foundation. Do not use for production DNS until security, integration and operational readiness requirements are met.
