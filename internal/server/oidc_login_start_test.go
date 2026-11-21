package server

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "net/url"
 "testing"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

type loginStartStore struct{calls int;err error;tx auth.BrowserLoginTransaction;binding string}
func (s *loginStartStore) Create(_ context.Context,tx auth.BrowserLoginTransaction,binding string,ttl time.Duration)error{
 s.calls++;s.tx=tx;s.binding=binding
 return s.err
}
func loginStartRequest(method string)*http.Request{
 r:=httptest.NewRequest(method,"https://domains.example/login/start",nil)
 r.Header.Set("Origin","https://domains.example")
 return r
}
func TestOIDCLoginStartHandler(t *testing.T){
 store:=&loginStartStore{}
 cfg:=OIDCLoginStartConfig{AuthorizationEndpoint:"https://issuer.example/authorize",ClientID:"client",RedirectURI:"https://domains.example/login/callback",TTL:5*time.Minute}
 rec:=httptest.NewRecorder()
 NewOIDCLoginStartHandler(cfg,store).ServeHTTP(rec,loginStartRequest("POST"))
 if rec.Code!=303||store.calls!=1{t.Fatalf("status=%d calls=%d",rec.Code,store.calls)}
 u,err:=url.Parse(rec.Header().Get("Location"));if err!=nil{t.Fatal(err)}
 if u.Query().Get("state")!=store.tx.State||u.Query().Get("nonce")!=store.tx.Nonce||u.Query().Get("code_challenge_method")!="S256"{t.Fatal("incorrect OIDC redirect")}
 cookies:=rec.Result().Cookies()
 if len(cookies)!=1||cookies[0].Value!=store.binding||!cookies[0].HttpOnly||!cookies[0].Secure{t.Fatal("missing secure browser binding cookie")}
}
func TestOIDCLoginStartFailsClosed(t *testing.T){
 cfg:=OIDCLoginStartConfig{AuthorizationEndpoint:"https://issuer.example/authorize",ClientID:"client",RedirectURI:"https://domains.example/callback",TTL:5*time.Minute}
 for _,tc:=range []struct{name string;change func(*http.Request);err error;want int}{
  {"wrong method",func(r *http.Request){r.Method="GET"},nil,405},
  {"missing origin",func(r *http.Request){r.Header.Del("Origin")},nil,403},
  {"foreign origin",func(r *http.Request){r.Header.Set("Origin","https://evil.example")},nil,403},
  {"insecure",func(r *http.Request){r.TLS=nil},nil,403},
  {"store failure",func(r *http.Request){},errors.New("db down"),503},
 }{
  t.Run(tc.name,func(t *testing.T){
   store:=&loginStartStore{err:tc.err}
   req:=loginStartRequest("POST");tc.change(req)
   rec:=httptest.NewRecorder()
   NewOIDCLoginStartHandler(cfg,store).ServeHTTP(rec,req)
   if rec.Code!=tc.want{t.Fatalf("got %d want %d",rec.Code,tc.want)}
   if rec.Header().Get("Location")!=""||len(rec.Result().Cookies())!=0{t.Fatal("redirect or cookie on rejected request")}
  })
 }
}
