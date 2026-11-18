-- Verify a non-owner application login can read required audit data
-- without acquiring audit history mutation privileges.
BEGIN;
DO $$
BEGIN
 IF NOT EXISTS(SELECT 1 FROM pg_roles WHERE rolname='ploos_audit_api_ci') THEN
  CREATE ROLE ploos_audit_api_ci NOLOGIN;
 END IF;
END $$;
GRANT ploos_audit_reader TO ploos_audit_api_ci;
GRANT USAGE ON SCHEMA public TO ploos_audit_api_ci;
GRANT SELECT ON public.memberships,public.user_oidc_identities TO ploos_audit_api_ci;
SET LOCAL ROLE ploos_audit_api_ci;
DO $$
BEGIN
 IF NOT has_table_privilege(current_user,'public.audit_events','SELECT')
 OR has_table_privilege(current_user,'public.audit_events','UPDATE')
 OR has_table_privilege(current_user,'public.audit_events','DELETE')
 OR has_table_privilege(current_user,'public.audit_events','TRUNCATE')
 OR NOT has_table_privilege(current_user,'public.memberships','SELECT')
 OR NOT has_table_privilege(current_user,'public.user_oidc_identities','SELECT')
 THEN RAISE EXCEPTION 'audit API role privilege boundary violated';
 END IF;
END $$;
RESET ROLE;
ROLLBACK;
