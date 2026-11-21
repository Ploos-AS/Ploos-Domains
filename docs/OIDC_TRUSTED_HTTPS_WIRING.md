# OIDC handlers behind Caddy — HTTPS boundary wiring

The unmounted login-start and callback handlers now accept an optional `TrustedHTTPSBoundary`. When provided, the guard validates a fixed public host and either direct TLS or a specifically trusted proxy peer with strict forwarded headers. Untrusted proxy requests fail before creating or consuming a login transaction.

**Not production-ready:** the boundary is optional for backwards-compatible direct-TLS tests. Production wiring MUST supply the boundary, bind the application to a private/loopback interface, and configure Caddy to overwrite forwarded headers. The login-start Origin check still requires exact `https://<public host>`. No routes are mounted by this PR. End-to-end OIDC and session tests, explicit provider configuration, CSRF/logout handling and deployment review remain required.
