package auth

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

type AuditEventType string
const (
 AuditInvitationIssued AuditEventType = "invitation.issued"
 AuditInvitationRedeemed AuditEventType = "invitation.redeemed"
 AuditMembershipGranted AuditEventType = "membership.granted"
 AuditMembershipRevoked AuditEventType = "membership.revoked"
 AuditMembershipRoleChanged AuditEventType = "membership.role_changed"
)

var ErrInvalidAuditEvent = errors.New("invalid audit event")

type AuditEvent struct {
 OrganizationID string
 ActorUserID string // optional; blank for system actions
 TargetUserID string // optional
 Type AuditEventType
 OccurredAt time.Time
}

func (e AuditEvent) Validate() error {
 if e.OrganizationID=="" || e.OccurredAt.IsZero(){return ErrInvalidAuditEvent}
 switch e.Type {
 case AuditInvitationIssued,AuditInvitationRedeemed,AuditMembershipGranted,
  AuditMembershipRevoked,AuditMembershipRoleChanged:
  return nil
 default:return ErrInvalidAuditEvent
 }
}

type AuditRecorder struct { DB *sql.DB }

// RecordTx writes an event in the caller's transaction. If the write fails,
// the caller must roll back the security-sensitive operation as well.
func (a AuditRecorder) RecordTx(ctx context.Context,tx *sql.Tx,event AuditEvent) error {
 if a.DB==nil||tx==nil||event.Validate()!=nil{return ErrInvalidAuditEvent}
 id,err:=randomUUID()
 if err!=nil{return fmt.Errorf("audit event ID: %w",err)}
 var actor,target interface{}
 if event.ActorUserID!=""{actor=event.ActorUserID}
 if event.TargetUserID!=""{target=event.TargetUserID}
 _,err=tx.ExecContext(ctx,`
 INSERT INTO audit_events(id,organization_id,actor_user_id,target_user_id,event_type,occurred_at)
 VALUES($1,$2::uuid,$3::uuid,$4::uuid,$5,$6)
 `,id,event.OrganizationID,actor,target,string(event.Type),event.OccurredAt)
 if err!=nil{return fmt.Errorf("record audit event: %w",err)}
 return nil
}
