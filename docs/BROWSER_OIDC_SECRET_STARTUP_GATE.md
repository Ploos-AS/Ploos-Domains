# Browser OIDC secret startup gate

When `PLOOS_BROWSER_OIDC_LOGIN_ENABLED=true`, startup validates browser OIDC configuration, then reads `PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE` with the protected file loader before contacting the provider or PostgreSQL. The file must be an absolute path to a regular, non-symlink file with no group/other permissions and a bounded nonempty secret. Startup rejects missing or unsafe files without returning secret contents in errors.

The secret is validated but **not yet used for code exchange**. Provider discovery and database checks remain required, and startup **always fails closed** after readiness because login handlers are not mounted. When disabled, no secret file is read. No production deployment.
