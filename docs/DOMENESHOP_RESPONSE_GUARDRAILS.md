# Domeneshop transport response guardrails

The optional `ReadResponse` helper rejects non-200 responses and caps response bodies at a caller-chosen size (maximum 4 MiB). It never includes remote response bodies in error messages. The optional `NewRestrictedHTTPClient` disables automatic redirects and enforces a timeout.

**Important:** the existing `NewTransport` still accepts an arbitrary injected `HTTPDoer`; an external client may follow redirects. Production wiring must use a redirect-restricted client and a verified, allowlisted API origin. These helpers do not authorize calling any Domeneshop endpoint.

No API endpoints, authentication, registration operations or live requests are implemented. Provider-specific status codes and schemas must be confirmed from official documentation before use.
