BEGIN;
CREATE TABLE domain_verification_challenges (
 id uuid PRIMARY KEY,
 domain_id uuid NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
 token_sha256 bytea NOT NULL CHECK (octet_length(token_sha256)=32),
 record_name text NOT NULL,
 expires_at timestamptz NOT NULL,
 consumed_at timestamptz,
 created_at timestamptz NOT NULL DEFAULT now(),
 CONSTRAINT challenge_expiry CHECK (expires_at > created_at)
);
CREATE INDEX domain_verification_pending_idx
 ON domain_verification_challenges(domain_id,expires_at)
 WHERE consumed_at IS NULL;
COMMIT;
