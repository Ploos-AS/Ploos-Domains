package auth

import (
 "crypto/rand"
 "encoding/base64"
 "errors"
 "net/http"
 "time"
)

const BrowserLoginBindingCookie = "__Host-ploos_oidc_binding"

var ErrInvalidBrowserLoginBinding = errors.New("invalid OIDC browser binding")

// NewBrowserLoginBinding generates a secret independent of state, nonce and
// PKCE. Store only its digest with the login transaction.
func NewBrowserLoginBinding()(string,error){
 var raw [32]byte
 if _,err:=rand.Read(raw[:]);err!=nil{return "",err}
 return base64.RawURLEncoding.EncodeToString(raw[:]),nil
}

// SetBrowserLoginBindingCookie is intended for HTTPS login-start handlers.
// SameSite=Lax permits top-level GET redirects back from the OIDC provider.
// A callback must clear this cookie after consuming the transaction.
func SetBrowserLoginBindingCookie(w http.ResponseWriter,binding string,ttl time.Duration)error{
 if w==nil||len(binding)!=43||ttl<time.Minute||ttl>10*time.Minute{
  return ErrInvalidBrowserLoginBinding
 }
 if _,err:=base64.RawURLEncoding.DecodeString(binding);err!=nil{return ErrInvalidBrowserLoginBinding}
 http.SetCookie(w,&http.Cookie{
  Name:BrowserLoginBindingCookie,Value:binding,Path:"/",Secure:true,
  HttpOnly:true,SameSite:http.SameSiteLaxMode,
  MaxAge:int(ttl.Seconds()),Expires:time.Now().Add(ttl),
 })
 return nil
}

func ClearBrowserLoginBindingCookie(w http.ResponseWriter){
 if w==nil{return}
 http.SetCookie(w,&http.Cookie{
  Name:BrowserLoginBindingCookie,Path:"/",Secure:true,HttpOnly:true,
  SameSite:http.SameSiteLaxMode,MaxAge:-1,
 })
}

// BrowserLoginBindingFromRequest rejects missing or duplicate cookies.
// The caller must additionally validate state via atomic store.Consume.
func BrowserLoginBindingFromRequest(r *http.Request)(string,error){
 if r==nil{return "",ErrInvalidBrowserLoginBinding}
 cookies:=r.CookiesNamed(BrowserLoginBindingCookie)
 if len(cookies)!=1||len(cookies[0].Value)!=43{return "",ErrInvalidBrowserLoginBinding}
 if _,err:=base64.RawURLEncoding.DecodeString(cookies[0].Value);err!=nil{return "",ErrInvalidBrowserLoginBinding}
 return cookies[0].Value,nil
}
