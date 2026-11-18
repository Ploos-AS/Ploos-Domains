# Roadmap

## Phase 0 — Foundation

- [x] Define product boundaries and MIT licensing.
- [x] Document modular architecture and authorization principles.
- [ ] Establish Go project layout, static checks and CI.
- [ ] Agree on API contract and initial schema.

## Phase 1 — Organizations and subdomains (MVP)

- [ ] Organization/unit hierarchy and scoped membership.
- [ ] Domain inventory and proof of control.
- [ ] Atomic allocation of chapter subdomains; reject collisions.
- [ ] Role and resource policy tests (including cross-tenant attacks).
- [ ] DNS changes as dry-run proposals, review and immutable audit records.
- [ ] First external DNS provider adapter in a sandbox.
- [ ] English and Norwegian UI.

## Phase 2 — DNS platform

- [ ] Provider abstraction and additional adapters.
- [ ] Record templates, change history and safe restore workflows.
- [ ] Monitoring via Voqryn; DNSSEC and email-DNS health checks.
- [ ] Self-service organization portal.

## Phase 3 — Domain commerce

- [ ] Registrar/reseller partner evaluation and agreements.
- [ ] Domain registration, transfer, renewal and billing adapters.
- [ ] Customer ownership and lifecycle controls.

## Phase 4 — Advanced infrastructure

- [ ] Optional redundant authoritative DNS infrastructure.
- [ ] Genuine NS delegation for chapter-owned zones.
- [ ] White-label services and developer APIs.

## Non-goals for v1

- Becoming an ICANN-accredited registrar on day one.
- Moving a customer's nameservers automatically.
- Running a production DNS network on one VPS.
- Charging per DNS record or subdomain as the primary pricing model.
