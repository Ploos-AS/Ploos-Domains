package auth

import (
 "context"
 "crypto/rand"
 "crypto/sha256"
 "database/sql"
 "encoding/hex"
 "errors"
 "fmt"
 "time"
)

var ErrInvitationUnavailable = errors.New("invitation invalid, expired, consumed or not intended for this identity")

// NewInvitationToken returns a 256-bit random bearer secret. Persist only its
// digest; deliver the token through an authenticated, out-of-band workflow.
func NewInvitationToken() (string,[32]byte,error) {
 var secret [32]byte
 if _,err:=rand.Read(secret[:]);err!=nil{return "",[32]byte{},err}
 token:=hex.EncodeToString(secret[:])
 return token,sha256.Sum256([]byte(token)),nil
}

type InvitationStore struct{ DB *sql.DB }

// ConsumeInvitation checks a previously authenticated internal user and
// verified OIDC identity against an invitation, then consumes it once.
// It intentionally does not grant any membership or role.
func (s InvitationStore) ConsumeInvitation(ctx context.Context, token string, identity Identity, userID string, now time.Time) (string,error) {
 if s.DB==nil || identity.Validate()!=nil || userID=="" || len(token)!=64 {return "",ErrInvalidIdentity}
 if _,err:=hex.DecodeString(token);err!=nil{return "",ErrInvitationUnavailable}
 digest:=sha256.Sum256([]byte(token))
 tx,err:=s.DB.BeginTx(ctx,nil)
 if err!=nil{return "",fmt.Errorf("begin enrollment: %w",err)}
 defer func(){_ = tx.Rollback()}()
 var orgID string
 // Only accept the invitation for an identity already linked to this user.
 err=tx.QueryRowContext(ctx,`
 SELECT e.organization_id::text
 FROM enrollment_invitations e
 JOIN user_oidc_identities i ON i.issuer=e.invited_issuer
   AND i.subject=e.invited_subject
 WHERE e.token_sha256=$1 AND e.invited_issuer=$2
   AND e.invited_subject=$3 AND i.user_id=$4::uuid
   AND e.consumed_at IS NULL AND e.expires_at>$5
 FOR UPDATE OF e
 `,digest[:],identity.Issuer,identity.Subject,userID,now).Scan(&orgID)
 if errors.Is(err,sql.ErrNoRows){return "",ErrInvitationUnavailable}
 if err!=nil{return "",fmt.Errorf("lock invitation: %w",err)}
 result,err:=tx.ExecContext(ctx,`
 UPDATE enrollment_invitations SET consumed_at=$2,consumed_by=$3::uuid
 WHERE token_sha256=$1 AND consumed_at IS NULL AND expires_at>$2
 `,digest[:],now,userID)
 if err!=nil{return "",fmt.Errorf("consume invitation: %w",err)}
 n,err:=result.RowsAffected()
 if err!=nil{return "",err}
 if n!=1{return "",ErrInvitationUnavailable}
 if err=tx.Commit();err!=nil{return "",fmt.Errorf("commit invitation: %w",err)}
 return orgID,nil
}
