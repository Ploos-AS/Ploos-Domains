# .no domain commerce and parking — product vision

**Ambition:** make Ploos Domains Norway's simplest, most polished place to **find, register, park, transfer and manage .no domains**. This is a product goal, not a claim about current market position or launched capabilities.

Ploos Domains remains an international, bilingual Norwegian/English platform. **.no is the first commercial launch focus**; the architecture must allow additional TLDs, registrars and DNS providers without hard-coded Norwegian rules.

## Customer experience principles

1. **Search first:** fast domain lookup with clear availability, eligibility and exact first-year/renewal pricing. Availability and prices must be fetched from an authoritative registrar/reseller integration; never invent them.
2. **Minimal checkout:** transparent taxes, recurring charges, renewal terms and identity/eligibility requirements before purchase.
3. **Domain-only is a complete product:** no compulsory web hosting, mailbox or Workspace purchase.
4. **Safe defaults:** optional parking, explicit DNS consent, clear renewal notices and straightforward export/transfer-out.
5. **Friendly DNS:** visual record editing, templates, previews, change confirmation, audit history and advanced mode for experienced users.
6. **Fast portfolio management:** domain overview, expiry dates, auto-renewal settings, roles and bulk operations with explicit confirmation.
7. **Accessible and bilingual:** Norwegian/English from the beginning, responsive design and WCAG-aligned usability.

## Customer journeys

### Discover and register

Search a domain → obtain live availability and eligibility → show registration and renewal price, including taxes → authenticate and collect required registrant details → complete payment authorization → submit registrar order → reconcile asynchronous outcome → show domain and next steps. **Never claim registration succeeded before the registrar confirms it.**

### Park a domain

Choose parking after registration or on an existing verified domain → select a branded default page or customizable landing page → preview → explicitly approve DNS routing changes → verify resolution and TLS → manage or disable parking independently of registration. If DNS is externally hosted, present instructions or request authorization; never change nameservers or records silently.

### Manage and transfer

Show current registrar, nameservers, DNS provider, expiry and renewal settings when supported by verified integrations. Support registrar transfer-in/out and DNS-hosting changes as **distinct workflows**. Explain any eligibility, authorization or lock requirements and do not imply all .no transfer processes match generic gTLD flows.

## Product boundaries

- **Ploos Domains:** domain search, registration/resale, portfolio, renewal, transfers, parking and delegated organization access.
- **Ploos DNS:** independently usable DNS management and optional managed authoritative DNS in a later phase. Customers may keep external registrars or DNS providers.
- **Ploos Web / Mail / Workspace:** optional separately purchased services; purchasing Web cannot change MX, SPF, DKIM or DMARC records without explicit approval.
- **Ploos Parking:** a Domains feature/product entitlement that can be offered without Web or Mail. Potential tiers: basic default page; branded landing page; domain-for-sale contact page. Pricing and availability remain undecided.

## Registrar and Norwegian policy dependencies

- Decide whether Ploos AS will become a Norid registrar or contract with an existing registrar/reseller. Confirm eligibility, API access, billing, customer support obligations and contractual terms before coding checkout.
- Validate current Norid .no eligibility, registrant identification, transfer, renewal, deletion and contact requirements against official rules and chosen partner. Do not embed unverified assumptions in product code.
- Define secure payment authorization, refunds, VAT handling, domain lifecycle reconciliation, fraud controls, privacy/data retention and dispute handling.
- Test duplicate purchase attempts, asynchronous registration failures, pending payments, expired sessions and registrar outages.

## Parking security and trust

- No parking page or for-sale listing should imply ownership beyond the verified registrant's authorization.
- HTTPS certificates and redirect loops require explicit handling; do not assume all external DNS can be modified.
- For-sale contact forms need spam/abuse controls and privacy disclosures.
- Keep domain registration, parking and DNS service cancellation independent.

## First commercial release — proposed acceptance criteria

- Live .no availability and pricing through an approved registrar partner.
- Verified eligibility and registrant identity workflow where required.
- Idempotent registration order, explicit status and failure reconciliation.
- Clear renewal price and notifications, plus transfer-out support.
- Domain-only purchase without hosting upsell requirements.
- Basic parking with preview, explicit DNS approval, HTTPS and reversible disable.
- DNS record changes audited and authorized; no silent MX/NS changes.
- Norwegian/English mobile-friendly, accessible checkout and portfolio.
- Monitoring, support escalation, backups, GDPR documentation and security testing.

**Status:** proposal and acceptance criteria only; no registrar agreement, pricing, payment integration or production release is asserted.
