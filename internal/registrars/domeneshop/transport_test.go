package domeneshop

import (
 "context"
 "errors"
 "io"
 "net/http"
 "strings"
 "testing"
 "time"
)

type fakeDoer struct{calls int;check func(*http.Request)}
func (f *fakeDoer) Do(req *http.Request)(*http.Response,error){
 f.calls++
 if f.check!=nil {f.check(req)}
 return &http.Response{StatusCode:200,Body:io.NopCloser(strings.NewReader("ok"))},nil
}
func TestTransportReadOnlyAndBounded(t *testing.T){
 f:=&fakeDoer{check:func(req *http.Request){
  if req.Method!=http.MethodGet || req.URL.String()!="https://api.example.test/availability" {t.Errorf("unexpected request: %s %s",req.Method,req.URL)}
  if _,ok:=req.Context().Deadline();!ok {t.Error("missing deadline")}
 }}
 tr,err:=NewTransport(APIConfig{BaseURL:"https://api.example.test",Timeout:time.Second},f)
 if err!=nil {t.Fatal(err)}
 for _,segment:=range []string{"","../x","/absolute","x?secret=1","x#frag","https://evil.test","x%2Fy","."}{
  if _,err:=tr.Get(context.Background(),segment);!errors.Is(err,ErrInvalidRequest){t.Errorf("accepted %q: %v",segment,err)}
 }
 if f.calls!=0 {t.Fatal("invalid requests reached transport")}
 resp,err:=tr.Get(context.Background(),"availability")
 if err!=nil {t.Fatal(err)}
 defer resp.Body.Close()
 if f.calls!=1 {t.Fatal("expected one call")}
}
func TestTransportRequiresExplicitClient(t *testing.T){
 cfg:=APIConfig{BaseURL:"https://api.example.test",Timeout:time.Second}
 if _,err:=NewTransport(cfg,nil);!errors.Is(err,ErrTransportUnavailable){t.Fatalf("got %v",err)}
 var tr *Transport
 if _,err:=tr.Get(context.Background(),"anything");!errors.Is(err,ErrTransportUnavailable){t.Fatalf("got %v",err)}
}
