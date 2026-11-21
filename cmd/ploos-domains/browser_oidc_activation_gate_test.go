package main

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "testing"
)

// Regression guard: neither browser OIDC route may be exposed by the default
// application, even if a caller sends plausible OIDC query parameters.
func TestBrowserOIDCDisabledDoesNotExposeLoginRoutes(t *testing.T){
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED","false")
 t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
 handler,closeFn,err:=buildHandler(context.Background())
 if err!=nil{t.Fatal(err)}
 defer closeFn()
 for _,tc:=range []struct{method,path string}{
  {http.MethodPost,"/login/start"},
  {http.MethodGet,"/login/callback?state=test&code=test"},
  {http.MethodPost,"/auth/login"},
  {http.MethodGet,"/auth/callback?state=test&code=test"},
 }{
  t.Run(tc.method+" "+tc.path,func(t *testing.T){
   req:=httptest.NewRequest(tc.method,"https://domains.example"+tc.path,nil)
   req.Header.Set("Origin","https://domains.example")
   rec:=httptest.NewRecorder()
   handler.ServeHTTP(rec,req)
   if rec.Code!=http.StatusNotFound{
    t.Fatalf("unexpected browser login route: %d",rec.Code)
   }
   if rec.Header().Get("Location")!=""||len(rec.Result().Cookies())!=0{
    t.Fatal("disabled login issued redirect or cookie")
   }
  })
 }
}

func TestBrowserOIDCOptInNeverReturnsLiveHandlerWithMissingSecret(t *testing.T){
 setValidBrowserOIDCEnv(t)
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED","true")
 t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
 t.Setenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE","")
 handler,closeFn,err:=buildHandler(context.Background())
 if handler!=nil||closeFn!=nil||!errors.Is(err,errBrowserOIDCSecret){
  t.Fatalf("browser OIDC opt-in unexpectedly started: %v",err)
 }
}

func TestBrowserOIDCInvalidFeatureFlagFailsClosed(t *testing.T){
 for _,flag:=range []string{"TRUE","1","yes","enabled"}{
  t.Run(flag,func(t *testing.T){
   t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED",flag)
   t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
   handler,closeFn,err:=buildHandler(context.Background())
   if handler!=nil||closeFn!=nil||!errors.Is(err,errBrowserLoginNotConfigured){
    t.Fatalf("invalid feature flag accepted: %v",err)
   }
  })
 }
}
