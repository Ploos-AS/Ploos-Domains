package auth

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "testing"
)

type fakeVerifier struct {identity Identity; err error; calls int}
func (f *fakeVerifier) VerifyToken(_ context.Context,_ string)(Identity,error){f.calls++;return f.identity,f.err}

func TestAuthMiddleware(t *testing.T){
 valid:=Identity{Issuer:"https://id.example.org",Subject:"user-1",Audience:"ploos-domains"}
 cases:=[]struct{name,header string; verifier *fakeVerifier; status int;calls int}{
  {"valid","Bearer token",&fakeVerifier{identity:valid},200,1},
  {"missing","",&fakeVerifier{identity:valid},401,0},
  {"basic","Basic token",&fakeVerifier{identity:valid},401,0},
  {"empty bearer","Bearer",&fakeVerifier{identity:valid},401,0},
  {"extra fields","Bearer one two",&fakeVerifier{identity:valid},401,0},
  {"invalid token","Bearer token",&fakeVerifier{err:errors.New("invalid signature")},401,1},
  {"missing issuer","Bearer token",&fakeVerifier{identity:Identity{Subject:"user",Audience:"app"}},401,1},
  {"missing audience","Bearer token",&fakeVerifier{identity:Identity{Issuer:"https://id.example.org",Subject:"user"}},401,1},
  {"missing verifier","Bearer token",nil,503,0},
 }
 for _,tc:=range cases {
  t.Run(tc.name,func(t *testing.T){
   called:=false
   next:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
    called=true
    id,ok:=FromContext(r.Context())
    if !ok||id!=valid{t.Fatal("missing verified identity")}
    w.WriteHeader(200)
   })
   var verifier Verifier
   if tc.verifier!=nil { verifier=tc.verifier }
   h:=Middleware(verifier,next)
   req:=httptest.NewRequest("GET","/",nil)
   if tc.header!=""{req.Header.Set("Authorization",tc.header)}
   rec:=httptest.NewRecorder()
   h.ServeHTTP(rec,req)
   if rec.Code!=tc.status{t.Fatalf("status %d want %d",rec.Code,tc.status)}
   if called!=(tc.status==200){t.Fatalf("handler called=%v",called)}
   if tc.verifier!=nil && tc.verifier.calls!=tc.calls{t.Fatalf("calls=%d want=%d",tc.verifier.calls,tc.calls)}
  })
 }
}

func TestMultipleAuthorizationHeadersDenied(t *testing.T){
 f:=&fakeVerifier{identity:Identity{Issuer:"i",Subject:"s",Audience:"a"}}
 h:=Middleware(f,http.HandlerFunc(func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(200)}))
 req:=httptest.NewRequest("GET","/",nil)
 req.Header.Add("Authorization","Bearer first")
 req.Header.Add("Authorization","Bearer second")
 rec:=httptest.NewRecorder()
 h.ServeHTTP(rec,req)
 if rec.Code!=401||f.calls!=0{t.Fatalf("status=%d verifier calls=%d",rec.Code,f.calls)}
}
