# .no domain search foundation

This is an **internal Go contract**, not a public search API or live registrar integration.

- `CanonicalNOQuery` validates and normalizes **ASCII** second-level `.no` domain inputs. Unicode internationalized domains require a separate IDNA/eligibility layer; do not silently reject them in a future customer UI without guidance.
- `AvailabilityProvider` defines a future authorized registrar/reseller adapter.
- `SearchNO` validates input before contacting the provider and refuses mismatched, malformed or missing-source responses.
- Availability is `available`, `unavailable` or `unknown`; unknown and provider errors must never be displayed as purchasable.
- DNS lookup, RDAP absence or a failed HTTP request are **not** substitutes for a registrar availability response.
- This module does not return prices, check registrant eligibility, reserve domains or execute purchases.

## Next milestones

1. Verify current Norid rules and choose registrar/reseller partner with permission to query availability.
2. Add Unicode IDNA handling and TLD-specific policy validation with appropriate tests.
3. Add provider adapter, rate limiting, caching/TTL, retry/backoff and a provider-outage UX.
4. Define an API and accessible Norwegian/English search UI.
5. Implement a separate quote/checkout flow with confirmed live pricing, identity and registrar order reconciliation.

Never expose an availability result as a guaranteed reservation.
