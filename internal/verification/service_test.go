package verification

import (
 "context"
 "crypto/sha256"
 "errors"
 "testing"
 "time"
)

type fakeAuth struct { allowed bool; calls int }
func (a *fakeAuth) CanVerifyDomain(_ context.Context,_,_,_ string)(bool,error){a.calls++;return a.allowed,nil}

type fakeRepo struct { challenge PendingChallenge; loads,completes int }
func (r *fakeRepo) LoadPending(_ context.Context,_,_ string)(PendingChallenge,error){r.loads++;return r.challenge,nil}
func (r *fakeRepo) Complete(_ context.Context,_,_ string,_ time.Time)error{r.completes++;return nil}

func TestVerificationService(t *testing.T){
 now:=time.Date(2026,11,18,12,0,0,0,time.UTC)
 token:="ploos-domain-verification=test-token"
 repo:=&fakeRepo{challenge:PendingChallenge{
  ID:"challenge",DomainID:"domain",OrganizationID:"org",
  RecordName:"_ploos-verify.example.org",TokenHash:sha256.Sum256([]byte(token)),
  ExpiresAt:now.Add(time.Hour),
 }}
 auth:=&fakeAuth{allowed:true}
 dns:=&fakeResolver{values:[]string{token}}
 service:=Service{Repository:repo,Authorizer:auth,Resolver:dns,Clock:func()time.Time{return now}}
 if err:=service.VerifyDomain(context.Background(),"actor","org","domain","challenge");err!=nil{t.Fatal(err)}
 if repo.completes!=1{t.Fatal("expected exactly one completion")}
 auth.allowed=false
 if err:=service.VerifyDomain(context.Background(),"actor","org","domain","challenge");!errors.Is(err,ErrUnauthorized){t.Fatalf("got %v",err)}
 if repo.loads!=1 || repo.completes!=1{t.Fatal("unauthorized actor reached persistence")}
 auth.allowed=true
 repo.challenge.OrganizationID="other"
 if err:=service.VerifyDomain(context.Background(),"actor","org","domain","challenge");!errors.Is(err,ErrChallengeUnavailable){t.Fatalf("got %v",err)}
 if repo.completes!=1{t.Fatal("cross-organization completion")}
 repo.challenge.OrganizationID="org"
 dns.values=[]string{"wrong"}
 if err:=service.VerifyDomain(context.Background(),"actor","org","domain","challenge");!errors.Is(err,ErrMismatch){t.Fatalf("got %v",err)}
 if repo.completes!=1{t.Fatal("mismatched DNS completed")}
}
