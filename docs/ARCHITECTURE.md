# Architecture and security

## Scope

Ploos Domains starts as a **modular monolith** written in Go. PostgreSQL stores tenancy, permissions, domain inventory, desired DNS state and immutable audit events. Domain registration, DNS provider actions and monitoring integrations are separate adapter interfaces.

## Core resources

- **Organization**: top-level tenant with owners and policy.
- **Unit**: hierarchical chapter/department, with stable ID and parent organization.
- **Membership**: identity, role, and scope (organization or unit).
- **Domain**: registrant metadata, ownership verification status, registrar adapter.
- **Zone**: authoritative DNS zone; may be external.
- **Subdomain allocation**: a unique DNS name reserved to a unit.
- **DNS change set**: proposed record operations with validation, review, status and audit.
- **Audit event**: append-only principal/action/resource/timestamp/change summary.

## Trust and permissions

All sensitive operations must check the **server-side effective scope**. A unit admin can request changes only within allocated DNS name(s); organization admins manage allocations; domain/zone administration requires separate elevated permission. Prevent scope escalation via CNAME, NS delegation, wildcard records, DNAME, apex records and record types that create control outside an allocation.

Verification of domain control (e.g. short-lived TXT challenge) is required before any external DNS operation. Changes must be validated against the current zone state and provider capabilities; show a diff before applying. Treat DNS caching, TTL, DNSSEC and external delegation as operational concerns, not instant rollback guarantees.

## Deployment and isolation

Keep API, worker, PostgreSQL and credentials logically separated. Encrypt sensitive configuration at rest; use minimal provider tokens. Add rate limits, idempotency keys, CSRF protections for browser flows, MFA/SSO, secure session handling and complete auditability before customer-facing launch.

Authoritative DNS requires redundant, geographically independent nameservers and operational runbooks. It is deliberately deferred. Existing customer DNS providers are supported through adapters where practical.

## API direction

Versioned REST API: `/api/v1`. First contract covers organizations, units, membership, domain verification, allocations, DNS change previews and approvals. API errors use stable codes with localized user messages.

## i18n

Initial locales: `en` and `nb`. Avoid hard-coding user-facing strings into Go handlers.

## Integrations

Use Voqryn for monitoring rather than recreating checks. Keep Ploos Deploy optional for provisioned hosting; Forenio, Ploform and qix integrations remain separate.
