# Manual Authelia OIDC smoke runner

From the repository root, configure the six `PLOOS_TEST_LIVE_OIDC_*` public test values described in [AUTHELIA_OIDC_TEST_ENVIRONMENT.md](AUTHELIA_OIDC_TEST_ENVIRONMENT.md), and explicitly set `PLOOS_TEST_LIVE_OIDC_ENABLED=true`. Run:

```sh
bash scripts/test-authelia-oidc-smoke.sh
```

The runner refuses to run without explicit opt-in and checks that all expected configuration values are present before running the existing three live provider smoke tests. It does **not** provision Authelia, request a login, exchange an authorization code, fetch user data or mount any application routes. It does not require a client secret. Use only an isolated test provider and verify that the configured URLs refer to your intended nonproduction endpoints.

Passing these checks establishes only public Discovery, verifier initialization and JWKS availability; it does not establish full login readiness.
