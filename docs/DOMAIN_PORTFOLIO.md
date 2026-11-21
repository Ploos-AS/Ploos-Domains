# Domain portfolio foundation

This PR adds a provider-neutral, **local data model** for a customer's domain portfolio. It does not read Domeneshop, claim ownership, register a domain or expose a customer API.

Fields: owning Ploos organization ID, canonical domain name, registrar label, locally recorded lifecycle status, optional expiration and last-verification timestamps.

- `unknown` means the record has not been confirmed; `active` must not be interpreted as live registry proof without provenance.
- `FilterPortfolio` returns records for exactly one organization in stable order, with independent timestamp copies.
- **Authorization is not implemented by this helper.** Callers must check organization membership and permissions before retrieving or filtering records.
- Persistence, verified registrar sync, source-of-truth provenance, pagination, transfer workflows and customer UI remain separate milestones.
- Do not expose or populate portfolio records from untrusted clients without server-side organization ownership checks.

Next: design an organization-authorized read service and persistence migration, then a customer-facing portfolio endpoint and accessible bilingual UI.
