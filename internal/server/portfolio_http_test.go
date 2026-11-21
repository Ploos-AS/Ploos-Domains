package server

import (
 "context"
 "encoding/json"
 "errors"
 "net/http"
 "net/http/httptest"
 "testing"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
 "github.com/Ploos-AS/Ploos-Domains/internal/domains"
)

type portfolioTestVerifier struct{}
func (portfolioTestVerifier) VerifyToken(_ context.Context,token string)(auth.Identity,error){
 if token!="test-valid"{return auth.Identity{},auth.ErrUnauthenticated}
 return auth.Identity{Issuer:"https://issuer.example",Subject:"subject",Audience:"client"},nil
}
type portfolioTestResolver struct{ err error }
func (r portfolioTestResolver) ResolveUserID(_ context.Context,_ auth.Identity)(string,error){
 if r.err!=nil{return "",r.err}
 return "internal-user",nil
}
type portfolioTestLister struct{ err error; called *bool }
func (l portfolioTestLister) ListForUser(_ context.Context,user,org string)([]domains.PortfolioEntry,error){
 if l.called!=nil{*l.called=true}
 if user!="internal-user"||org!="org-a"{return nil,errors.New("unexpected scope")}
 if l.err!=nil{return nil,l.err}
 return []domains.PortfolioEntry{{OrganizationID:"org-a",Domain:"example.no",Registrar:"external",Status:domains.PortfolioUnknown}},nil
}
func TestPortfolioHTTPAuthenticationAndErrors(t *testing.T){
 cases:=[]struct{name,authorization string;resolverErr,listerErr error;want int;calls bool}{
  {"missing token","",nil,nil,401,false},
  {"invalid token","Bearer bad",nil,nil,401,false},
  {"unlinked identity","Bearer test-valid",auth.ErrIdentityNotLinked,nil,403,false},
  {"resolver failure","Bearer test-valid",errors.New("db down"),nil,503,false},
  {"portfolio failure","Bearer test-valid",nil,errors.New("db down"),503,true},
  {"valid","Bearer test-valid",nil,nil,200,true},
 }
 for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){
  called:=false
  h:=NewHandlerWithPortfolio(portfolioTestVerifier{},portfolioTestResolver{err:tc.resolverErr},portfolioTestLister{err:tc.listerErr,called:&called})
  req:=httptest.NewRequest("GET","/api/v1/organizations/org-a/portfolio",nil)
  if tc.authorization!=""{req.Header.Set("Authorization",tc.authorization)}
  rec:=httptest.NewRecorder();h.ServeHTTP(rec,req)
  if rec.Code!=tc.want{t.Fatalf("status %d want %d: %s",rec.Code,tc.want,rec.Body.String())}
  if called!=tc.calls{t.Fatalf("lister called %v want %v",called,tc.calls)}
  if rec.Header().Get("Cache-Control")!="no-store"{t.Fatal("missing no-store")}
  if tc.want==200{
   var body struct{Domains []domains.PortfolioEntry `json:"domains"`}
   if err:=json.Unmarshal(rec.Body.Bytes(),&body);err!=nil{t.Fatal(err)}
   if len(body.Domains)!=1||body.Domains[0].Domain!="example.no"{t.Fatalf("unexpected payload %+v",body)}
  }
 })}
}
func TestPortfolioHTTPRejectsUnknownQueryAndWrongMethod(t *testing.T){
 called:=false
 h:=NewHandlerWithPortfolio(portfolioTestVerifier{},portfolioTestResolver{},portfolioTestLister{called:&called})
 for _,tc:=range []struct{method,url string;want int}{
  {"GET","/api/v1/organizations/org-a/portfolio?unexpected=1",400},
  {"POST","/api/v1/organizations/org-a/portfolio",405},
 }{
  req:=httptest.NewRequest(tc.method,tc.url,nil)
  req.Header.Set("Authorization","Bearer test-valid")
  rec:=httptest.NewRecorder();h.ServeHTTP(rec,req)
  if rec.Code!=tc.want{t.Fatalf("%s %s: %d want %d",tc.method,tc.url,rec.Code,tc.want)}
 }
 if called{t.Fatal("lister called for invalid request")}
}
func TestPortfolioNotMountedWithoutDependencies(t *testing.T){
 h:=NewHandlerWithPortfolio(nil,portfolioTestResolver{},portfolioTestLister{})
 rec:=httptest.NewRecorder()
 h.ServeHTTP(rec,httptest.NewRequest(http.MethodGet,"/api/v1/organizations/org-a/portfolio",nil))
 if rec.Code!=404{t.Fatalf("got %d want 404",rec.Code)}
}
