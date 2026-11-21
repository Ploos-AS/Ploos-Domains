BEGIN;
-- Independent of the existing verified domains inventory: customers may track
-- domains held at external registrars without claiming DNS control.
CREATE TABLE domain_portfolio (
 organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 fqdn text NOT NULL,
 registrar text NOT NULL,
 status text NOT NULL DEFAULT 'unknown',
 expires_at timestamptz,
 last_verified_at timestamptz,
 created_at timestamptz NOT NULL DEFAULT now(),
 updated_at timestamptz NOT NULL DEFAULT now(),
 PRIMARY KEY (organization_id,fqdn),
 CONSTRAINT portfolio_domain_canonical CHECK (
  length(fqdn) BETWEEN 3 AND 253
  AND fqdn = lower(fqdn)
  AND fqdn ~ '^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$'
  AND fqdn LIKE '%.%'
  AND fqdn !~ '[.][.]'
 ),
 CONSTRAINT portfolio_registrar_nonempty CHECK (
  length(btrim(registrar)) BETWEEN 1 AND 200 AND registrar=btrim(registrar)
 ),
 CONSTRAINT portfolio_status_valid CHECK (
  status IN ('unknown','active','pending','expired','transferred')
 )
);
CREATE INDEX domain_portfolio_org_expiry_idx ON domain_portfolio(organization_id,expires_at)
 WHERE expires_at IS NOT NULL;
COMMIT;
