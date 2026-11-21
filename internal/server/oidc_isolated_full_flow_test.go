package server

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "net/url"
 "sync"
 "testing"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

// flowTransactionStore models the one-time, browser-bound contract in memory.
// PostgreSQL transaction semantics are covered separately by integration tests.
type flowTransactionStore struct {
 mu sync.Mutex
 tx auth.BrowserLoginTransaction
 binding string
 consumed bool
}
func (s *flowTransactionStore)Create(_ context.Context,tx auth.BrowserLoginTransaction,binding string,_ time.Duration)error{
 s.mu.Lock();defer s.mu.Unlock()
 s.tx=tx;s.binding=binding;s.consumed=false
 return nil
}
func (s *flowTransactionStore)Consume(_ context.Context,state,binding string)(auth.BrowserLoginTransaction,error){
 s.mu.Lock();defer s.mu.Unlock()
 if s.consumed||s.tx.State==""||state!=s.tx.State||binding!=s.binding{
  return auth.BrowserLoginTransaction{},errors.New("invalid or consumed transaction")
 }
 s.consumed=true
 return s.tx,nil
}
type flowExchanger struct{expectedPKCE string;calls int}
func (e *flowExchanger)ExchangeCode(_ context.Context,code,pkce string)(string,error){
 e.calls++
 if code!="valid-code"||pkce!=e.expectedPKCE{return "",errors.New("invalid code or PKCE")}
 return "id-token",nil
}
type flowVerifier struct{expectedNonce string;calls int}
func (v *flowVerifier)VerifyBrowserIDToken(_ context.Context,token,nonce string)(auth.Identity,error){
 v.calls++
 if token!="id-token"||nonce!=v.expectedNonce{return auth.Identity{},errors.New("invalid nonce")}
 return auth.Identity{Issuer:"https://issuer.example",Subject:"subject",Audience:"client"},nil
}
func TestIsolatedOIDCCompleteFlowAndReplayRejection(t *testing.T){
 cfg:=validOIDCDeploymentConfig()
 transactions:=&flowTransactionStore{}
 exchanger:=&flowExchanger{}
 verifier:=&flowVerifier{}
 sessions:=&callbackSessions{}
 deps:=OIDCCallbackDependencies{
  Transactions:transactions,Exchanger:exchanger,Verifier:verifier,
  Users:portfolioTestResolver{},Sessions:sessions,SessionTTL:cfg.SessionTTL,
 }
 handlers,err:=PrepareOIDCLoginHandlers(context.Background(),cfg,&startupDiscovery{},transactions,deps)
 if err!=nil{t.Fatal(err)}
 router:=isolatedOIDCRouter(handlers)
 start:=httptest.NewRequest(http.MethodPost,"https://domains.example/login/start",nil)
 start.Header.Set("Origin","https://domains.example")
 startRec:=httptest.NewRecorder()
 router.ServeHTTP(startRec,start)
 if startRec.Code!=http.StatusSeeOther{t.Fatalf("start: %d",startRec.Code)}
 redirect,err:=url.Parse(startRec.Header().Get("Location"))
 if err!=nil{t.Fatal(err)}
 if redirect.Query().Get("state")!=transactions.tx.State||redirect.Query().Get("nonce")!=transactions.tx.Nonce{
  t.Fatal("authorization redirect lost transaction binding")
 }
 exchanger.expectedPKCE=transactions.tx.PKCEVerifier
 verifier.expectedNonce=transactions.tx.Nonce
 var binding *http.Cookie
 for _,cookie:=range startRec.Result().Cookies(){
  if cookie.Name==auth.BrowserLoginBindingCookie{binding=cookie}
 }
 if binding==nil{t.Fatal("start did not issue binding cookie")}
 callbackURL:="https://domains.example/login/callback?state="+url.QueryEscape(transactions.tx.State)+"&code=valid-code"
 call:=func(cookie *http.Cookie)*httptest.ResponseRecorder{
  req:=httptest.NewRequest(http.MethodGet,callbackURL,nil)
  if cookie!=nil{req.AddCookie(cookie)}
  rec:=httptest.NewRecorder()
  router.ServeHTTP(rec,req)
  return rec
 }
 missing:=call(nil)
 if missing.Code!=http.StatusBadRequest||transactions.consumed{
  t.Fatal("missing binding consumed transaction")
 }
 success:=call(binding)
 if success.Code!=http.StatusSeeOther||sessions.calls!=1||exchanger.calls!=1||verifier.calls!=1{
  t.Fatalf("callback=%d sessions=%d exchange=%d verify=%d",success.Code,sessions.calls,exchanger.calls,verifier.calls)
 }
 hasSession,clearedBinding:=false,false
 for _,cookie:=range success.Result().Cookies(){
  if cookie.Name=="__Host-ploos_session"&&cookie.Secure&&cookie.HttpOnly{hasSession=true}
  if cookie.Name==auth.BrowserLoginBindingCookie&&cookie.MaxAge<0{clearedBinding=true}
 }
 if !hasSession||!clearedBinding{t.Fatal("successful callback cookie contract failed")}
 replay:=call(binding)
 if replay.Code!=http.StatusUnauthorized||sessions.calls!=1||exchanger.calls!=1{
  t.Fatalf("replay accepted: status=%d sessions=%d exchange=%d",replay.Code,sessions.calls,exchanger.calls)
 }
 for _,cookie:=range replay.Result().Cookies(){
  if cookie.Name=="__Host-ploos_session"{t.Fatal("replay issued session")}
 }
}
