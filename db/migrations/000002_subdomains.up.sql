BEGIN;

-- Inventory only: verification and DNS-provider adapters are later work.
CREATE TABLE domains (
 id uuid PRIMARY KEY,
 organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 fqdn text NOT NULL CHECK (fqdn = lower(fqdn) AND fqdn !~ '[.]$'),
 verified_at timestamptz,
 created_at timestamptz NOT NULL DEFAULT now(),
 UNIQUE (organization_id, id),
 UNIQUE (fqdn),
 CONSTRAINT domains_fqdn_basic CHECK (
  length(fqdn) BETWEEN 3 AND 253
  AND fqdn ~ '^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$'
  AND fqdn !~ '[.][.]'
 )
);

CREATE TABLE subdomain_allocations (
 id uuid PRIMARY KEY,
 organization_id uuid NOT NULL,
 domain_id uuid NOT NULL,
 unit_id uuid NOT NULL,
 fqdn text NOT NULL,
 created_at timestamptz NOT NULL DEFAULT now(),
 CONSTRAINT allocations_domain_same_org FOREIGN KEY (organization_id,domain_id)
  REFERENCES domains(organization_id,id) ON DELETE CASCADE,
 CONSTRAINT allocations_unit_same_org FOREIGN KEY (organization_id,unit_id)
  REFERENCES units(organization_id,id) ON DELETE CASCADE,
 CONSTRAINT allocations_name_unique UNIQUE (fqdn),
 CONSTRAINT allocations_name_canonical CHECK (fqdn = lower(fqdn) AND fqdn !~ '[.]$'),
 CONSTRAINT allocations_name_basic CHECK (
  length(fqdn) BETWEEN 3 AND 253
  AND fqdn ~ '^[a-z0-9]([a-z0-9.-]*[a-z0-9])?$'
  AND fqdn !~ '[.][.]'
 )
);
CREATE INDEX subdomain_allocations_unit_idx ON subdomain_allocations(organization_id,unit_id);
CREATE INDEX subdomain_allocations_domain_idx ON subdomain_allocations(domain_id);

-- Prevent allocations outside the domain and before verification, including
-- writes that bypass application code. Domain rename/verification changes
-- are checked separately; domain FQDN is immutable after creation.
CREATE FUNCTION check_subdomain_allocation() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE parent_name text; verified timestamptz;
BEGIN
 SELECT d.fqdn,d.verified_at INTO parent_name,verified
 FROM domains d WHERE d.id=NEW.domain_id AND d.organization_id=NEW.organization_id
 FOR SHARE;
 IF NOT FOUND OR verified IS NULL THEN
  RAISE EXCEPTION 'domain not verified or not found' USING ERRCODE='23514';
 END IF;
 IF length(NEW.fqdn) <= length(parent_name)+1 OR
    right(NEW.fqdn,length(parent_name)+1) <> '.' || parent_name THEN
  RAISE EXCEPTION 'allocation must be below verified domain' USING ERRCODE='23514';
 END IF;
 RETURN NEW;
END $$;

CREATE TRIGGER subdomain_allocation_guard
 BEFORE INSERT OR UPDATE ON subdomain_allocations
 FOR EACH ROW EXECUTE FUNCTION check_subdomain_allocation();

CREATE FUNCTION prevent_domain_identity_change() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
 IF NEW.fqdn IS DISTINCT FROM OLD.fqdn OR
    NEW.organization_id IS DISTINCT FROM OLD.organization_id OR
    (OLD.verified_at IS NOT NULL AND NEW.verified_at IS NULL) THEN
  RAISE EXCEPTION 'domain identity and verification cannot be changed in place' USING ERRCODE='23514';
 END IF;
 RETURN NEW;
END $$;
CREATE TRIGGER domain_identity_guard
 BEFORE UPDATE ON domains
 FOR EACH ROW EXECUTE FUNCTION prevent_domain_identity_change();

COMMIT;
