package server

import (
 "context"
 "crypto/tls"
 "errors"
 "net/http"
 "net/http/httptest"
 "testing"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

type callbackRejectResolver struct{calls int}
func (r *callbackRejectResolver) ResolveUserID(_ context.Context,_ auth.Identity)(string,error){
 r.calls++
 return "",errors.New("unknown identity")
}

func TestOIDCCallbackFailureDoesNotIssueSession(t *testing.T){
 for _,tc:=range []struct{name string;transactionErr,exchangeErr,verifyErr,sessionErr error;rejectUser bool;want int;consume,exchange,verify,session int}{
  {"state rejected",errors.New("state mismatch"),nil,nil,nil,false,401,1,0,0,0},
  {"code exchange failed",nil,errors.New("token endpoint down"),nil,nil,false,503,1,1,0,0},
  {"invalid ID token nonce",nil,nil,errors.New("nonce mismatch"),nil,false,401,1,1,1,0},
  {"identity not mapped",nil,nil,nil,nil,true,403,1,1,1,0},
  {"session persistence failed",nil,nil,nil,errors.New("database down"),false,503,1,1,1,1},
 }{
  t.Run(tc.name,func(t *testing.T){
   transactions:=&callbackTransactions{tx:auth.BrowserLoginTransaction{Nonce:"expected-nonce",PKCEVerifier:"pkce"},err:tc.transactionErr}
   exchanger:=&callbackExchanger{err:tc.exchangeErr}
   verifier:=&callbackVerifier{err:tc.verifyErr}
   sessions:=&callbackSessions{err:tc.sessionErr}
   var resolver AuditIdentityResolver=portfolioTestResolver{}
   rejecting:=&callbackRejectResolver{}
   if tc.rejectUser{resolver=rejecting}
   handler:=NewOIDCCallbackHandler(OIDCCallbackDependencies{
    Transactions:transactions,Exchanger:exchanger,Verifier:verifier,Users:resolver,Sessions:sessions,SessionTTL:time.Hour,
   })
   binding,err:=auth.NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
   req:=httptest.NewRequest(http.MethodGet,"https://domains.example/callback?state=expected-state&code=expected-code",nil)
   req.TLS=&tls.ConnectionState{}
   req.AddCookie(&http.Cookie{Name:auth.BrowserLoginBindingCookie,Value:binding})
   rec:=httptest.NewRecorder()
   handler.ServeHTTP(rec,req)
   if rec.Code!=tc.want{t.Fatalf("status %d want %d",rec.Code,tc.want)}
   if transactions.calls!=tc.consume||exchanger.calls!=tc.exchange||verifier.calls!=tc.verify||sessions.calls!=tc.session{
    t.Fatalf("calls consume/exchange/verify/session: %d/%d/%d/%d",transactions.calls,exchanger.calls,verifier.calls,sessions.calls)
   }
   if rec.Header().Get("Location")!=""{t.Fatal("redirected after failed authentication")}
   if rec.Header().Get("Cache-Control")!="no-store"{t.Fatal("missing no-store")}
   cleared:=false
   for _,cookie:=range rec.Result().Cookies(){
    if cookie.Name==browserSessionCookie{t.Fatal("session cookie issued after failed authentication")}
    if cookie.Name==auth.BrowserLoginBindingCookie&&cookie.MaxAge<0{cleared=true}
   }
   if !cleared{t.Fatal("browser binding not cleared after failed callback")}
  })
 }
}

func TestOIDCCallbackRejectsAmbiguousParametersBeforeConsume(t *testing.T){
 for _,query:=range []string{
  "?state=expected-state&state=other&code=expected-code",
  "?state=expected-state&code=expected-code&next=https://evil.example",
  "?state=&code=expected-code",
  "?state=expected-state&code=",
 }{
  t.Run(query,func(t *testing.T){
   transactions:=&callbackTransactions{}
   sessions:=&callbackSessions{}
   handler:=NewOIDCCallbackHandler(OIDCCallbackDependencies{
    Transactions:transactions,Exchanger:&callbackExchanger{},Verifier:&callbackVerifier{},
    Users:portfolioTestResolver{},Sessions:sessions,SessionTTL:time.Hour,
   })
   binding,err:=auth.NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
   req:=httptest.NewRequest(http.MethodGet,"https://domains.example/callback"+query,nil)
   req.TLS=&tls.ConnectionState{}
   req.AddCookie(&http.Cookie{Name:auth.BrowserLoginBindingCookie,Value:binding})
   rec:=httptest.NewRecorder()
   handler.ServeHTTP(rec,req)
   if rec.Code!=400||transactions.calls!=0||sessions.calls!=0{
    t.Fatalf("status=%d consumed=%d sessions=%d",rec.Code,transactions.calls,sessions.calls)
   }
  })
 }
}
