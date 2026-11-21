# Domeneshop API capability audit (2026-11-21)

## Official public sources

- Domeneshop API v0: https://api.domeneshop.no/docs/
- Domeneshop FAQ: https://domene.shop/faq?id=370

These are **public API capabilities**, not confirmation of a reseller agreement, wholesale terms or access to a private reseller API.

## Confirmed from public documentation

| Operation | Public v0 API | Notes |
| --- | --- | --- |
| List existing account domains | Documented | GET domain listing/details |
| Read invoices | Documented | Existing account invoices |
| DNS records | Documented | Read/create/update/delete, requires explicit authorization |
| HTTP forwarding | Documented | Read/create/update/delete |
| Dynamic DNS | Documented | Update endpoint; not a domain registration feature |
| Availability lookup for unregistered domains | **Not documented** | Do not infer from domain listing or DNS |
| New domain registration | **Not documented** | FAQ directs users to web registration flow |
| Renewal, transfer or registrant changes via API | **Not documented** | Do not assume these are supported |
| Registrar reseller/wholesale access | **Unverified** | Requires direct commercial confirmation |

Public API v0 uses HTTP Basic Auth with an API token and secret. Official FAQ states the API is currently available only on the legacy platform, not the new platform. Confirm current applicability with Domeneshop before deploying.

## Architectural implications

- The existing `internal/registrars/domeneshop` code is a fail-closed foundation; it is **not** a domain-commerce integration.
- Do not mark `OperationAvailability`, `OperationRegistration`, `OperationRenewal` or transfer capabilities as supported by Domeneshop solely based on public v0 documentation.
- A future *account management* integration (read-only domain inventory and DNS) must be kept separate from the *domain commerce* registrar adapter.
- No fallback to RDAP absence, DNS resolution or website availability as proof that a domain is registrable.
- If reseller registration APIs are unavailable, require a separate approved partner or an explicit manual fulfillment model; never present manual orders as automated.
- Keep the Norway-first multi-TLD strategy, but treat each TLD and operation as unverified until contractually confirmed.

## Questions for Domeneshop

1. Does Domeneshop offer a reseller/partner agreement allowing Ploos AS to register domains for end customers and invoice them under its own brand?
2. Is there a separate private API for availability, registration, renewal, transfer, nameserver and registrant operations?
3. Which TLDs are included and which operations are available per TLD?
4. Are sandbox access, API documentation, pricing, rate limits and credentials available to approved partners?
5. Who is the registrar of record, who is the legal registrant, and what are the invoicing, customer support and liability arrangements?
6. How would a later move to direct Norid accreditation work for new and existing .no registrations?

## Decision gate

**Do not implement live domain registration or promise automated checkout until Domeneshop confirms reseller permissions and a supported integration method in writing.** No production API calls, credentials or domain changes are introduced by this audit.
