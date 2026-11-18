package auth

import (
 "context"
 "errors"
 "testing"
)

func TestOIDCConfig(t *testing.T){
 for _,tc:=range []struct{issuer,audience string;valid bool}{
  {"https://auth.example.org","ploos-domains",true},
  {"https://auth.example.org/realms/ploos","ploos-domains",true},
  {"http://auth.example.org","ploos-domains",false},
  {"https://auth.example.org/","ploos-domains",false},
  {"https://user:pass@auth.example.org","ploos-domains",false},
  {"https://auth.example.org?x=1","ploos-domains",false},
  {"https://auth.example.org#frag","ploos-domains",false},
  {"","ploos-domains",false},
  {"https://auth.example.org","",false},
 }{
  err:=ValidateOIDCConfig(tc.issuer,tc.audience)
  if (err==nil)!=tc.valid{t.Errorf("issuer=%q audience=%q err=%v",tc.issuer,tc.audience,err)}
 }
}

func TestUnconfiguredVerifierRejectsTokens(t *testing.T){
 var verifier *OIDCVerifier
 if _,err:=verifier.VerifyToken(context.Background(),"header.payload.signature");!errors.Is(err,ErrUnauthenticated){t.Fatal(err)}
 verifier=&OIDCVerifier{}
 if _,err:=verifier.VerifyToken(context.Background(),"header.payload.signature");!errors.Is(err,ErrUnauthenticated){t.Fatal(err)}
}
