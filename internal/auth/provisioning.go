package auth

import (
 "context"
 "crypto/rand"
 "database/sql"
 "encoding/hex"
 "errors"
 "fmt"
)

// UserProvisioner creates an unprivileged user and binds a verified OIDC
// issuer+subject. Only call it after cryptographic token verification and
// an explicit enrollment policy decision; never trust raw client claims.
type UserProvisioner struct { DB *sql.DB }

func randomUUID() (string,error) {
 b:=make([]byte,16)
 if _,err:=rand.Read(b);err!=nil{return "",err}
 b[6]=(b[6]&0x0f)|0x40
 b[8]=(b[8]&0x3f)|0x80
 h:=hex.EncodeToString(b)
 return fmt.Sprintf("%s-%s-%s-%s-%s",h[:8],h[8:12],h[12:16],h[16:20],h[20:]),nil
}

// Provision returns the existing internal user for a verified identity, or
// atomically creates a new user and link. No memberships or roles are added.
// A unique (issuer,subject) constraint serializes concurrent registrations.
func (p UserProvisioner) Provision(ctx context.Context,identity Identity)(userID string,err error){
 if p.DB==nil || identity.Validate()!=nil{return "",ErrInvalidIdentity}
 newID,err:=randomUUID()
 if err!=nil{return "",fmt.Errorf("generate user ID: %w",err)}
 tx,err:=p.DB.BeginTx(ctx,nil)
 if err!=nil{return "",fmt.Errorf("begin provisioning: %w",err)}
 defer func(){_ = tx.Rollback()}()

 // Preserve the legacy NOT NULL unique external_subject until it can be
 // removed in a separately reviewed schema migration. Never use an OIDC
 // subject alone as this legacy identifier.
 _,err=tx.ExecContext(ctx,
  "INSERT INTO users(id,external_subject) VALUES($1,$2)",
  newID,"oidc-provisioned:"+newID)
 if err!=nil{return "",fmt.Errorf("create user: %w",err)}

 var linkedID string
 err=tx.QueryRowContext(ctx,`
  INSERT INTO user_oidc_identities(issuer,subject,user_id)
  VALUES($1,$2,$3)
  ON CONFLICT (issuer,subject) DO UPDATE
    SET subject=EXCLUDED.subject
  RETURNING user_id::text
 `,identity.Issuer,identity.Subject,newID).Scan(&linkedID)
 if err!=nil{return "",fmt.Errorf("link OIDC identity: %w",err)}

 if linkedID!=newID {
  // Existing identity won the race. Do not retain an orphan user.
  if _,err=tx.ExecContext(ctx,"DELETE FROM users WHERE id=$1",newID);err!=nil{
   return "",fmt.Errorf("discard duplicate user: %w",err)
  }
 }
 if err=tx.Commit();err!=nil{return "",fmt.Errorf("commit provisioning: %w",err)}
 if linkedID==""{return "",errors.New("provisioning returned empty user ID")}
 return linkedID,nil
}
