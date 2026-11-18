BEGIN;
CREATE TABLE audit_events (
 id uuid PRIMARY KEY,
 organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE RESTRICT,
 actor_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
 target_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
 event_type text NOT NULL CHECK (event_type IN (
  'invitation.issued','invitation.redeemed','membership.granted',
  'membership.revoked','membership.role_changed'
 )),
 occurred_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX audit_events_org_time_idx ON audit_events(organization_id,occurred_at DESC,id);
CREATE INDEX audit_events_actor_time_idx ON audit_events(actor_user_id,occurred_at DESC)
 WHERE actor_user_id IS NOT NULL;
COMMIT;
