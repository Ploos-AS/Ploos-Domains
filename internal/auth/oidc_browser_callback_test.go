package auth

import (
 "context"
 "errors"
 "testing"

 "github.com/coreos/go-oidc/v3/oidc"
)

func TestBrowserIDTokenNonceAndIdentity(t *testing.T){
 token:=&oidc.IDToken{Issuer:"https://issuer.example",Subject:"subject-1",Nonce:"expected"}
 if err:=validateBrowserIDTokenClaims(token,"https://issuer.example","expected");err!=nil{t.Fatal(err)}
 for _,tc:=range []struct{name,nonce string}{
  {"wrong","wrong"},{"missing",""},
 }{
  t.Run(tc.name,func(t *testing.T){
   if err:=validateBrowserIDTokenClaims(token,"https://issuer.example",tc.nonce);!errors.Is(err,ErrOIDCNonceMismatch){t.Fatalf("got %v",err)}
  })
 }
 missingNonce:=&oidc.IDToken{Issuer:token.Issuer,Subject:token.Subject}
 if err:=validateBrowserIDTokenClaims(missingNonce,token.Issuer,"expected");!errors.Is(err,ErrOIDCNonceMismatch){t.Fatalf("got %v",err)}
 if err:=validateBrowserIDTokenClaims(token,"https://other.example","expected");!errors.Is(err,ErrInvalidIdentity){t.Fatalf("wrong issuer: %v",err)}
 if err:=validateBrowserIDTokenClaims(&oidc.IDToken{Issuer:token.Issuer,Nonce:"expected"},token.Issuer,"expected");!errors.Is(err,ErrInvalidIdentity){t.Fatalf("missing subject: %v",err)}
}
func TestBrowserVerifierFailsClosedWithoutConfiguration(t *testing.T){
 var verifier *OIDCVerifier
 if _,err:=verifier.VerifyBrowserIDToken(context.Background(),"token","nonce");!errors.Is(err,ErrUnauthenticated){t.Fatal(err)}
 verifier=&OIDCVerifier{}
 if _,err:=verifier.VerifyBrowserIDToken(context.Background(),"token","nonce");!errors.Is(err,ErrUnauthenticated){t.Fatal(err)}
}
