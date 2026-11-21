package auth

import (
 "crypto/rand"
 "crypto/sha256"
 "encoding/base64"
 "errors"
 "net/url"
 "strings"
)

var ErrInvalidLoginTransaction=errors.New("invalid OIDC login transaction")

type BrowserLoginTransaction struct{
 State string
 Nonce string
 PKCEVerifier string
}

// NewBrowserLoginTransaction generates independent high-entropy state,
// nonce and PKCE verifier values. The server must store them securely with
// a short TTL and bind them to the initiating browser before redirecting.
func NewBrowserLoginTransaction()(BrowserLoginTransaction,error){
 var parts [3]string
 for i:=range parts{
  var b [32]byte
  if _,err:=rand.Read(b[:]);err!=nil{return BrowserLoginTransaction{},err}
  parts[i]=base64.RawURLEncoding.EncodeToString(b[:])
 }
 return BrowserLoginTransaction{State:parts[0],Nonce:parts[1],PKCEVerifier:parts[2]},nil
}

func (t BrowserLoginTransaction) PKCEChallenge()string{
 sum:=sha256.Sum256([]byte(t.PKCEVerifier))
 return base64.RawURLEncoding.EncodeToString(sum[:])
}

// AuthorizationURL builds an OIDC authorization-code redirect using only a
// trusted configured authorization endpoint. It does not exchange tokens.
func (t BrowserLoginTransaction) AuthorizationURL(endpoint,clientID,redirectURI string)(string,error){
 if t.State==""||t.Nonce==""||len(t.PKCEVerifier)<43||clientID==""||redirectURI==""{
  return "",ErrInvalidLoginTransaction
 }
 u,err:=url.Parse(endpoint)
 if err!=nil||u.Scheme!="https"||u.Host==""||u.User!=nil||u.Fragment!=""||u.RawQuery!=""{
  return "",ErrInvalidLoginTransaction
 }
 callback,err:=url.Parse(redirectURI)
 if err!=nil||callback.Scheme!="https"||callback.Host==""||callback.User!=nil||callback.Fragment!=""{
  return "",ErrInvalidLoginTransaction
 }
 if strings.TrimSpace(clientID)!=clientID{return "",ErrInvalidLoginTransaction}
 q:=url.Values{}
 q.Set("response_type","code")
 q.Set("client_id",clientID)
 q.Set("redirect_uri",redirectURI)
 q.Set("scope","openid")
 q.Set("state",t.State)
 q.Set("nonce",t.Nonce)
 q.Set("code_challenge",t.PKCEChallenge())
 q.Set("code_challenge_method","S256")
 u.RawQuery=q.Encode()
 return u.String(),nil
}
