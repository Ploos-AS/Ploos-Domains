-- Run only against a disposable migrated CI database as its schema owner.
-- Transaction rollback ensures the temporary role and grants do not persist.
BEGIN;
CREATE ROLE ploos_browser_oidc_ci NOLOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
GRANT USAGE ON SCHEMA public TO ploos_browser_oidc_ci;
GRANT SELECT ON public.users, public.user_oidc_identities TO ploos_browser_oidc_ci;
GRANT SELECT,INSERT,UPDATE,DELETE ON public.oidc_login_transactions,public.browser_sessions TO ploos_browser_oidc_ci;

SET LOCAL ROLE ploos_browser_oidc_ci;
DO $$
DECLARE
 item text;
 action text;
BEGIN
 IF (SELECT rolsuper FROM pg_roles WHERE rolname=current_user) THEN
  RAISE EXCEPTION 'runtime role is superuser';
 END IF;
 FOREACH item IN ARRAY ARRAY['public.users','public.user_oidc_identities','public.oidc_login_transactions','public.browser_sessions'] LOOP
  IF (SELECT pg_has_role(current_user,c.relowner,'MEMBER') FROM pg_class c WHERE c.oid=item::regclass) THEN
   RAISE EXCEPTION 'runtime role owns %',item;
  END IF;
  IF NOT has_table_privilege(current_user,item,'SELECT') THEN
   RAISE EXCEPTION 'missing SELECT on %',item;
  END IF;
 END LOOP;
 FOREACH item IN ARRAY ARRAY['public.oidc_login_transactions','public.browser_sessions'] LOOP
  FOREACH action IN ARRAY ARRAY['INSERT','UPDATE','DELETE'] LOOP
   IF NOT has_table_privilege(current_user,item,action) THEN
    RAISE EXCEPTION 'missing % on %',action,item;
   END IF;
  END LOOP;
 END LOOP;
 IF has_table_privilege(current_user,'public.users','UPDATE')
  OR has_table_privilege(current_user,'public.user_oidc_identities','DELETE')
  OR has_table_privilege(current_user,'public.browser_sessions','TRUNCATE') THEN
  RAISE EXCEPTION 'runtime role has excessive privileges';
 END IF;
END $$;
RESET ROLE;

-- Confirm that removing a required privilege makes the role unsuitable.
REVOKE INSERT ON public.browser_sessions FROM ploos_browser_oidc_ci;
SET LOCAL ROLE ploos_browser_oidc_ci;
DO $$
BEGIN
 IF has_table_privilege(current_user,'public.browser_sessions','INSERT') THEN
  RAISE EXCEPTION 'revoked privilege remains effective';
 END IF;
END $$;
RESET ROLE;
ROLLBACK;
