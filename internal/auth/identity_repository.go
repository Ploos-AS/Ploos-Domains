package auth

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
)

var ErrIdentityNotLinked = errors.New("OIDC identity is not linked to an internal user")

// IdentityRepository resolves a verified issuer+subject to an internal user.
// Linking and provisioning are separate privileged operations.
type IdentityRepository struct { DB *sql.DB }

func (r IdentityRepository) ResolveUserID(ctx context.Context, identity Identity) (string,error) {
 if r.DB==nil || identity.Validate()!=nil {return "",ErrInvalidIdentity}
 var id string
 err:=r.DB.QueryRowContext(ctx,`
  SELECT user_id::text FROM user_oidc_identities
  WHERE issuer=$1 AND subject=$2
 `,identity.Issuer,identity.Subject).Scan(&id)
 if errors.Is(err,sql.ErrNoRows){return "",ErrIdentityNotLinked}
 if err!=nil{return "",fmt.Errorf("resolve OIDC identity: %w",err)}
 return id,nil
}
