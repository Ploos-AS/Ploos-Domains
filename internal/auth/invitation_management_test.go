package auth

import (
 "context"
 "errors"
 "testing"
 "time"
)

type invitationAuthStub struct{ allowed bool; calls int }
func (a *invitationAuthStub) CanManageInvitations(_ context.Context,_,_ string)(bool,error){a.calls++;return a.allowed,nil}
type invitationWriterStub struct{ calls int; digest [32]byte }
func (w *invitationWriterStub) Create(_ context.Context,_,_,_ string,digest [32]byte,_,_ time.Time)error{
 w.calls++;w.digest=digest;return nil
}
func TestInvitationManager(t *testing.T){
 now:=time.Date(2026,11,18,12,0,0,0,time.UTC)
 a:=&invitationAuthStub{allowed:true}
 w:=&invitationWriterStub{}
 m:=InvitationManager{Authorizer:a,Writer:w,Clock:func()time.Time{return now}}
 token,err:=m.Issue(context.Background(),"actor","org","https://id.example.org","subject",time.Hour)
 if err!=nil||len(token)!=64||w.calls!=1{t.Fatalf("token=%q calls=%d err=%v",token,w.calls,err)}
 a.allowed=false
 _,err=m.Issue(context.Background(),"actor","org","https://id.example.org","subject",time.Hour)
 if !errors.Is(err,ErrInvitationForbidden)||w.calls!=1{t.Fatalf("denial err=%v writes=%d",err,w.calls)}
}
func TestInvitationManagerRejectsBadInput(t *testing.T){
 a:=&invitationAuthStub{allowed:true}
 w:=&invitationWriterStub{}
 m:=InvitationManager{Authorizer:a,Writer:w,Clock:time.Now}
 for _,tc:=range []struct{issuer,subject string;ttl time.Duration}{
  {"http://id.example.org","person",time.Hour},
  {"https://id.example.org","",time.Hour},
  {"https://id.example.org","person",0},
  {"https://id.example.org","person",25*time.Hour},
 }{
  _,err:=m.Issue(context.Background(),"actor","org",tc.issuer,tc.subject,tc.ttl)
  if !errors.Is(err,ErrInvalidIdentity){t.Fatalf("unexpected err: %v",err)}
 }
 if a.calls!=0||w.calls!=0{t.Fatal("invalid input reached authorization or persistence")}
}
