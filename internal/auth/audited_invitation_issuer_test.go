package auth

import (
 "context"
 "errors"
 "testing"
 "time"
)

func TestAuditedInvitationIssuerRejectsBadInput(t *testing.T){
 s:=AuditedInvitationIssuer{}
 _,err:=s.Issue(context.Background(),"actor","org","https://id.example.org","subject",time.Hour)
 if !errors.Is(err,ErrInvalidIdentity){t.Fatal(err)}
 s.Clock=time.Now
 _,err=s.Issue(context.Background(),"actor","org","https://id.example.org","subject",25*time.Hour)
 if !errors.Is(err,ErrInvalidIdentity){t.Fatal(err)}
}
