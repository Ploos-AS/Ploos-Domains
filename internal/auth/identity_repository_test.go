package auth

import (
 "context"
 "errors"
 "testing"
)

func TestIdentityRepositoryRejectsMissingDB(t *testing.T) {
 repo:=IdentityRepository{}
 _,err:=repo.ResolveUserID(context.Background(),Identity{
  Issuer:"https://id.example.org",Subject:"user-1",Audience:"ploos-domains",
 })
 if !errors.Is(err,ErrInvalidIdentity){t.Fatalf("got %v",err)}
}
