# Domeneshop adapter foundation

This package is an **inert, fail-closed placeholder**, not a live Domeneshop integration. It implements the existing `domains.AvailabilityProvider` interface but does not make HTTP requests, accept credentials or return a positive availability result.

The current capability catalogue is operator-provided. **No Domeneshop TLD permissions or endpoints have been verified**, and no extension is enabled by default.

## Required before implementing API calls

1. Confirm reseller/partner agreement and permission for Ploos AS to sell domains to its customers.
2. Obtain current authorized API documentation and verify whether reseller operations differ from customer-account API access.
3. Confirm availability lookup, supported TLDs, rate limits, authentication, errors and test environment.
4. Add credential handling through secure runtime configuration, never committed secrets.
5. Add explicit HTTP timeouts, retries with backoff, rate limiting and mocked integration tests.
6. Keep availability, quote, registration, transfer and renewal as distinct capabilities and workflows.
7. Preserve fail-closed behavior for unsupported TLDs, API outages and ambiguous responses.

**No live domains, DNS records, registrar orders or payments are changed by this PR.**
