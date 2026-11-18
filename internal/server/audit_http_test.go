package server

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "testing"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

type auditResolverStub struct{ userID string; err error; calls int }
func (s *auditResolverStub) ResolveUserID(_ context.Context,_ auth.Identity)(string,error){s.calls++;return s.userID,s.err}
type auditListerStub struct{ err error; actor,org string; calls int }
func (s *auditListerStub) List(_ context.Context,actor,org string,_ int,_ *auth.AuditCursor)(auth.AuditPage,error){
 s.calls++;s.actor=actor;s.org=org
 return auth.AuditPage{Events:[]auth.AuditRow{}},s.err
}
type auditVerifierStub struct{}
func (auditVerifierStub) VerifyToken(_ context.Context,token string)(auth.Identity,error){
 if token!="good"{return auth.Identity{},auth.ErrUnauthenticated}
 return auth.Identity{Issuer:"https://id.example.org",Subject:"subject",Audience:"ploos-domains"},nil
}
func TestAuditHandlerSecurity(t *testing.T){
 cases:=[]struct{name,authorization,query string; resolverErr,listerErr error;want,calls int}{
  {"missing bearer","","",nil,nil,401,0},
  {"bad bearer","Bearer bad","",nil,nil,401,0},
  {"unlinked identity","Bearer good","",auth.ErrIdentityNotLinked,nil,403,0},
  {"forbidden","Bearer good","",nil,auth.ErrAuditAccessDenied,403,1},
  {"bad limit","Bearer good","?limit=101",nil,nil,400,0},
  {"bad cursor","Bearer good","?cursor=bad!",nil,nil,400,0},
  {"success","Bearer good","?limit=2",nil,nil,200,1},
 }
 for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){
  resolver:=&auditResolverStub{userID:"trusted-user",err:tc.resolverErr}
  lister:=&auditListerStub{err:tc.listerErr}
  mux:=http.NewServeMux()
  mux.Handle(AuditRoute,auth.Middleware(auditVerifierStub{},NewAuditHandler(resolver,lister)))
  req:=httptest.NewRequest(http.MethodGet,"/api/v1/organizations/trusted-org/audit-events"+tc.query,nil)
  if tc.authorization!=""{req.Header.Set("Authorization",tc.authorization)}
  rr:=httptest.NewRecorder();mux.ServeHTTP(rr,req)
  if rr.Code!=tc.want{t.Fatalf("status=%d want=%d body=%s",rr.Code,tc.want,rr.Body.String())}
  if lister.calls!=tc.calls{t.Fatalf("list calls=%d want=%d",lister.calls,tc.calls)}
  if tc.calls==1&&(lister.actor!="trusted-user"||lister.org!="trusted-org"){t.Fatal("untrusted actor or organization")}
  if tc.want==200&&rr.Header().Get("Cache-Control")!="no-store"{t.Fatal("missing no-store")}
 })}
 if !errors.Is(auth.ErrAuditAccessDenied,auth.ErrAuditAccessDenied){t.Fatal("unexpected")}
}
