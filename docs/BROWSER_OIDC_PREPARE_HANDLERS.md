# Browser OIDC handler composition

`prepareBrowserOIDCLoginHandlers` combines a validated deployment configuration, an already-open restricted PostgreSQL connection, a trusted signed-ID-token verifier, and a discovery validator. It builds the transaction stores, identity lookup, session writer and protected-file confidential code exchanger, then delegates to `server.PrepareOIDCLoginHandlers`.

The returned start and callback handlers are **not mounted**. The helper does not open or check the database, initialize the OIDC verifier, change startup behavior, or enable browser login. These remain separate prerequisites before activation. No production deployment.
