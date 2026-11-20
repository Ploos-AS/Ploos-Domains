# Domeneshop HTTP transport foundation

This is a **read-only, dependency-injected transport primitive**, not a live Domeneshop API integration. It accepts a verified HTTPS origin and a caller-provided HTTP client. No default HTTP client, credentials, endpoint names, or provider response parsing are configured.

The `Get` method accepts a single relative path segment, disallowing absolute URLs, path traversal, queries, fragments and percent escapes. Requests have a bounded context deadline; callers must close response bodies to release resources.

## Integration gates

- Confirm a reseller agreement, API origin, authentication scheme, allowed endpoints and per-TLD capabilities from official Domeneshop documentation.
- Define a dedicated authenticated client with redirect policy, response-size limits, status/error mapping, retry rules, redaction and rate limiting before any production network calls.
- Implement availability as a documented operation only, with mock tests and no assumption that retail customer API rights imply reseller rights.
- Do not use this generic transport as a public proxy or permit user-controlled path segments.

This PR does not enable domain orders, payments, nameserver edits or live registrar requests.
