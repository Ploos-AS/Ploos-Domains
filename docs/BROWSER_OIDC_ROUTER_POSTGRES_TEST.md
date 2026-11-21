# Browser OIDC router and PostgreSQL integration

The integration-tagged test exercises the **test-only** login router against a disposable, migrated PostgreSQL database. It checks browser-bound state, PKCE and nonce handoff, successful session persistence, wrong/missing binding rejection, callback replay rejection, session resolution and revocation.

The existing callback and proxy PostgreSQL integration tests remain in place. Token exchange and signed ID-token verification are still faked; no real identity provider is contacted. This test does not register production routes, enable browser login or deploy anything.
