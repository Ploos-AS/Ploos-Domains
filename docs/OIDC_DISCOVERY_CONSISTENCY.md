# OIDC discovery consistency gate

`OIDCDiscoveryClient.Validate` retrieves discovery metadata from the fixed HTTPS issuer and rejects mismatches in issuer, authorization endpoint or token endpoint. It also requires an HTTPS JWKS URI, limits response size and rejects redirects. No URL is taken from an incoming browser request.

**Scope and limitations:** this first implementation supports only issuers without path components. A provider with a path-based issuer must be handled using the correct OIDC Discovery URL construction before enabling it. HTTP client transport must use normal certificate validation and protected network settings. This is a separate validation primitive: it is not yet wired as a mandatory startup gate, and the login routes remain unmounted. It does not validate JWKS key contents, which are handled by the existing OIDC token verifier. No deployment.
