# OIDC deployment configuration validation

The new `OIDCDeploymentConfig` validates fixed HTTPS issuer, authorization and token endpoints, client ID, callback host, public host, trusted proxy peer IPs, and login/session lifetimes. It constructs an unmounted login-start configuration with an explicit trusted HTTPS boundary. No credentials are stored in this configuration type.

**Security and deployment limitations:** this is structural validation, **not** proof that authorization/token endpoints belong to the issuer. Before enabling login, resolve and compare endpoints using trusted OIDC discovery metadata and enforce approved provider metadata. Keep client secrets in protected runtime secret storage, not Git. Configure Caddy to overwrite forwarded headers and bind the app to loopback/private network. Real-provider tests, callback wiring, logout and CSRF review remain required. This PR does not mount endpoints or deploy anything.
