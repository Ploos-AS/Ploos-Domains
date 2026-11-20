# Registrar strategy: Domeneshop first, Norid direct later

**Decision (2026-11-20):** Ploos Domains will target **Domeneshop as its initial registrar/reseller integration**, with the long-term goal that **Ploos AS becomes a direct Norid registrar**. This is a product/architecture decision, **not confirmation of an agreement or API access**.

## Non-negotiable boundaries

- Ploos Domains owns the storefront, customer portal, payment UX, domain portfolio, parking and support experience.
- Registrar integrations are replaceable adapters behind internal contracts; never expose provider credentials, raw response formats or provider-specific IDs as the primary public domain identity.
- The registrar of record and its actual status must remain visible and auditable internally; do not represent Ploos AS as an accredited registrar until accredited.
- Ploos DNS is independently usable regardless of domain registrar. Never silently change nameservers or MX records.
- Migration from Domeneshop to a future direct Norid registrar must be a **separate, consented and tested registrar-transfer/migration process**, not a provider configuration switch that silently moves existing domains.

## Domeneshop discovery gate — verify with the provider

Before building a live adapter, obtain written confirmation of:

1. Reseller/partner eligibility and whether Ploos AS may register and manage domains for its own end customers.
2. API capabilities and permissions for live .no availability, pricing, registration, renewal, transfers, contacts and nameserver management; distinguish public/customer API from reseller APIs.
3. Authentication, credential scoping, rate limits, sandbox/test environment, webhooks or polling, and service availability.
4. Domain ownership/registrant identity handling, Norid-specific eligibility and documentation requirements.
5. Wholesale/retail pricing, invoicing, VAT, refunds, settlement, support, abuse and liability responsibilities.
6. Contractual and operational implications of future migration to a direct Norid registrar.

**Do not guess endpoint paths or implement undocumented order requests.** An adapter should be built only against verified API documentation and authorized credentials.

## Internal adapter contracts — proposed

Separate capabilities so a partner's partial API support does not block the rest of the product:

- `AvailabilityProvider` — already introduced in `internal/domains/search.go`. Returns `available`, `unavailable` or `unknown` with provenance.
- `QuoteProvider` — obtains an expiring, currency- and tax-aware registration/renewal quote; never derive checkout prices from an availability response.
- `RegistrationProvider` — idempotent submit and query of order state; distinguish accepted/pending/registered/failed.
- `LifecycleProvider` — renewal, transfer-in/out, registrant updates and status retrieval, each with explicit authorization and audit records.
- `NameserverProvider` — optional registrar-side nameserver updates, never silently invoked by Ploos DNS.

Keep domain identity, order state, customer ownership and provider references in Ploos-owned persistence. Use explicit capability checks and fail closed when the chosen provider does not support an operation.

## Phased implementation

1. Confirm Domeneshop commercial terms, API documentation and .no workflows.
2. Add provider-neutral typed contracts and deterministic fake-provider tests.
3. Implement authorized Domeneshop availability adapter with secure credentials, timeouts, rate limits and failure semantics.
4. Add pricing/quote contract and transparent customer checkout; then idempotent registrar orders and reconciliation.
5. Add renewals, transfer workflows and parking independently of registrar integration.
6. Research Norid direct registrar accreditation, economics and operational requirements.
7. When approved, add a direct Norid adapter, migrate **new** orders first, then offer controlled migration of existing portfolios where lawful and practical.

## Launch guardrails

No live registration, billing or provider-specific API calls are authorized by this document. Treat the provider relationship, access rights, Norid rules and pricing as **unverified until confirmed**.
