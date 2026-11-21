package server

import (
 "net/http"
 "net/http/httptest"
 "testing"
)

// isolatedOIDCRouter is test-only. Production NewHandler never registers
// these routes. The explicit paths are not an activation mechanism.
func isolatedOIDCRouter(handlers OIDCLoginHandlers)http.Handler{
 mux:=http.NewServeMux()
 mux.Handle("POST /login/start",handlers.Start)
 mux.Handle("GET /login/callback",handlers.Callback)
 return mux
}

func TestIsolatedOIDCRoutesMethodAndHostBoundary(t *testing.T){
 cfg:=validOIDCDeploymentConfig()
 store:=&loginStartStore{}
 handlers,err:=PrepareOIDCLoginHandlers(t.Context(),cfg,&startupDiscovery{},store,startupCallbackDeps())
 if err!=nil{t.Fatal(err)}
 router:=isolatedOIDCRouter(handlers)
 for _,tc:=range []struct{name,method,path,origin string;want int}{
  {"login start",http.MethodPost,"/login/start","https://domains.example",http.StatusSeeOther},
  {"start wrong method",http.MethodGet,"/login/start","https://domains.example",http.StatusMethodNotAllowed},
  {"start missing origin",http.MethodPost,"/login/start","",http.StatusForbidden},
  {"callback wrong method",http.MethodPost,"/login/callback","https://domains.example",http.StatusMethodNotAllowed},
  {"callback missing params",http.MethodGet,"/login/callback","",http.StatusBadRequest},
  {"unknown path",http.MethodGet,"/login/unknown","",http.StatusNotFound},
 }{
  t.Run(tc.name,func(t *testing.T){
   req:=httptest.NewRequest(tc.method,"https://domains.example"+tc.path,nil)
   if tc.origin!=""{req.Header.Set("Origin",tc.origin)}
   rec:=httptest.NewRecorder()
   router.ServeHTTP(rec,req)
   if rec.Code!=tc.want{t.Fatalf("status=%d want=%d",rec.Code,tc.want)}
   if tc.want!=http.StatusSeeOther&&(rec.Header().Get("Location")!=""||len(rec.Result().Cookies())!=0){
    t.Fatal("rejected request produced redirect or cookie")
   }
  })
 }
 if store.calls!=1{t.Fatalf("unexpected transaction writes: %d",store.calls)}
}

func TestIsolatedOIDCRoutesNeverAppearInDefaultHandler(t *testing.T){
 handler:=NewHandler()
 for _,tc:=range []struct{method,path string}{
  {http.MethodPost,"/login/start"},
  {http.MethodGet,"/login/callback?state=example&code=example"},
 }{
  req:=httptest.NewRequest(tc.method,"https://domains.example"+tc.path,nil)
  req.Header.Set("Origin","https://domains.example")
  rec:=httptest.NewRecorder()
  handler.ServeHTTP(rec,req)
  if rec.Code!=http.StatusNotFound{t.Fatalf("default handler exposed %s: %d",tc.path,rec.Code)}
 }
}
