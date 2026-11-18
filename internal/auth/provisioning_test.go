package auth

import (
 "context"
 "errors"
 "testing"
)

func TestProvisioningRejectsMissingDatabase(t *testing.T){
 p:=UserProvisioner{}
 _,err:=p.Provision(context.Background(),Identity{Issuer:"https://id.example.org",Subject:"person",Audience:"ploos-domains"})
 if !errors.Is(err,ErrInvalidIdentity){t.Fatalf("got %v",err)}
}

func TestProvisioningRejectsUnverifiedClaims(t *testing.T){
 p:=UserProvisioner{}
 _,err:=p.Provision(context.Background(),Identity{})
 if !errors.Is(err,ErrInvalidIdentity){t.Fatalf("got %v",err)}
}
