-- PostgreSQL verification persistence smoke tests. Runs in rollback-only transaction.
BEGIN;
INSERT INTO organizations(id,name) VALUES
 ('00000000-0000-0000-0000-000000000001','Seniornett'),
 ('00000000-0000-0000-0000-000000000002','Other');
INSERT INTO domains(id,organization_id,fqdn) VALUES
 ('00000000-0000-0000-0000-000000000201','00000000-0000-0000-0000-000000000001','seniornett.no'),
 ('00000000-0000-0000-0000-000000000202','00000000-0000-0000-0000-000000000002','example.org');
INSERT INTO domain_verification_challenges
 (id,domain_id,token_sha256,record_name,expires_at) VALUES
 ('00000000-0000-0000-0000-000000000301','00000000-0000-0000-0000-000000000201',
  decode(repeat('ab',32),'hex'),'_ploos-verify.seniornett.no',now()+interval '1 hour'),
 ('00000000-0000-0000-0000-000000000302','00000000-0000-0000-0000-000000000202',
  decode(repeat('cd',32),'hex'),'_ploos-verify.example.org',now()+interval '1 hour');
DO $$
DECLARE n integer;
BEGIN
 SELECT count(*) INTO n FROM domain_verification_challenges c
 JOIN domains d ON d.id=c.domain_id
 WHERE c.id='00000000-0000-0000-0000-000000000301'
 AND d.organization_id='00000000-0000-0000-0000-000000000001';
 IF n<>1 THEN RAISE EXCEPTION 'expected own challenge'; END IF;
 SELECT count(*) INTO n FROM domain_verification_challenges c
 JOIN domains d ON d.id=c.domain_id
 WHERE c.id='00000000-0000-0000-0000-000000000302'
 AND d.organization_id='00000000-0000-0000-0000-000000000001';
 IF n<>0 THEN RAISE EXCEPTION 'cross-organization challenge visible'; END IF;
END $$;
UPDATE domain_verification_challenges SET consumed_at=now()
 WHERE id='00000000-0000-0000-0000-000000000301';
DO $$
DECLARE n integer;
BEGIN
 SELECT count(*) INTO n FROM domain_verification_challenges
 WHERE id='00000000-0000-0000-0000-000000000301' AND consumed_at IS NULL;
 IF n<>0 THEN RAISE EXCEPTION 'consumed challenge still pending'; END IF;
END $$;
ROLLBACK;
