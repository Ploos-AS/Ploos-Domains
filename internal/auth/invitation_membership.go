package auth

import (
 "context"
 "crypto/sha256"
 "database/sql"
 "encoding/hex"
 "errors"
 "fmt"
 "time"
)

// InvitationMembership redeems an identity-bound invitation and grants only
// an organization-level viewer membership in the same database transaction.
// The caller must supply a cryptographically verified identity and its
// authenticated, resolved internal user ID.
type InvitationMembership struct { DB *sql.DB }

func (s InvitationMembership) Redeem(ctx context.Context, token string, identity Identity, userID string, now time.Time)(string,error){
 if s.DB==nil||identity.Validate()!=nil||userID==""||len(token)!=64{return "",ErrInvalidIdentity}
 if _,err:=hex.DecodeString(token);err!=nil{return "",ErrInvitationUnavailable}
 digest:=sha256.Sum256([]byte(token))
 tx,err:=s.DB.BeginTx(ctx,nil)
 if err!=nil{return "",fmt.Errorf("begin redemption: %w",err)}
 defer func(){_ = tx.Rollback()}()
 var orgID string
 err=tx.QueryRowContext(ctx,`
 SELECT e.organization_id::text FROM enrollment_invitations e
 JOIN user_oidc_identities i ON i.issuer=e.invited_issuer
  AND i.subject=e.invited_subject AND i.user_id=$4::uuid
 WHERE e.token_sha256=$1 AND e.invited_issuer=$2
  AND e.invited_subject=$3 AND e.consumed_at IS NULL
  AND e.expires_at>now() AND e.expires_at>$5
 FOR UPDATE OF e
 `,digest[:],identity.Issuer,identity.Subject,userID,now).Scan(&orgID)
 if errors.Is(err,sql.ErrNoRows){return "",ErrInvitationUnavailable}
 if err!=nil{return "",fmt.Errorf("lock invitation: %w",err)}
 membershipID,err:=randomUUID()
 if err!=nil{return "",fmt.Errorf("membership ID: %w",err)}
 var created string
 err=tx.QueryRowContext(ctx,`
 INSERT INTO memberships(id,user_id,organization_id,unit_id,role)
 VALUES($1,$2::uuid,$3::uuid,NULL,'viewer')
 ON CONFLICT (organization_id,user_id) WHERE unit_id IS NULL DO NOTHING
 RETURNING id::text
 `,membershipID,userID,orgID).Scan(&created)
 if errors.Is(err,sql.ErrNoRows){return "",ErrMembershipExists}
 if err!=nil{return "",fmt.Errorf("create viewer membership: %w",err)}
 result,err:=tx.ExecContext(ctx,`
 UPDATE enrollment_invitations SET consumed_at=now(),consumed_by=$2::uuid
 WHERE token_sha256=$1 AND consumed_at IS NULL AND expires_at>now()
 `,digest[:],userID)
 if err!=nil{return "",fmt.Errorf("consume invitation: %w",err)}
 n,err:=result.RowsAffected()
 if err!=nil{return "",err}
 if n!=1{return "",ErrInvitationUnavailable}
 if err=tx.Commit();err!=nil{return "",fmt.Errorf("commit redemption: %w",err)}
 return orgID,nil
}
