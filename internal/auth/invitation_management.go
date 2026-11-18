package auth

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

var ErrInvitationForbidden = errors.New("not authorized to manage invitations")

// InvitationAuthorizer checks trusted internal actor IDs against an organization.
type InvitationAuthorizer interface {
 CanManageInvitations(ctx context.Context, actorID, organizationID string) (bool,error)
}

// PostgresInvitationAuthorizer permits organization-scoped owners and admins.
// Unit administrators and viewers cannot issue invitations for the organization.
type PostgresInvitationAuthorizer struct { DB *sql.DB }
func (a PostgresInvitationAuthorizer) CanManageInvitations(ctx context.Context, actorID, organizationID string)(bool,error){
 if a.DB==nil || actorID=="" || organizationID=="" {return false,ErrInvalidIdentity}
 var allowed bool
 err:=a.DB.QueryRowContext(ctx,`
 SELECT EXISTS(SELECT 1 FROM memberships
 WHERE user_id=$1::uuid AND organization_id=$2::uuid
 AND unit_id IS NULL AND role IN ('owner','admin'))`,actorID,organizationID).Scan(&allowed)
 if err!=nil{return false,fmt.Errorf("authorize invitation management: %w",err)}
 return allowed,nil
}

type InvitationWriter interface {
 Create(ctx context.Context, organizationID, issuer, subject string, digest [32]byte, created, expires time.Time) error
}

type PostgresInvitationWriter struct { DB *sql.DB }
func (w PostgresInvitationWriter) Create(ctx context.Context, organizationID, issuer, subject string, digest [32]byte, created, expires time.Time) error {
 if w.DB==nil {return ErrInvalidIdentity}
 id,err:=randomUUID()
 if err!=nil{return fmt.Errorf("generate invitation ID: %w",err)}
 _,err=w.DB.ExecContext(ctx,`
 INSERT INTO enrollment_invitations
 (id,token_sha256,organization_id,invited_issuer,invited_subject,created_at,expires_at)
 VALUES ($1,$2,$3::uuid,$4,$5,$6,$7)`,id,digest[:],organizationID,issuer,subject,created,expires)
 if err!=nil{return fmt.Errorf("create invitation: %w",err)}
 return nil
}

type InvitationManager struct {
 Authorizer InvitationAuthorizer
 Writer InvitationWriter
 Clock func()time.Time
}

// Issue returns the raw invitation token exactly once. Callers must transmit
// it securely and never log it. No membership or role is granted.
func (m InvitationManager) Issue(ctx context.Context, actorID, organizationID, issuer, subject string, ttl time.Duration)(string,error){
 if actorID=="" || organizationID=="" || issuer=="" || subject=="" ||
  ttl<=0 || ttl>24*time.Hour || m.Authorizer==nil || m.Writer==nil || m.Clock==nil {
  return "",ErrInvalidIdentity
 }
 if err:=ValidateOIDCConfig(issuer,"invitation");err!=nil{return "",ErrInvalidIdentity}
 allowed,err:=m.Authorizer.CanManageInvitations(ctx,actorID,organizationID)
 if err!=nil{return "",fmt.Errorf("authorize invitation issuance: %w",err)}
 if !allowed{return "",ErrInvitationForbidden}
 token,digest,err:=NewInvitationToken()
 if err!=nil{return "",fmt.Errorf("generate invitation: %w",err)}
 now:=m.Clock()
 if err:=m.Writer.Create(ctx,organizationID,issuer,subject,digest,now,now.Add(ttl));err!=nil{
  return "",fmt.Errorf("persist invitation: %w",err)
 }
 return token,nil
}
