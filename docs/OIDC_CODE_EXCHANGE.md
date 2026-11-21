# OIDC authorization-code exchange

The unmounted `OIDCCodeExchange` implements the callback's `OIDCCodeExchanger` interface. It exchanges a code with a trusted, fixed HTTPS token endpoint, includes PKCE verifier and redirect URI, authenticates the confidential client using HTTP Basic, rejects redirects and non-200/non-JSON/oversized responses, and returns the ID token for **separate cryptographic verification**.

The client secret must come from protected server configuration, not source control. Configure the token endpoint from trusted OIDC discovery; never from the browser. In production use a bounded timeout and a transport with normal TLS certificate validation. This implementation is for confidential clients using client_secret_basic; providers requiring another token authentication method need an explicitly reviewed adapter.

**Not mounted, not deployed.** Still required: trusted reverse-proxy HTTPS/origin handling, provider configuration, end-to-end integration tests, logout/CSRF review and UI wiring.
