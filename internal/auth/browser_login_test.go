package auth

import (
 "net/url"
 "testing"
)

func TestBrowserLoginTransaction(t *testing.T){
 a,err:=NewBrowserLoginTransaction();if err!=nil{t.Fatal(err)}
 b,err:=NewBrowserLoginTransaction();if err!=nil{t.Fatal(err)}
 if a.State==a.Nonce||a.State==a.PKCEVerifier||a.State==b.State||len(a.PKCEVerifier)!=43{t.Fatal("non-independent login secrets")}
 raw,err:=a.AuthorizationURL("https://issuer.example/authorize","ploos-client","https://domains.example/callback")
 if err!=nil{t.Fatal(err)}
 u,err:=url.Parse(raw);if err!=nil{t.Fatal(err)}
 q:=u.Query()
 for key,want:=range map[string]string{"response_type":"code","scope":"openid","state":a.State,"nonce":a.Nonce,"code_challenge_method":"S256","code_challenge":a.PKCEChallenge()}{
  if q.Get(key)!=want{t.Fatalf("%s=%q want %q",key,q.Get(key),want)}
 }
}
func TestBrowserLoginRejectsUnsafeEndpoints(t *testing.T){
 tx,err:=NewBrowserLoginTransaction();if err!=nil{t.Fatal(err)}
 for _,endpoint:=range []string{"http://issuer.example/authorize","https://user:pass@issuer.example/auth","https://issuer.example/auth?evil=1","/authorize"}{
  if _,err:=tx.AuthorizationURL(endpoint,"client","https://domains.example/callback");err==nil{t.Fatalf("accepted endpoint %q",endpoint)}
 }
 if _,err:=tx.AuthorizationURL("https://issuer.example/authorize","client","http://domains.example/callback");err==nil{t.Fatal("accepted HTTP callback")}
}
