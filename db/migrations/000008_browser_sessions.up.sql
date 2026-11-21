BEGIN;
CREATE TABLE browser_sessions (
 token_sha256 bytea PRIMARY KEY CHECK (octet_length(token_sha256)=32),
 user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 created_at timestamptz NOT NULL DEFAULT now(),
 expires_at timestamptz NOT NULL,
 revoked_at timestamptz,
 CONSTRAINT browser_sessions_expiry CHECK (expires_at>created_at),
 CONSTRAINT browser_sessions_revocation CHECK (revoked_at IS NULL OR revoked_at>=created_at)
);
CREATE INDEX browser_sessions_user_idx ON browser_sessions(user_id);
CREATE INDEX browser_sessions_expiry_idx ON browser_sessions(expires_at);
COMMIT;
