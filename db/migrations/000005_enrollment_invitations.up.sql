BEGIN;
CREATE TABLE enrollment_invitations (
 id uuid PRIMARY KEY,
 token_sha256 bytea NOT NULL UNIQUE CHECK (octet_length(token_sha256)=32),
 organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 invited_issuer text NOT NULL CHECK (length(btrim(invited_issuer)) > 0),
 invited_subject text NOT NULL CHECK (length(btrim(invited_subject)) > 0),
 expires_at timestamptz NOT NULL,
 consumed_at timestamptz,
 consumed_by uuid REFERENCES users(id) ON DELETE SET NULL,
 created_at timestamptz NOT NULL DEFAULT now(),
 CONSTRAINT enrollment_invitation_expiry CHECK (expires_at > created_at),
 CONSTRAINT enrollment_consumption_pair CHECK ((consumed_at IS NULL) = (consumed_by IS NULL))
);
CREATE INDEX enrollment_invitations_pending_idx ON enrollment_invitations(expires_at) WHERE consumed_at IS NULL;
COMMIT;
