package server

import (
 "crypto/rand"
 "crypto/sha256"
 "encoding/base64"
 "errors"
 "net/http"
 "time"
)

var ErrInvalidBrowserSession = errors.New("invalid browser session configuration")

const browserSessionCookie = "__Host-ploos_session"

// BrowserSessionConfig defines cookie and server-side session primitives.
// A future OIDC callback must validate state, nonce, PKCE and token claims
// before issuing a session. This file does not implement login.
type BrowserSessionConfig struct {
 TTL time.Duration
}

func (c BrowserSessionConfig) Validate() error {
 if c.TTL<time.Minute || c.TTL>24*time.Hour {return ErrInvalidBrowserSession}
 return nil
}

// NewBrowserSessionToken creates an opaque random token. Persist only its
// SHA-256 digest in a server-side store, never the raw cookie value.
func NewBrowserSessionToken()(string,[32]byte,error){
 var raw [32]byte
 if _,err:=rand.Read(raw[:]);err!=nil{return "",[32]byte{},err}
 token:=base64.RawURLEncoding.EncodeToString(raw[:])
 return token,sha256.Sum256([]byte(token)),nil
}

// SetBrowserSessionCookie requires an HTTPS response at the deployment
// boundary. The __Host- prefix forbids Domain and requires Path=/ and Secure.
func SetBrowserSessionCookie(w http.ResponseWriter,token string,c BrowserSessionConfig)error{
 if w==nil||token==""||c.Validate()!=nil{return ErrInvalidBrowserSession}
 http.SetCookie(w,&http.Cookie{
  Name:browserSessionCookie,Value:token,Path:"/",HttpOnly:true,
  Secure:true,SameSite:http.SameSiteLaxMode,
  MaxAge:int(c.TTL.Seconds()),Expires:time.Now().Add(c.TTL),
 })
 return nil
}

func ClearBrowserSessionCookie(w http.ResponseWriter){
 if w==nil{return}
 http.SetCookie(w,&http.Cookie{
  Name:browserSessionCookie,Path:"/",HttpOnly:true,Secure:true,
  SameSite:http.SameSiteLaxMode,MaxAge:-1,
 })
}
