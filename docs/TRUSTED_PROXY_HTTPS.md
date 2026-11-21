# Trusted HTTPS boundary for Caddy

`TrustedHTTPSBoundary` accepts direct TLS for the configured public host, or a proxy-forwarded HTTPS request **only** when the TCP peer's IP is explicitly configured as trusted and exactly one `X-Forwarded-Proto: https` and `X-Forwarded-Host` matching the public host are present. Unknown peers, missing/duplicate headers and chains fail closed.

**Important deployment constraints:** restrict application listener to loopback/private networking; configure Caddy to overwrite, not pass through, incoming forwarded headers; set trusted proxy addresses to actual Caddy peer IPs only. Do not trust arbitrary CIDR ranges or Internet clients. Public host must be fixed in trusted configuration.

This PR introduces the boundary primitive and tests only. Existing OIDC handlers still check `r.TLS` directly and are not mounted. A follow-up must inject this boundary into login-start and callback handlers, and test end-to-end behavior behind Caddy. No production deployment.
