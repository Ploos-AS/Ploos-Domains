-- Run as a database administrator/schema owner AFTER all migrations.
-- Never use the schema owner or migration role as the application login.
-- Roles are NOLOGIN group roles; provision LOGIN credentials separately.
BEGIN;
DO $$
BEGIN
 IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname='ploos_audit_writer') THEN
  CREATE ROLE ploos_audit_writer NOLOGIN;
 END IF;
 IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname='ploos_audit_reader') THEN
  CREATE ROLE ploos_audit_reader NOLOGIN;
 END IF;
END $$;

REVOKE ALL ON TABLE public.audit_events FROM PUBLIC;
REVOKE ALL ON TABLE public.audit_events FROM ploos_audit_writer;
REVOKE ALL ON TABLE public.audit_events FROM ploos_audit_reader;
GRANT USAGE ON SCHEMA public TO ploos_audit_writer,ploos_audit_reader;
GRANT INSERT ON TABLE public.audit_events TO ploos_audit_writer;
-- Readers MUST be further scoped by a trusted service; never grant this
-- role to an untrusted tenant-facing login.
GRANT SELECT ON TABLE public.audit_events TO ploos_audit_reader;
COMMIT;
