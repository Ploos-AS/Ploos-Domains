-- Initial organization schema for PostgreSQL 15+.
BEGIN;

CREATE TABLE organizations (
 id uuid PRIMARY KEY,
 name text NOT NULL CHECK (length(btrim(name)) BETWEEN 1 AND 200),
 created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE users (
 id uuid PRIMARY KEY,
 external_subject text NOT NULL UNIQUE CHECK (length(btrim(external_subject)) > 0),
 created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE units (
 id uuid PRIMARY KEY,
 organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 parent_id uuid,
 name text NOT NULL CHECK (length(btrim(name)) BETWEEN 1 AND 200),
 created_at timestamptz NOT NULL DEFAULT now(),
 UNIQUE (organization_id, id),
 CONSTRAINT units_not_self_parent CHECK (parent_id IS NULL OR parent_id <> id),
 CONSTRAINT units_parent_same_org FOREIGN KEY (organization_id, parent_id)
  REFERENCES units (organization_id, id) DEFERRABLE INITIALLY IMMEDIATE
);

CREATE INDEX units_parent_idx ON units (organization_id, parent_id);

CREATE TABLE memberships (
 id uuid PRIMARY KEY,
 user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 unit_id uuid,
 role text NOT NULL CHECK (role IN ('owner', 'admin', 'unit_admin', 'viewer')),
 created_at timestamptz NOT NULL DEFAULT now(),
 CONSTRAINT memberships_scope_role CHECK (
  (unit_id IS NULL AND role IN ('owner', 'admin', 'viewer'))
  OR (unit_id IS NOT NULL AND role IN ('unit_admin', 'viewer'))
 ),
 CONSTRAINT memberships_unit_same_org FOREIGN KEY (organization_id, unit_id)
  REFERENCES units (organization_id, id) ON DELETE CASCADE
);

-- NULL unit_id needs a separate unique index to avoid duplicate org-level memberships.
CREATE UNIQUE INDEX memberships_org_user_unique
 ON memberships (organization_id, user_id) WHERE unit_id IS NULL;
CREATE UNIQUE INDEX memberships_unit_user_unique
 ON memberships (organization_id, unit_id, user_id) WHERE unit_id IS NOT NULL;
CREATE INDEX memberships_user_idx ON memberships (user_id);

COMMIT;
