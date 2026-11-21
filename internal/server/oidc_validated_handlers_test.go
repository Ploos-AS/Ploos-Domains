package server

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "testing"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

type startupDiscovery struct{calls int;err error}
func (s *startupDiscovery) Validate(_ context.Context,_ OIDCDeploymentConfig)error{s.calls++;return s.err}

func startupCallbackDeps()OIDCCallbackDependencies{
 return OIDCCallbackDependencies{
  Transactions:&callbackTransactions{},Exchanger:&callbackExchanger{},
  Verifier:&callbackVerifier{},Users:portfolioTestResolver{},
  Sessions:&callbackSessions{},SessionTTL:time.Hour,
 }
}
func TestOIDCStartupGateRejectsDiscoveryFailure(t *testing.T){
 cfg:=validOIDCDeploymentConfig()
 d:=&startupDiscovery{err:errors.New("metadata mismatch")}
 got,err:=PrepareOIDCLoginHandlers(context.Background(),cfg,d,&loginStartStore{},startupCallbackDeps())
 if !errors.Is(err,ErrOIDCLoginNotReady)||got.Start!=nil||got.Callback!=nil||d.calls!=1{
  t.Fatalf("gate did not fail closed: %+v %v calls=%d",got,err,d.calls)
 }
}
func TestOIDCStartupGateRejectsMissingDependencies(t *testing.T){
 cfg:=validOIDCDeploymentConfig()
 for _,tc:=range []struct{name string;mutate func(*OIDCCallbackDependencies)}{
  {"missing verifier",func(d *OIDCCallbackDependencies){d.Verifier=nil}},
  {"wrong session TTL",func(d *OIDCCallbackDependencies){d.SessionTTL=2*time.Hour}},
  {"missing session writer",func(d *OIDCCallbackDependencies){d.Sessions=nil}},
 }{
  t.Run(tc.name,func(t *testing.T){
   d:=&startupDiscovery{}
   deps:=startupCallbackDeps();tc.mutate(&deps)
   got,err:=PrepareOIDCLoginHandlers(context.Background(),cfg,d,&loginStartStore{},deps)
   if !errors.Is(err,ErrOIDCLoginNotReady)||got.Start!=nil||got.Callback!=nil||d.calls!=0{
    t.Fatalf("accepted invalid dependencies: %+v %v calls=%d",got,err,d.calls)
   }
  })
 }
}
func TestOIDCStartupGateProducesUnmountedHandlers(t *testing.T){
 cfg:=validOIDCDeploymentConfig()
 discovery:=&startupDiscovery{}
 store:=&loginStartStore{}
 handlers,err:=PrepareOIDCLoginHandlers(context.Background(),cfg,discovery,store,startupCallbackDeps())
 if err!=nil||handlers.Start==nil||handlers.Callback==nil||discovery.calls!=1{
  t.Fatalf("startup failed: %+v %v",handlers,err)
 }
 req:=httptest.NewRequest(http.MethodPost,"http://domains.example/login/start",nil)
 req.RemoteAddr="127.0.0.1:1234"
 req.Header.Set("Origin","https://domains.example")
 req.Header.Set("X-Forwarded-Proto","https")
 req.Header.Set("X-Forwarded-Host","domains.example")
 rec:=httptest.NewRecorder()
 handlers.Start.ServeHTTP(rec,req)
 if rec.Code!=303||store.calls!=1{t.Fatalf("start: status=%d calls=%d",rec.Code,store.calls)}
 // A different public host is rejected even when the forwarded headers match.
 req2:=httptest.NewRequest(http.MethodPost,"http://other.example/login/start",nil)
 req2.RemoteAddr="127.0.0.1:1234"
 req2.Header.Set("Origin","https://other.example")
 req2.Header.Set("X-Forwarded-Proto","https")
 req2.Header.Set("X-Forwarded-Host","domains.example")
 rec2:=httptest.NewRecorder();handlers.Start.ServeHTTP(rec2,req2)
 if rec2.Code!=403||store.calls!=1{t.Fatalf("wrong host: status=%d calls=%d",rec2.Code,store.calls)}
 _=auth.BrowserLoginBindingCookie
}
