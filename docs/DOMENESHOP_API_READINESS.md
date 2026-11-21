# Domeneshop API integration readiness

PR #39 introduced an inert adapter. This follow-up adds an explicit HTTPS-only, bounded-timeout connection configuration contract, **without creating an HTTP client or calling any endpoint**.

See [the public API capability audit](DOMENESHOP_API_AUDIT.md) for documented features and the outstanding registrar-commerce gap.\n\n## Verified facts vs open questions

The architecture decision is to use Domeneshop first. **The following remain unverified** until reviewed against current official Domeneshop documentation and the commercial agreement:

- Reseller permission and registration/renewal/transfer authority for Ploos AS.
- Supported TLDs and capability differences by TLD.
- Public v0 API documentation identifies HTTP Basic Auth and domain/DNS/invoice/forwarding endpoints. A separate registrar-commerce API, if any, is unverified.
- API rate limits, sandbox availability, error codes, operational SLAs and terms.
- Whether the normal customer API differs from the reseller/partner API.

Never treat the placeholder host used in tests as a production URL.

## Implementation gates

1. Obtain written reseller authorization and current official API documentation.
2. Document supported operations and TLDs in a verified capability matrix.
3. Define credential loading from deployment secrets and redaction rules; never commit secrets.
4. Implement a transport abstraction with context cancellation, strict timeout, no credentials in logs, and mocked tests.
5. Only then implement documented read-only availability calls; leave ordering, DNS and payment untouched.
6. Add provider-specific contract tests with recorded sanitized fixtures or a sanctioned sandbox.

No changes to domains, DNS, invoices, payments, server deployments or provider accounts occur in this PR.
