package server

import (
 "context"
 "net/http"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

// OIDCCallbackDependencies must be supplied by trusted server-side components.
// Exchanger must use the configured HTTPS token endpoint, fixed redirect URI,
// confidential-client credentials where required, and PKCE verifier.
type OIDCCodeExchanger interface {
 ExchangeCode(context.Context,string,string)(string,error)
}
type OIDCTransactionConsumer interface {
 Consume(context.Context,string,string)(auth.BrowserLoginTransaction,error)
}
type OIDCBrowserTokenVerifier interface {
 VerifyBrowserIDToken(context.Context,string,string)(auth.Identity,error)
}
type OIDCSessionWriter interface {
 Create(context.Context,string,string,time.Time,time.Duration)error
}

type OIDCCallbackDependencies struct{
 Transactions OIDCTransactionConsumer
 Exchanger OIDCCodeExchanger
 Verifier OIDCBrowserTokenVerifier
 Users AuditIdentityResolver
 Sessions OIDCSessionWriter
 SessionTTL time.Duration
}

// NewOIDCCallbackHandler is NOT mounted by NewHandler. It requires a trusted
// TLS-terminating deployment boundary and a reviewed OAuth code exchanger.
func NewOIDCCallbackHandler(d OIDCCallbackDependencies)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("X-Content-Type-Options","nosniff")
  if r.Method!=http.MethodGet{
   w.Header().Set("Allow","GET")
   http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return
  }
  if d.Transactions==nil||d.Exchanger==nil||d.Verifier==nil||d.Users==nil||d.Sessions==nil||(BrowserSessionConfig{TTL:d.SessionTTL}).Validate()!=nil{
   http.Error(w,"login unavailable",http.StatusServiceUnavailable);return
  }
  // This primitive expects direct TLS. A reverse-proxy-aware mount must
  // verify the trusted external HTTPS origin before using this handler.
  if r.TLS==nil{http.Error(w,"forbidden",http.StatusForbidden);return}
  q:=r.URL.Query()
  if len(q)!=2||len(q["state"])!=1||len(q["code"])!=1||q.Get("state")==""||q.Get("code")==""{
   http.Error(w,"invalid callback",http.StatusBadRequest);return
  }
  binding,err:=auth.BrowserLoginBindingFromRequest(r)
  if err!=nil{http.Error(w,"invalid callback",http.StatusBadRequest);return}
  // Clear the binding before any response headers are written. Consume is atomic and single-use.
  auth.ClearBrowserLoginBindingCookie(w)
  tx,err:=d.Transactions.Consume(r.Context(),q.Get("state"),binding)
  if err!=nil{http.Error(w,"invalid or expired login",http.StatusUnauthorized);return}
  rawIDToken,err:=d.Exchanger.ExchangeCode(r.Context(),q.Get("code"),tx.PKCEVerifier)
  if err!=nil||rawIDToken==""{http.Error(w,"login unavailable",http.StatusServiceUnavailable);return}
  identity,err:=d.Verifier.VerifyBrowserIDToken(r.Context(),rawIDToken,tx.Nonce)
  if err!=nil||identity.Validate()!=nil{http.Error(w,"unauthorized",http.StatusUnauthorized);return}
  userID,err:=d.Users.ResolveUserID(r.Context(),identity)
  if err!=nil||userID==""{http.Error(w,"forbidden",http.StatusForbidden);return}
  token,_,err:=NewBrowserSessionToken()
  if err!=nil{http.Error(w,"login unavailable",http.StatusServiceUnavailable);return}
  if err:=d.Sessions.Create(r.Context(),userID,token,time.Now().UTC(),d.SessionTTL);err!=nil{
   http.Error(w,"login unavailable",http.StatusServiceUnavailable);return
  }
  if err:=SetBrowserSessionCookie(w,token,BrowserSessionConfig{TTL:d.SessionTTL});err!=nil{
   http.Error(w,"login unavailable",http.StatusServiceUnavailable);return
  }
  http.Redirect(w,r,"/",http.StatusSeeOther)
 })
}
