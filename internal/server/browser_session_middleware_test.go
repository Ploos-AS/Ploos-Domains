package server

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "testing"
)

type fakeBrowserSessionResolver struct{ user string; err error; calls *int }
func (f fakeBrowserSessionResolver) Resolve(_ context.Context,token string)(string,error){
 if f.calls!=nil{*f.calls++}
 if token!="valid"{return "",ErrBrowserSessionInvalid}
 return f.user,f.err
}

func TestBrowserSessionMiddleware(t *testing.T){
 cases:=[]struct{name,cookie,user string;err error;want int;calls int}{
  {"missing","", "u",nil,401,0},
  {"invalid","bad","u",nil,401,1},
  {"valid","valid","u",nil,200,1},
  {"revoked","valid","",ErrBrowserSessionInvalid,401,1},
  {"db failure","valid","",errors.New("db down"),503,1},
  {"empty user","valid","",nil,401,1},
 }
 for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){
  calls:=0
  next:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
   id,ok:=BrowserSessionUserID(r.Context())
   if !ok||id!="u"{t.Fatalf("unexpected session user %q",id)}
   w.WriteHeader(200)
  })
  h:=RequireBrowserSession(fakeBrowserSessionResolver{user:tc.user,err:tc.err,calls:&calls},next)
  req:=httptest.NewRequest("GET","https://domains.example/portal",nil)
  if tc.cookie!=""{req.AddCookie(&http.Cookie{Name:browserSessionCookie,Value:tc.cookie})}
  rec:=httptest.NewRecorder();h.ServeHTTP(rec,req)
  if rec.Code!=tc.want||calls!=tc.calls{t.Fatalf("status %d calls %d, want %d %d",rec.Code,calls,tc.want,tc.calls)}
  if rec.Header().Get("Cache-Control")!="no-store"{t.Fatal("missing no-store")}
 })}
}
func TestBrowserSessionMiddlewareRejectsDuplicateCookies(t *testing.T){
 calls:=0
 h:=RequireBrowserSession(fakeBrowserSessionResolver{user:"u",calls:&calls},http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){t.Fatal("handler reached")}))
 req:=httptest.NewRequest("GET","https://domains.example/portal",nil)
 req.Header.Set("Cookie",browserSessionCookie+"=valid; "+browserSessionCookie+"=valid")
 rec:=httptest.NewRecorder();h.ServeHTTP(rec,req)
 if rec.Code!=401||calls!=0{t.Fatalf("status %d calls %d",rec.Code,calls)}
}
