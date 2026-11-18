package auth

import (
 "context"
 "errors"
 "testing"
 "time"
)

func TestInvitationToken(t *testing.T){
 a,ha,err:=NewInvitationToken()
 if err!=nil{t.Fatal(err)}
 b,hb,err:=NewInvitationToken()
 if err!=nil{t.Fatal(err)}
 if len(a)!=64||len(b)!=64||a==b||ha==hb{t.Fatal("invalid invitation token entropy")}
}
func TestInvitationStoreRejectsMissingDB(t *testing.T){
 _,err:=(InvitationStore{}).ConsumeInvitation(context.Background(),"invalid",Identity{},"",time.Now())
 if !errors.Is(err,ErrInvalidIdentity){t.Fatal(err)}
}
