BEGIN;
INSERT INTO organizations(id,name) VALUES ('00000000-0000-0000-0000-000000000601','Audit test');
INSERT INTO users(id,external_subject) VALUES ('00000000-0000-0000-0000-000000000602','audit-test-user');
INSERT INTO audit_events(id,organization_id,actor_user_id,target_user_id,event_type)
VALUES ('00000000-0000-0000-0000-000000000603',
 '00000000-0000-0000-0000-000000000601',
 '00000000-0000-0000-0000-000000000602',NULL,'invitation.issued');
DO $$
BEGIN
 BEGIN
  INSERT INTO audit_events(id,organization_id,event_type)
  VALUES ('00000000-0000-0000-0000-000000000604',
   '00000000-0000-0000-0000-000000000601','secret.token.value');
  RAISE EXCEPTION 'unexpected audit event type accepted';
 EXCEPTION WHEN check_violation THEN NULL;
 END;
END $$;
ROLLBACK;
