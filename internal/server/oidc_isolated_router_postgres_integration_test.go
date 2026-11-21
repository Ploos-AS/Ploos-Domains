//go:build integration

package server

import (
 "context"
 "crypto/sha256"
 "database/sql"
 "net/http"
 "net/http/httptest"
 "net/url"
 "os"
 "testing"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
 _ "github.com/lib/pq"
)

// This covers the real PostgreSQL transaction and session stores through the
// test-only router. Provider code exchange and ID-token validation are fakes.
func TestIsolatedOIDCRouterPostgresLifecycle(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL to a disposable migrated PostgreSQL")}
 db,err:=sql.Open("postgres",dsn);if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second);defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 var userID string
 err=db.QueryRowContext(ctx,"INSERT INTO users(external_subject) VALUES($1) RETURNING id::text","oidc-router-"+time.Now().Format("20060102150405.000000000")).Scan(&userID)
 if err!=nil{t.Fatal(err)}
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",userID)})
 cfg:=validOIDCDeploymentConfig()
 transactions:=auth.BrowserLoginStore{DB:db}
 sessions:=BrowserSessionStore{DB:db}
 exchanger:=&postgresCallbackExchanger{}
 verifier:=&postgresCallbackVerifier{}
 deps:=OIDCCallbackDependencies{
  Transactions:transactions,Exchanger:exchanger,Verifier:verifier,
  Users:postgresCallbackResolver{userID:userID},Sessions:sessions,SessionTTL:cfg.SessionTTL,
 }
 handlers,err:=PrepareOIDCLoginHandlers(ctx,cfg,&startupDiscovery{},transactions,deps)
 if err!=nil{t.Fatal(err)}
 router:=isolatedOIDCRouter(handlers)
 start:=httptest.NewRequest(http.MethodPost,"https://domains.example/login/start",nil)
 start.Header.Set("Origin","https://domains.example")
 started:=httptest.NewRecorder();router.ServeHTTP(started,start)
 if started.Code!=303{t.Fatalf("start: %d",started.Code)}
 target,err:=url.Parse(started.Header().Get("Location"));if err!=nil{t.Fatal(err)}
 state:=target.Query().Get("state")
 nonce:=target.Query().Get("nonce")
 if state==""||nonce==""{t.Fatal("missing state or nonce")}
 digest:=sha256.Sum256([]byte(state))
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM oidc_login_transactions WHERE state_sha256=$1",digest[:])})
 if err:=db.QueryRowContext(ctx,"SELECT pkce_verifier FROM oidc_login_transactions WHERE state_sha256=$1",digest[:]).Scan(&exchanger.wantVerifier);err!=nil{t.Fatal(err)}
 verifier.wantNonce=nonce
 var binding *http.Cookie
 for _,cookie:=range started.Result().Cookies(){if cookie.Name==auth.BrowserLoginBindingCookie{binding=cookie}}
 if binding==nil{t.Fatal("missing browser binding")}
 callbackURL:="https://domains.example/login/callback?state="+url.QueryEscape(state)+"&code=expected-code"
 call:=func(cookie *http.Cookie)*httptest.ResponseRecorder{
  req:=httptest.NewRequest(http.MethodGet,callbackURL,nil)
  if cookie!=nil{req.AddCookie(cookie)}
  rec:=httptest.NewRecorder();router.ServeHTTP(rec,req)
  return rec
 }
 missing:=call(nil)
 if missing.Code!=400{t.Fatalf("missing binding: %d",missing.Code)}
 wrongBinding,err:=auth.NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
 wrong:=call(&http.Cookie{Name:auth.BrowserLoginBindingCookie,Value:wrongBinding})
 if wrong.Code!=401{t.Fatalf("wrong binding: %d",wrong.Code)}
 success:=call(binding)
 if success.Code!=303{t.Fatalf("callback: %d",success.Code)}
 var session *http.Cookie
 for _,cookie:=range success.Result().Cookies(){if cookie.Name=="__Host-ploos_session"{session=cookie}}
 if session==nil||!session.Secure||!session.HttpOnly{t.Fatal("missing secure session")}
 if got,err:=sessions.Resolve(ctx,session.Value);err!=nil||got!=userID{t.Fatalf("session lookup: %q %v",got,err)}
 replay:=call(binding)
 if replay.Code!=401||exchanger.calls!=1||verifier.calls!=1{t.Fatalf("replay: %d, exchange=%d, verify=%d",replay.Code,exchanger.calls,verifier.calls)}
 for _,cookie:=range replay.Result().Cookies(){if cookie.Name=="__Host-ploos_session"{t.Fatal("replay issued session")}}
 if err:=sessions.Revoke(ctx,session.Value);err!=nil{t.Fatal(err)}
 if _,err:=sessions.Resolve(ctx,session.Value);err!=ErrBrowserSessionInvalid{t.Fatalf("revoked session accepted: %v",err)}
}
