package server

import (
 "context"
 "errors"
 "net/http"
)

type browserSessionContextKey struct{}

// BrowserSessionUserID returns the internal user ID resolved from the session
// store. It is never populated from request parameters or untrusted headers.
func BrowserSessionUserID(ctx context.Context)(string,bool){
 id,ok:=ctx.Value(browserSessionContextKey{}).(string)
 return id,ok&&id!=""
}

type BrowserSessionResolver interface{
 Resolve(context.Context,string)(string,error)
}

// RequireBrowserSession validates a host-only session cookie using a trusted
// server-side store. This does not replace OIDC authorization or CSRF checks.
func RequireBrowserSession(store BrowserSessionResolver,next http.Handler)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  if store==nil||next==nil{http.Error(w,"session unavailable",http.StatusServiceUnavailable);return}
  cookies:=r.CookiesNamed(browserSessionCookie)
  if len(cookies)!=1||cookies[0].Value==""{http.Error(w,"unauthorized",http.StatusUnauthorized);return}
  userID,err:=store.Resolve(r.Context(),cookies[0].Value)
  if errors.Is(err,ErrBrowserSessionInvalid)||userID==""&&err==nil{
   http.Error(w,"unauthorized",http.StatusUnauthorized);return
  }
  if err!=nil{http.Error(w,"session unavailable",http.StatusServiceUnavailable);return}
  next.ServeHTTP(w,r.WithContext(context.WithValue(r.Context(),browserSessionContextKey{},userID)))
 })
}
