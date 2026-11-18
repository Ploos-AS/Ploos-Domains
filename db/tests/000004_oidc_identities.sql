BEGIN;

INSERT INTO users(id,external_subject) VALUES
 ('00000000-0000-0000-0000-000000000411','legacy-test-411'),
 ('00000000-0000-0000-0000-000000000412','legacy-test-412');

INSERT INTO user_oidc_identities(issuer,subject,user_id) VALUES
 ('https://id.example.org','same-subject','00000000-0000-0000-0000-000000000411'),
 ('https://other.example.org','same-subject','00000000-0000-0000-0000-000000000412');

DO $$
DECLARE matches integer;
BEGIN
 SELECT count(*) INTO matches FROM user_oidc_identities
 WHERE subject='same-subject';
 IF matches<>2 THEN RAISE EXCEPTION 'issuer+subject isolation failed'; END IF;
END $$;

DO $$
BEGIN
 BEGIN
  INSERT INTO user_oidc_identities(issuer,subject,user_id)
  VALUES ('https://id.example.org','same-subject','00000000-0000-0000-0000-000000000412');
  RAISE EXCEPTION 'duplicate identity was accepted';
 EXCEPTION WHEN unique_violation THEN NULL;
 END;
END $$;

ROLLBACK;
