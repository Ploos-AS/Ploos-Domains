# Multi-TLD search foundation

The provider-neutral `TLDRegistry` explicitly enables supported extensions. This is an **internal Go foundation**, not evidence of Domeneshop reseller API access, supported TLD catalogue, or a public domain-search service.

- The current implementation supports ASCII second-level domains with a single-label TLD, such as `example.no` and `example.com`.
- No TLD is enabled by default; enabled entries must come from verified commercial/provider capabilities.
- Disabled and unknown TLDs fail before the provider is called.
- Provider results must match the queried domain, identify their source, and use the known available/unavailable/unknown states.
- Unicode IDN handling, public suffixes such as `co.uk`, per-TLD registry rules, price quoting and registration are separate future milestones.
- Existing `SearchNO` remains unchanged for compatibility.

## Next steps

1. Verify Domeneshop reseller agreement and actual supported TLD/API operations.
2. Define a richer TLD policy catalogue (eligibility, IDN, period, transfer/renewal and provider capabilities).
3. Implement standards-compliant IDNA normalization and multi-label suffix handling.
4. Implement an authorized Domeneshop adapter with timeouts, rate limiting and tests.
5. Expose a secure API and accessible bilingual domain-search UI.

Availability is advisory; never treat it as a reservation or confirmed order.
