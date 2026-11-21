# Optional live OIDC verifier readiness check

In addition to the public Discovery metadata smoke test, the integration suite can initialize the **real** `auth.NewOIDCVerifier` against an explicitly configured test issuer and client ID. This exercises provider discovery through the same Go OIDC library used by the application, without any client secret or user login.

Run with `PLOOS_TEST_LIVE_OIDC_ENABLED=true` and `PLOOS_TEST_LIVE_OIDC_ISSUER` plus `PLOOS_TEST_LIVE_OIDC_CLIENT_ID` set to a dedicated test provider. Command: `go test -tags=integration -run '^TestBrowserOIDCLiveVerifierInitialization$' -count=1 ./internal/server`.

This does **not** fetch and validate an actual signed ID token, exercise authorization code + PKCE, or prove the JWKS endpoint is accessible. No production routes are mounted and CI skips the live check by default.
