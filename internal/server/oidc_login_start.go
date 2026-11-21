package server

import (
 "context"
 "net/http"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

// OIDCLoginTransactionWriter stores one-time browser-bound login attempts.
type OIDCLoginTransactionWriter interface {
 Create(context.Context,auth.BrowserLoginTransaction,string,time.Duration)error
}

type OIDCLoginStartConfig struct{
 AuthorizationEndpoint string
 ClientID string
 RedirectURI string
 TTL time.Duration
}

// NewOIDCLoginStartHandler is intentionally NOT mounted by NewHandler.
// The eventual public route must implement origin/CSRF protection for POST
// and be served only over HTTPS before enabling it.
func NewOIDCLoginStartHandler(cfg OIDCLoginStartConfig,store OIDCLoginTransactionWriter)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("X-Content-Type-Options","nosniff")
  if r.Method!=http.MethodPost{
   w.Header().Set("Allow","POST")
   http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
   return
  }
  if store==nil||cfg.TTL<time.Minute||cfg.TTL>10*time.Minute{
   http.Error(w,"login unavailable",http.StatusServiceUnavailable)
   return
  }
  // Reject cross-site POST and ambiguous/missing browser origin.
  if r.TLS==nil||r.Header.Get("Origin")!="https://"+r.Host||r.Host==""{
   http.Error(w,"forbidden",http.StatusForbidden)
   return
  }
  tx,err:=auth.NewBrowserLoginTransaction()
  if err!=nil{http.Error(w,"login unavailable",http.StatusServiceUnavailable);return}
  binding,err:=auth.NewBrowserLoginBinding()
  if err!=nil{http.Error(w,"login unavailable",http.StatusServiceUnavailable);return}
  target,err:=tx.AuthorizationURL(cfg.AuthorizationEndpoint,cfg.ClientID,cfg.RedirectURI)
  if err!=nil{http.Error(w,"login unavailable",http.StatusServiceUnavailable);return}
  if err:=store.Create(r.Context(),tx,binding,cfg.TTL);err!=nil{
   http.Error(w,"login unavailable",http.StatusServiceUnavailable)
   return
  }
  if err:=auth.SetBrowserLoginBindingCookie(w,binding,cfg.TTL);err!=nil{
   http.Error(w,"login unavailable",http.StatusServiceUnavailable)
   return
  }
  http.Redirect(w,r,target,http.StatusSeeOther)
 })
}
