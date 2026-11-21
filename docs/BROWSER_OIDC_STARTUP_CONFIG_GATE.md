# Browser OIDC startup configuration gate

When `PLOOS_BROWSER_OIDC_LOGIN_ENABLED=true`, startup now parses and validates the operator-controlled browser OIDC environment using `loadBrowserOIDCConfig`. Invalid settings fail startup with `errBrowserOIDCConfig`. Even with valid settings, startup deliberately fails with `errBrowserLoginNotConfigured` until provider discovery, PostgreSQL dependencies, protected client secrets, and handler registration are wired and tested.

With the flag unset or `false`, browser OIDC configuration is ignored and existing startup behavior is unchanged. All other flag values are rejected. No OIDC browser routes are mounted and no production deployment occurs.
