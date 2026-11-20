package domeneshop

import (
 "context"
 "errors"
 "io"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"
)

type responseDoer struct{ status int; body string }
func (d responseDoer) Do(_ *http.Request)(*http.Response,error) {
 return &http.Response{StatusCode:d.status,Body:io.NopCloser(strings.NewReader(d.body))},nil
}
func TestReadResponseGuards(t *testing.T){
 cfg:=APIConfig{BaseURL:"https://api.example.test",Timeout:time.Second}
 for _,tc:=range []struct{status int;body string;limit int64;want error}{
  {200,"ok",2,nil},
  {200,"too long",3,ErrResponseTooLarge},
  {302,"redirect",100,ErrUnexpectedStatus},
  {401,"secret",100,ErrUnexpectedStatus},
  {500,"internal",100,ErrUnexpectedStatus},
  {200,"ok",0,ErrInvalidResponseLimit},
  {200,"ok",4<<20+1,ErrInvalidResponseLimit},
 } {
  tr,err:=NewTransport(cfg,responseDoer{tc.status,tc.body})
  if err!=nil {t.Fatal(err)}
  body,err:=tr.ReadResponse(context.Background(),"status",tc.limit)
  if tc.want==nil {
   if err!=nil || string(body)!=tc.body {t.Errorf("got %q %v",body,err)}
  } else if !errors.Is(err,tc.want) {t.Errorf("status %d: want %v got %v",tc.status,tc.want,err)}
 }
}
func TestRestrictedClientDoesNotFollowRedirects(t *testing.T){
 target:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.WriteHeader(http.StatusOK)}))
 defer target.Close()
 redirected:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){http.Redirect(w,r,target.URL,http.StatusFound)}))
 defer redirected.Close()
 client,err:=NewRestrictedHTTPClient(APIConfig{BaseURL:"https://api.example.test",Timeout:time.Second})
 if err!=nil {t.Fatal(err)}
 resp,err:=client.Get(redirected.URL)
 if err!=nil {t.Fatal(err)}
 defer resp.Body.Close()
 if resp.StatusCode!=http.StatusFound {t.Fatalf("redirect was followed: %d",resp.StatusCode)}
}
