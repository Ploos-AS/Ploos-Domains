//go:build integration

package server

import (
 "context"
 "crypto/sha256"
 "database/sql"
 "errors"
 "net/http"
 "net/http/httptest"
 "os"
 "testing"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
 _ "github.com/lib/pq"
)

// Uses real PostgreSQL transaction/session stores and HTTP handlers. The
// exchanger, token verifier and identity resolver are intentionally stubbed;
// a live OIDC provider test remains a separate deployment gate.
func TestOIDCCallbackPostgresSessionLifecycle(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL to an isolated migrated PostgreSQL")}
 db,err:=sql.Open("postgres",dsn);if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second);defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 var userID string
 err=db.QueryRowContext(ctx,`INSERT INTO users(external_subject) VALUES($1) RETURNING id::text`,"oidc-callback-integration-"+time.Now().Format("20060102150405.000000000")).Scan(&userID)
 if err!=nil{t.Fatal(err)}
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",userID)})
 tx,err:=auth.NewBrowserLoginTransaction();if err!=nil{t.Fatal(err)}
 binding,err:=auth.NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
 transactions:=auth.BrowserLoginStore{DB:db}
 if err:=transactions.Create(ctx,tx,binding,5*time.Minute);err!=nil{t.Fatal(err)}
 stateHash:=sha256.Sum256([]byte(tx.State))
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM oidc_login_transactions WHERE state_sha256=$1",stateHash[:])})
 sessions:=BrowserSessionStore{DB:db}
 exchanger:=&postgresCallbackExchanger{wantVerifier:tx.PKCEVerifier}
 verifier:=&postgresCallbackVerifier{wantNonce:tx.Nonce}
 handler:=NewOIDCCallbackHandler(OIDCCallbackDependencies{
  Transactions:transactions,Exchanger:exchanger,Verifier:verifier,
  Users:postgresCallbackResolver{userID:userID},Sessions:sessions,SessionTTL:time.Hour,
 })
 request:=func() *http.Request{
  r:=httptest.NewRequest(http.MethodGet,"https://domains.example/callback?state="+tx.State+"&code=expected-code",nil)
  r.AddCookie(&http.Cookie{Name:auth.BrowserLoginBindingCookie,Value:binding})
  return r
 }
 first:=httptest.NewRecorder();handler.ServeHTTP(first,request())
 if first.Code!=303{t.Fatalf("first callback: %d %s",first.Code,first.Body.String())}
 var sessionCookie *http.Cookie
 for _,cookie:=range first.Result().Cookies(){if cookie.Name=="__Host-ploos_session"{sessionCookie=cookie}}
 if sessionCookie==nil{t.Fatal("successful callback did not issue session cookie")}
 got,err:=sessions.Resolve(ctx,sessionCookie.Value)
 if err!=nil||got!=userID{t.Fatalf("session lookup %q: %v",got,err)}
 var storedRaw int
 if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM browser_sessions WHERE token_sha256=$1",[]byte(sessionCookie.Value)).Scan(&storedRaw);err!=nil{t.Fatal(err)}
 if storedRaw!=0{t.Fatal("raw session token stored in database")}
 second:=httptest.NewRecorder();handler.ServeHTTP(second,request())
 if second.Code!=401{t.Fatalf("replayed callback: %d",second.Code)}
 for _,cookie:=range second.Result().Cookies(){if cookie.Name=="__Host-ploos_session"{t.Fatal("replay issued session")}}
 if exchanger.calls!=1||verifier.calls!=1{t.Fatalf("replay reached token flow: %d/%d",exchanger.calls,verifier.calls)}
 if err:=sessions.Revoke(ctx,sessionCookie.Value);err!=nil{t.Fatal(err)}
 if _,err:=sessions.Resolve(ctx,sessionCookie.Value);!errors.Is(err,ErrBrowserSessionInvalid){t.Fatalf("revoked session still valid: %v",err)}
}

type postgresCallbackExchanger struct{wantVerifier string;calls int}
func (e *postgresCallbackExchanger) ExchangeCode(_ context.Context,code,verifier string)(string,error){
 e.calls++
 if code!="expected-code"||verifier!=e.wantVerifier{return "",errors.New("bad code or PKCE")}
 return "test-id-token",nil
}
type postgresCallbackVerifier struct{wantNonce string;calls int}
func (v *postgresCallbackVerifier) VerifyBrowserIDToken(_ context.Context,raw,nonce string)(auth.Identity,error){
 v.calls++
 if raw!="test-id-token"||nonce!=v.wantNonce{return auth.Identity{},errors.New("bad nonce")}
 return auth.Identity{Issuer:"https://issuer.example",Subject:"test-user",Audience:"client"},nil
}
type postgresCallbackResolver struct{userID string}
func (r postgresCallbackResolver) ResolveUserID(_ context.Context,_ auth.Identity)(string,error){return r.userID,nil}
