//go:build integration

package server

import (
 "context"
 "crypto/sha256"
 "database/sql"
 "net/http"
 "net/http/httptest"
 "net/netip"
 "net/url"
 "os"
 "testing"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
 _ "github.com/lib/pq"
)

// This test simulates Caddy forwarding from a trusted loopback peer. It does
// not launch Caddy or an external OIDC provider.
func TestOIDCProxyStartCallbackPostgres(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL to a disposable migrated PostgreSQL")}
 db,err:=sql.Open("postgres",dsn);if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second);defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 var userID string
 err=db.QueryRowContext(ctx,`INSERT INTO users(external_subject) VALUES($1) RETURNING id::text`,"oidc-proxy-test-"+time.Now().Format("20060102150405.000000000")).Scan(&userID)
 if err!=nil{t.Fatal(err)}
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",userID)})
 boundary:=&TrustedHTTPSBoundary{PublicHost:"domains.example",ProxyAddresses:[]netip.Addr{netip.MustParseAddr("127.0.0.1")}}
 transactions:=auth.BrowserLoginStore{DB:db}
 sessions:=BrowserSessionStore{DB:db}
 start:=NewOIDCLoginStartHandler(OIDCLoginStartConfig{
  AuthorizationEndpoint:"https://issuer.example/authorize",ClientID:"client",
  RedirectURI:"https://domains.example/callback",TTL:5*time.Minute,HTTPSBoundary:boundary,
 },transactions)
 proxyRequest:=func(method,target,peer,proto string)*http.Request{
  r:=httptest.NewRequest(method,target,nil)
  r.TLS=nil
  r.RemoteAddr=peer
  r.Header.Set("X-Forwarded-Proto",proto)
  r.Header.Set("X-Forwarded-Host","domains.example")
  return r
 }
 badStart:=proxyRequest(http.MethodPost,"http://domains.example/login/start","203.0.113.10:3000","https")
 badStart.Header.Set("Origin","https://domains.example")
 rejected:=httptest.NewRecorder();start.ServeHTTP(rejected,badStart)
 if rejected.Code!=403{t.Fatalf("untrusted login start: %d",rejected.Code)}
 goodStart:=proxyRequest(http.MethodPost,"http://domains.example/login/start","127.0.0.1:3000","https")
 goodStart.Header.Set("Origin","https://domains.example")
 started:=httptest.NewRecorder();start.ServeHTTP(started,goodStart)
 if started.Code!=303{t.Fatalf("trusted login start: %d %s",started.Code,started.Body.String())}
 location,err:=url.Parse(started.Header().Get("Location"));if err!=nil{t.Fatal(err)}
 state:=location.Query().Get("state")
 nonce:=location.Query().Get("nonce")
 if state==""||nonce==""||location.Query().Get("code_challenge_method")!="S256"{t.Fatal("invalid authorization redirect")}
 var binding *http.Cookie
 for _,cookie:=range started.Result().Cookies(){if cookie.Name==auth.BrowserLoginBindingCookie{binding=cookie}}
 if binding==nil||!binding.Secure||!binding.HttpOnly{t.Fatal("missing secure binding cookie")}
 stateHash:=sha256.Sum256([]byte(state))
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM oidc_login_transactions WHERE state_sha256=$1",stateHash[:])})
 var verifierSecret string
 err=db.QueryRowContext(ctx,"SELECT pkce_verifier FROM oidc_login_transactions WHERE state_sha256=$1",stateHash[:]).Scan(&verifierSecret)
 if err!=nil{t.Fatal(err)}
 exchanger:=&postgresCallbackExchanger{wantVerifier:verifierSecret}
 verifier:=&postgresCallbackVerifier{wantNonce:nonce}
 callback:=NewOIDCCallbackHandler(OIDCCallbackDependencies{
  Transactions:transactions,Exchanger:exchanger,Verifier:verifier,
  Users:postgresCallbackResolver{userID:userID},Sessions:sessions,SessionTTL:time.Hour,HTTPSBoundary:boundary,
 })
 makeCallback:=func(peer,proto string,b *http.Cookie)*http.Request{
  r:=proxyRequest(http.MethodGet,"http://domains.example/callback?state="+state+"&code=expected-code",peer,proto)
  r.AddCookie(b)
  return r
 }
 spoofed:=httptest.NewRecorder()
 callback.ServeHTTP(spoofed,makeCallback("203.0.113.10:3001","https",binding))
 if spoofed.Code!=403{t.Fatalf("untrusted callback: %d",spoofed.Code)}
 wrongBinding,err:=auth.NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
 wrong:=httptest.NewRecorder()
 callback.ServeHTTP(wrong,makeCallback("127.0.0.1:3001","https",&http.Cookie{Name:auth.BrowserLoginBindingCookie,Value:wrongBinding}))
 if wrong.Code!=401{t.Fatalf("wrong browser binding: %d",wrong.Code)}
 accepted:=httptest.NewRecorder()
 callback.ServeHTTP(accepted,makeCallback("127.0.0.1:3001","https",binding))
 if accepted.Code!=303{t.Fatalf("trusted callback: %d %s",accepted.Code,accepted.Body.String())}
 var session *http.Cookie
 for _,cookie:=range accepted.Result().Cookies(){if cookie.Name=="__Host-ploos_session"{session=cookie}}
 if session==nil||!session.Secure||!session.HttpOnly{t.Fatal("missing secure session cookie")}
 if got,err:=sessions.Resolve(ctx,session.Value);err!=nil||got!=userID{t.Fatalf("session resolve: %q %v",got,err)}
 replay:=httptest.NewRecorder()
 callback.ServeHTTP(replay,makeCallback("127.0.0.1:3001","https",binding))
 if replay.Code!=401{t.Fatalf("replay accepted: %d",replay.Code)}
 if exchanger.calls!=1||verifier.calls!=1{t.Fatalf("replay reached token verification: %d/%d",exchanger.calls,verifier.calls)}
}
