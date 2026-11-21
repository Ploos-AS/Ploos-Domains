# Browser OIDC callback dependency preparation

`prepareBrowserOIDCCallbackDependencies` validates the deployment config and requires transaction consumer, ID-token verifier, identity resolver and session writer before constructing the confidential code exchanger from the protected client secret file. It refuses a pre-supplied exchanger to avoid accidentally overriding trusted configuration. Session TTL is taken from validated deployment settings.

This is an **unmounted composition helper**; it does not open database connections, perform token exchange or activate browser login. The actual startup integration and runtime dependency provisioning remain future work. No production deployment.
