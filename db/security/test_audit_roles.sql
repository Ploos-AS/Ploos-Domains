-- Must run as the schema owner / PostgreSQL test administrator.
BEGIN;
CREATE TEMP TABLE audit_role_assertions(ok boolean NOT NULL CHECK(ok));
INSERT INTO audit_role_assertions VALUES
 (has_table_privilege('ploos_audit_writer','public.audit_events','INSERT')),
 (NOT has_table_privilege('ploos_audit_writer','public.audit_events','SELECT')),
 (NOT has_table_privilege('ploos_audit_writer','public.audit_events','UPDATE')),
 (NOT has_table_privilege('ploos_audit_writer','public.audit_events','DELETE')),
 (has_table_privilege('ploos_audit_reader','public.audit_events','SELECT')),
 (NOT has_table_privilege('ploos_audit_reader','public.audit_events','INSERT')),
 (NOT has_table_privilege('ploos_audit_reader','public.audit_events','UPDATE')),
 (NOT has_table_privilege('ploos_audit_reader','public.audit_events','DELETE'));
DO $$
BEGIN
 IF (SELECT count(*) FROM audit_role_assertions)<>8 THEN
  RAISE EXCEPTION 'audit role assertion count mismatch';
 END IF;
END $$;
ROLLBACK;
