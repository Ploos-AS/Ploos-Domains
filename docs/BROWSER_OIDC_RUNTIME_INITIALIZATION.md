# Browser OIDC runtime initialization

`initializeBrowserOIDCRuntime` is a **standalone, uncalled** initializer. It validates the deployment configuration and trusted OIDC Discovery, opens a dedicated PostgreSQL connection, verifies least-privilege access to required migrated tables, initializes a signed OIDC ID-token verifier, and composes unmounted start/callback handlers. Failure closes the database; success transfers ownership to `browserOIDCRuntime.Close`.

It uses a bounded initialization context and generic errors to avoid leaking DSNs or client secrets. The current `buildHandler` startup path remains unchanged and **always fails closed** when browser login is explicitly enabled. This PR does not register routes, enable login or deploy to production.
