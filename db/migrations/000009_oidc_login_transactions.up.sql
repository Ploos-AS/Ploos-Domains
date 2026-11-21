BEGIN;
CREATE TABLE oidc_login_transactions (
 state_sha256 bytea PRIMARY KEY CHECK (octet_length(state_sha256)=32),
 browser_binding_sha256 bytea NOT NULL CHECK (octet_length(browser_binding_sha256)=32),
 nonce text NOT NULL,
 pkce_verifier text NOT NULL,
 created_at timestamptz NOT NULL DEFAULT now(),
 expires_at timestamptz NOT NULL,
 consumed_at timestamptz,
 CONSTRAINT oidc_login_transaction_expiry CHECK (expires_at>created_at)
);
CREATE INDEX oidc_login_transactions_expiry_idx ON oidc_login_transactions(expires_at);
COMMIT;
