package server

import (
 "net/http"
 "net/http/httptest"
 "net/netip"
 "testing"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

func TestOIDCLoginStartTrustedProxy(t *testing.T){
 b:=&TrustedHTTPSBoundary{PublicHost:"domains.example",ProxyAddresses:[]netip.Addr{netip.MustParseAddr("127.0.0.1")}}
 cfg:=OIDCLoginStartConfig{AuthorizationEndpoint:"https://issuer.example/authorize",ClientID:"client",RedirectURI:"https://domains.example/callback",TTL:5*time.Minute,HTTPSBoundary:b}
 for _,tc:=range []struct{name,remote,proto,host string;want int}{
  {"trusted","127.0.0.1:1234","https","domains.example",303},
  {"untrusted","203.0.113.10:1234","https","domains.example",403},
  {"wrong proto","127.0.0.1:1234","http","domains.example",403},
  {"wrong host","127.0.0.1:1234","https","evil.example",403},
 }{
  t.Run(tc.name,func(t *testing.T){
   store:=&loginStartStore{}
   r:=loginStartRequest(http.MethodPost)
   r.TLS=nil
   r.RemoteAddr=tc.remote
   r.Header.Set("X-Forwarded-Proto",tc.proto)
   r.Header.Set("X-Forwarded-Host",tc.host)
   rec:=httptest.NewRecorder()
   NewOIDCLoginStartHandler(cfg,store).ServeHTTP(rec,r)
   if rec.Code!=tc.want{t.Fatalf("got %d want %d",rec.Code,tc.want)}
   if tc.want!=303&&store.calls!=0{t.Fatal("stored login transaction on untrusted request")}
  })
 }
}

func TestOIDCCallbackRejectsUntrustedProxyBeforeConsume(t *testing.T){
 b:=&TrustedHTTPSBoundary{PublicHost:"domains.example",ProxyAddresses:[]netip.Addr{netip.MustParseAddr("127.0.0.1")}}
 tx:=&callbackTransactions{}
 deps:=OIDCCallbackDependencies{Transactions:tx,Exchanger:&callbackExchanger{},Verifier:&callbackVerifier{},Users:portfolioTestResolver{},Sessions:&callbackSessions{},SessionTTL:time.Hour,HTTPSBoundary:b}
 binding,err:=auth.NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
 for _,tc:=range []struct{name,remote,proto string;want int}{
  {"trusted","127.0.0.1:4444","https",401},
  {"spoofed","203.0.113.10:4444","https",403},
  {"insecure","127.0.0.1:4444","http",403},
 }{
  t.Run(tc.name,func(t *testing.T){
   tx.calls=0
   r:=httptest.NewRequest(http.MethodGet,"http://domains.example/callback?state=expected-state&code=expected-code",nil)
   r.RemoteAddr=tc.remote
   r.Header.Set("X-Forwarded-Proto",tc.proto)
   r.Header.Set("X-Forwarded-Host","domains.example")
   r.AddCookie(&http.Cookie{Name:auth.BrowserLoginBindingCookie,Value:binding})
   rec:=httptest.NewRecorder()
   NewOIDCCallbackHandler(deps).ServeHTTP(rec,r)
   if rec.Code!=tc.want{t.Fatalf("got %d want %d",rec.Code,tc.want)}
   if tc.want==403&&tx.calls!=0{t.Fatal("consumed transaction from untrusted proxy")}
  })
 }
}
