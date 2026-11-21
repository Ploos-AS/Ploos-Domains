-- Disposable CI database only. Never use this fixed password in production.
CREATE ROLE ploos_browser_oidc_runtime_ci LOGIN PASSWORD 'browser-oidc-ci-only-password' NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
GRANT CONNECT ON DATABASE ploos_domains_test TO ploos_browser_oidc_runtime_ci;
GRANT USAGE ON SCHEMA public TO ploos_browser_oidc_runtime_ci;
GRANT SELECT ON public.users, public.user_oidc_identities TO ploos_browser_oidc_runtime_ci;
GRANT SELECT,INSERT,UPDATE,DELETE ON public.oidc_login_transactions,public.browser_sessions TO ploos_browser_oidc_runtime_ci;
