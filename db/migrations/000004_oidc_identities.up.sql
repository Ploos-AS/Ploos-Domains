BEGIN;

-- Keep legacy users.external_subject intact until an explicit migration
-- strategy is approved. Never infer an OIDC issuer from a legacy subject.
CREATE TABLE user_oidc_identities (
 issuer text NOT NULL CHECK (
  length(issuer) BETWEEN 9 AND 2048 AND issuer LIKE 'https://%'
  AND issuer = btrim(issuer)
 ),
 subject text NOT NULL CHECK (
  length(subject) BETWEEN 1 AND 1024 AND subject = btrim(subject)
 ),
 user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 created_at timestamptz NOT NULL DEFAULT now(),
 PRIMARY KEY (issuer, subject),
 UNIQUE (user_id, issuer)
);

CREATE INDEX user_oidc_identities_user_idx ON user_oidc_identities(user_id);

COMMIT;
