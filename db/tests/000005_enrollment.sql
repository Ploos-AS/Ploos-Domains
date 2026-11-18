BEGIN;
INSERT INTO organizations(id,name) VALUES ('00000000-0000-0000-0000-000000000501','Enrollment test');
INSERT INTO users(id,external_subject) VALUES ('00000000-0000-0000-0000-000000000502','enrollment-test');
INSERT INTO enrollment_invitations
 (id,token_sha256,organization_id,invited_issuer,invited_subject,expires_at)
 VALUES ('00000000-0000-0000-0000-000000000503',decode(repeat('ab',32),'hex'),
 '00000000-0000-0000-0000-000000000501','https://id.example.org','person',now()+interval '1 hour');
DO $$
BEGIN
 BEGIN
  INSERT INTO enrollment_invitations
   (id,token_sha256,organization_id,invited_issuer,invited_subject,expires_at)
  VALUES ('00000000-0000-0000-0000-000000000504',decode(repeat('ab',32),'hex'),
   '00000000-0000-0000-0000-000000000501','https://id.example.org','other',now()+interval '1 hour');
  RAISE EXCEPTION 'duplicate token digest accepted';
 EXCEPTION WHEN unique_violation THEN NULL;
 END;
 BEGIN
  UPDATE enrollment_invitations SET consumed_at=now()
  WHERE id='00000000-0000-0000-0000-000000000503';
  RAISE EXCEPTION 'incomplete consumption accepted';
 EXCEPTION WHEN check_violation THEN NULL;
 END;
END $$;
UPDATE enrollment_invitations SET consumed_at=now(),consumed_by='00000000-0000-0000-0000-000000000502'
WHERE id='00000000-0000-0000-0000-000000000503';
ROLLBACK;
