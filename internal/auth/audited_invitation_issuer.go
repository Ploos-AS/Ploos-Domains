package auth

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

// AuditedInvitationIssuer is the preferred PostgreSQL issuance path. Unlike
// the legacy InvitationManager, authorization, insertion and audit are atomic.
type AuditedInvitationIssuer struct { DB *sql.DB; Clock func()time.Time }

func (s AuditedInvitationIssuer) Issue(ctx context.Context, actorID, organizationID, issuer, subject string, ttl time.Duration)(string,error){
 if s.DB==nil||s.Clock==nil||actorID==""||organizationID==""||issuer==""||subject==""||ttl<=0||ttl>24*time.Hour{
  return "",ErrInvalidIdentity
 }
 if err:=ValidateOIDCConfig(issuer,"invitation");err!=nil{return "",ErrInvalidIdentity}
 tx,err:=s.DB.BeginTx(ctx,nil)
 if err!=nil{return "",fmt.Errorf("begin audited issuance: %w",err)}
 defer func(){_ = tx.Rollback()}()
 var role string
 err=tx.QueryRowContext(ctx,`
 SELECT role FROM memberships WHERE user_id=$1::uuid AND organization_id=$2::uuid
 AND unit_id IS NULL AND role IN ('owner','admin') FOR SHARE
 `,actorID,organizationID).Scan(&role)
 if errors.Is(err,sql.ErrNoRows){return "",ErrInvitationForbidden}
 if err!=nil{return "",fmt.Errorf("authorize audited issuance: %w",err)}
 token,digest,err:=NewInvitationToken()
 if err!=nil{return "",fmt.Errorf("generate invitation: %w",err)}
 id,err:=randomUUID()
 if err!=nil{return "",fmt.Errorf("generate invitation ID: %w",err)}
 now:=s.Clock()
 _,err=tx.ExecContext(ctx,`
 INSERT INTO enrollment_invitations
 (id,token_sha256,organization_id,invited_issuer,invited_subject,created_at,expires_at)
 VALUES($1,$2,$3::uuid,$4,$5,$6,$7)
 `,id,digest[:],organizationID,issuer,subject,now,now.Add(ttl))
 if err!=nil{return "",fmt.Errorf("persist invitation: %w",err)}
 if err:=(AuditRecorder{DB:s.DB}).RecordTx(ctx,tx,AuditEvent{
  OrganizationID:organizationID,ActorUserID:actorID,
  Type:AuditInvitationIssued,OccurredAt:now,
 });err!=nil{return "",fmt.Errorf("audit invitation issuance: %w",err)}
 if err:=tx.Commit();err!=nil{return "",fmt.Errorf("commit audited issuance: %w",err)}
 return token,nil
}
