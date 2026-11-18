package verification

import (
 "context"
 "crypto/sha256"
 "errors"
 "testing"
 "time"
)

type fakeResolver struct { values []string; err error; queried string }
func (f *fakeResolver) LookupTXT(_ context.Context,name string)([]string,error){ f.queried=name;return f.values,f.err }

func TestChallengeAndVerification(t *testing.T){
 now:=time.Date(2026,11,18,0,0,0,0,time.UTC)
 ch,err:=NewChallenge("Seniornett.NO.",now,15*time.Minute)
 if err!=nil {t.Fatal(err)}
 if ch.RecordName!="_ploos-verify.seniornett.no" {t.Fatal(ch.RecordName)}
 if ch.Token=="" || ch.TokenHash!=sha256.Sum256([]byte(ch.Token)){t.Fatal("invalid token")}
 resolver:=&fakeResolver{values:[]string{"other",ch.Token}}
 if err:=Verify(context.Background(),resolver,ch.RecordName,ch.TokenHash,ch.ExpiresAt,now);err!=nil {t.Fatal(err)}
 if resolver.queried!=ch.RecordName {t.Fatal(resolver.queried)}
 resolver.values=[]string{"wrong"}
 if err:=Verify(context.Background(),resolver,ch.RecordName,ch.TokenHash,ch.ExpiresAt,now);!errors.Is(err,ErrMismatch){t.Fatalf("got %v",err)}
 resolver.values=[]string{ch.Token}
 if err:=Verify(context.Background(),resolver,ch.RecordName,ch.TokenHash,ch.ExpiresAt,ch.ExpiresAt);!errors.Is(err,ErrExpired){t.Fatalf("got %v",err)}
}

func TestChallengeValidation(t *testing.T){
 now:=time.Now()
 for _,domain:=range []string{"", "*.example.org","bad..example.org"}{
  if _,err:=NewChallenge(domain,now,time.Hour);!errors.Is(err,ErrInvalid){t.Fatalf("accepted %q",domain)}
 }
 if _,err:=NewChallenge("example.org",now,25*time.Hour);!errors.Is(err,ErrInvalid){t.Fatal(err)}
 if _,err:=NewChallenge("example.org",now,0);!errors.Is(err,ErrInvalid){t.Fatal(err)}
}

func TestResolverFailure(t *testing.T){
 now:=time.Now()
 ch,err:=NewChallenge("example.org",now,time.Hour);if err!=nil{t.Fatal(err)}
 resolver:=&fakeResolver{err:errors.New("DNS timeout")}
 if err:=Verify(context.Background(),resolver,ch.RecordName,ch.TokenHash,ch.ExpiresAt,now);err==nil{t.Fatal("expected error")}
}
