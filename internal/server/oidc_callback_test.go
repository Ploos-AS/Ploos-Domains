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

type callbackTransactions struct{tx auth.BrowserLoginTransaction;err error;calls int}
func (s *callbackTransactions) Consume(_ context.Context,state,binding string)(auth.BrowserLoginTransaction,error){
 s.calls++
 if state!="expected-state"||binding==""{return auth.BrowserLoginTransaction{},errors.New("invalid state")}
 return s.tx,s.err
}
type callbackExchanger struct{calls int;err error}
func (s *callbackExchanger) ExchangeCode(_ context.Context,code,verifier string)(string,error){
 s.calls++
 if code!="expected-code"||verifier!="pkce"{return "",errors.New("invalid code")}
 return "id-token",s.err
}
type callbackVerifier struct{calls int;err error}
func (s *callbackVerifier) VerifyBrowserIDToken(_ context.Context,raw,nonce string)(auth.Identity,error){
 s.calls++
 if raw!="id-token"||nonce!="expected-nonce"{return auth.Identity{},errors.New("invalid nonce")}
 return auth.Identity{Issuer:"https://issuer.example",Subject:"subject",Audience:"client"},s.err
}
type callbackSessions struct{calls int;err error}
func (s *callbackSessions) Create(_ context.Context,user,token string,_ time.Time,_ time.Duration)error{
 s.calls++
 if user!="internal-user"||token==""{return errors.New("invalid session")}
 return s.err
}
func TestOIDCCallbackSuccess(t *testing.T){
 tx:=&callbackTransactions{tx:auth.BrowserLoginTransaction{Nonce:"expected-nonce",PKCEVerifier:"pkce"}}
 ex:=&callbackExchanger{};v:=&callbackVerifier{};sessions:=&callbackSessions{}
 h:=NewOIDCCallbackHandler(OIDCCallbackDependencies{Transactions:tx,Exchanger:ex,Verifier:v,Users:portfolioTestResolver{},Sessions:sessions,SessionTTL:time.Hour})
 binding,err:=auth.NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
 req:=httptest.NewRequest("GET","https://domains.example/callback?state=expected-state&code=expected-code",nil)
 req.TLS=&tls.ConnectionState{}
 req.AddCookie(&http.Cookie{Name:auth.BrowserLoginBindingCookie,Value:binding})
 rec:=httptest.NewRecorder();h.ServeHTTP(rec,req)
 if rec.Code!=303||tx.calls!=1||ex.calls!=1||v.calls!=1||sessions.calls!=1{t.Fatalf("status=%d calls=%d/%d/%d/%d",rec.Code,tx.calls,ex.calls,v.calls,sessions.calls)}
 if len(rec.Result().Cookies())!=2{t.Fatal("expected cleared binding and session cookie")}
}
func TestOIDCCallbackRejectsMissingBinding(t *testing.T){
 tx:=&callbackTransactions{}
 h:=NewOIDCCallbackHandler(OIDCCallbackDependencies{Transactions:tx,Exchanger:&callbackExchanger{},Verifier:&callbackVerifier{},Users:portfolioTestResolver{},Sessions:&callbackSessions{},SessionTTL:time.Hour})
 req:=httptest.NewRequest("GET","https://domains.example/callback?state=expected-state&code=expected-code",nil)
 req.TLS=&tls.ConnectionState{}
 rec:=httptest.NewRecorder();h.ServeHTTP(rec,req)
 if rec.Code!=400||tx.calls!=0{t.Fatalf("status=%d calls=%d",rec.Code,tx.calls)}
}
