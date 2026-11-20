# Ploos Domains roadmap

Ploos Domains is an independent product in the broader [Ploos Hosting portfolio](https://github.com/Ploos-AS/Ploos-Hosting). This roadmap tracks **implemented foundations** separately from the **unreleased product experience**. Checked items indicate repository-level implementation, **not production deployment or full feature readiness**.

## Phase 0 — Foundation and persistence

- [x] Define open-source MIT scope and modular architecture.
- [x] Establish Go application, CI, PostgreSQL migrations and integration tests.
- [x] Add organization, unit, user and membership schema.
- [x] Add subdomain allocation constraints and verification persistence.

## Phase 1 — Domain control and authorization foundations

- [x] Add DNS TXT challenge generation and domain-control verification service.
- [x] Add PostgreSQL challenge repository and concurrent integration tests.
- [x] Add organization-scoped domain authorization checks.
- [x] Add OIDC ID-token verification middleware and OIDC-to-internal-user mapping.
- [x] Add invitation-based enrollment, membership grants and audited invitation redemption.
- [x] Add append-oriented audit events and atomic audit recording in selected mutation paths.
- [x] Add organization-scoped, paginated audit query and authenticated HTTP endpoint.
- [x] Add opt-in server wiring, rate limiting, security headers and database privilege checks.
- [ ] Complete end-to-end tenant-security review, error handling and security tests.
- [ ] Review all legacy unaudited mutation paths and enforce approved audited paths.
- [ ] Confirm OAuth2 access-token vs OIDC ID-token model before exposing customer APIs.

## Phase 2 — Deployment and pilot

- [x] Add Dockerfile, loopback-bound Compose deployment skeleton and container build in CI.
- [x] Add GHCR image-publishing workflow and immutable digest guidance.
- [ ] Integrate verified image digests, secrets, health checks and rollback with Ploos Deploy.
- [ ] Set up real `domains.ploos.no` routing, TLS, monitoring and backups.
- [ ] Deploy PostgreSQL migrations, dedicated least-privilege users and test restores.
- [ ] Validate running service with end-to-end OIDC and domain verification tests.
- [ ] Define incident response, retention and operational ownership.

## Phase 3 — Customer-facing domain experience and organization/DNS MVP

- [ ] Define mobile-first multi-TLD search, transparent per-TLD pricing and domain-only checkout journeys, with .no prioritized for Norwegian customers.
- [ ] Specify parking flows (default, branded and for-sale), including preview, HTTPS and explicit DNS consent.
- [ ] Design customer domain portfolio, renewal notices, transfer-in/out and clear domain lifecycle status.
- [ ] Confirm .no-specific eligibility, identity and registrar process requirements with Norid and a chosen partner.
- [ ] Document first-release acceptance criteria in [the .no commerce and parking vision](NO_DOMAIN_COMMERCE.md).

- [ ] Complete OpenAPI contract, errors and API surface for organizations/domains.
- [ ] Deliver bilingual English/Norwegian UI and usable self-service flows.
- [ ] Complete scoped chapter/branch delegated administration and abuse tests.
- [ ] Implement DNS dry-run proposals, approval, safe execution and audited restore planning.
- [ ] Integrate first external DNS provider in sandbox; never change customer records silently.
- [ ] Confirm product-level activation and cancellation independence from Ploos Web/Mail.

## Phase 4 — Ploos DNS as a separately usable service

- [ ] Define Ploos DNS as independently purchasable/usable with external registrars.
- [ ] Separate registrar transfer, DNS provider change and authoritative DNS hosting in product flows.
- [ ] Provide safe record editor, preview, templates, approval and audit history.

## Phase 4a — DNS provider platform and integrations

- [ ] Provider abstraction and additional DNS adapters.
- [ ] DNS templates and safe provider synchronization.
- [ ] Voqryn-powered monitoring, DNSSEC checks and mail-DNS health checks.
- [ ] Integration with Ploos Web, Mail, Workspace and Apps through explicit approved workflows.

## Phase 5 — Domeneshop multi-TLD domain commerce and parking launch

- [ ] Verify Domeneshop reseller agreement, API permissions, supported TLD catalogue and per-TLD operations as the **first integration**; see [registrar strategy](REGISTRAR_STRATEGY.md).\n- [ ] Research direct Norid registrar accreditation as the **long-term target**, without assuming existing domains can be migrated automatically.
- [ ] Implement a provider-neutral registrar contract and deterministic test adapter.\n- [ ] Integrate authorized Domeneshop multi-TLD availability, eligibility and first-year/renewal pricing once access is confirmed; prioritize .no UX.
- [ ] Implement idempotent registration orders, payment handling and asynchronous registrar reconciliation.
- [ ] Add transfer-in/out, renewal, cancellation and billing workflows with independent customer lifecycle management.
- [ ] Launch basic parking with reversible DNS changes; evaluate branded and for-sale landing pages.
- [ ] Validate accessible bilingual checkout, transparent pricing, abuse prevention and customer support readiness.
- [ ] Launch additional supported TLDs incrementally as their API permissions, rules, prices and lifecycle workflows are verified.

## Phase 6 — Optional advanced infrastructure

- [ ] Evaluate redundant authoritative DNS hosting (not required for MVP).
- [ ] Support true NS delegation for chapter zones when justified.
- [ ] Add white-label and advanced developer APIs.

## Non-goals for v1

- ICANN accreditation or independent authoritative DNS as a launch prerequisite.
- Automatically transferring a customer's nameservers.
- Overwriting MX or other unrelated DNS records when provisioning web hosting.
- Running a production DNS network on one VPS.

**Status note (2026-11-18):** The checked implementation milestones above reflect the merged foundational PRs through #31. The published GHCR workflow is a capability, not evidence that a usable customer portal or production service has been deployed. The full Ploos Hosting product strategy is maintained separately.
