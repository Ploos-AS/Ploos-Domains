# Browser OIDC PostgreSQL readiness

`checkBrowserOIDCDatabase` checks for required migrated tables and privileges using the **current application database role**, without mutating the database. It rejects superusers and roles that own required tables. Required tables: `users` and `user_oidc_identities` (read); `oidc_login_transactions` and `browser_sessions` (read, insert, update, delete).

The check is a separate primitive and **is not yet invoked by buildHandler**. Browser OIDC activation remains disabled. The PostgreSQL integration test asserts that a superuser is rejected when the test database connection is a superuser; a dedicated restricted-role integration test is still needed.

Before enabling login, provision a separate least-privilege role, confirm column-level requirements and indexes, verify RLS/role policies as applicable, and connect the check to startup with bounded timeouts. No migration, route mounting or production deployment is included.
