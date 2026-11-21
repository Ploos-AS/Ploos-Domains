package auth

import (
 "context"
 "crypto/subtle"
 "errors"
 "fmt"

 "github.com/coreos/go-oidc/v3/oidc"
)

var ErrOIDCNonceMismatch=errors.New("OIDC nonce mismatch")

// VerifyBrowserIDToken validates an ID token's signature, expiry, issuer,
// audience and algorithm through the configured OIDC verifier, then enforces
// the nonce from a previously consumed, browser-bound login transaction.
// Never use this method with an untrusted caller-supplied expectedNonce.
func (v *OIDCVerifier) VerifyBrowserIDToken(ctx context.Context,rawToken,expectedNonce string)(Identity,error){
 if v==nil||v.verifier==nil||rawToken==""||expectedNonce==""{
  return Identity{},ErrUnauthenticated
 }
 token,err:=v.verifier.Verify(ctx,rawToken)
 if err!=nil{return Identity{},fmt.Errorf("%w: %v",ErrUnauthenticated,err)}
 if err:=validateBrowserIDTokenClaims(token,v.issuer,expectedNonce);err!=nil{return Identity{},err}
 return Identity{Issuer:token.Issuer,Subject:token.Subject,Audience:v.audience},nil
}

func validateBrowserIDTokenClaims(token *oidc.IDToken,issuer,expectedNonce string)error{
 if token==nil||token.Subject==""||token.Issuer!=issuer{return ErrInvalidIdentity}
 if token.Nonce==""||expectedNonce==""||subtle.ConstantTimeCompare([]byte(token.Nonce),[]byte(expectedNonce))!=1{
  return ErrOIDCNonceMismatch
 }
 return nil
}
