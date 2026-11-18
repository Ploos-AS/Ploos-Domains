package auth

import (
 "context"
 "errors"
 "testing"
 "time"
)

func TestInvitationMembershipRejectsInvalidInput(t *testing.T){
 s:=InvitationMembership{}
 _,err:=s.Redeem(context.Background(),"",Identity{},"",time.Now())
 if !errors.Is(err,ErrInvalidIdentity){t.Fatal(err)}
}
