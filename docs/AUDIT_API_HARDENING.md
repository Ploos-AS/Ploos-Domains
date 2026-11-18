# Audit API hardening

The audit route now has a process-local, shared token-bucket request limit of **60 requests burst / 60 requests per minute**. The limit applies **before** OIDC verification to reduce unauthenticated request load and does not use spoofable forwarding headers.

Every response from the mounted audit route includes:
- `Cache-Control: no-store`
- `X-Content-Type-Options: nosniff`
- `Referrer-Policy: no-referrer`

Rejected requests return HTTP 429 with `Retry-After`.

**Operational limitations:** This is a single, shared limiter per server process, not per-user or per-organization. Heavy legitimate usage by one tenant can throttle other tenants. Limits do not coordinate across replicas and are reset on restart. Configure additional IP-based and distributed per-identity limits at a trusted reverse proxy or gateway before production. Do not trust unvalidated `X-Forwarded-For` headers.

Future work: per-identity quotas after authentication, shared distributed state, audit metrics, end-to-end authentication tests and robust deployment configuration.
