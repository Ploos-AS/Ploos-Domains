package auth

import (
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"
)

func TestBrowserLoginBindingCookie(t *testing.T){
 binding,err:=NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
 other,err:=NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
 if binding==other||len(binding)!=43{t.Fatal("weak binding")}
 w:=httptest.NewRecorder()
 if err:=SetBrowserLoginBindingCookie(w,binding,5*time.Minute);err!=nil{t.Fatal(err)}
 cookie:=w.Header().Get("Set-Cookie")
 for _,part:=range []string{BrowserLoginBindingCookie+"=","Secure","HttpOnly","SameSite=Lax","Path=/"}{
  if !strings.Contains(cookie,part){t.Fatalf("missing %s in %s",part,cookie)}
 }
 if strings.Contains(cookie,"Domain="){t.Fatal("cookie must be host-only")}
 req:=httptest.NewRequest("GET","https://domains.example/callback",nil)
 req.AddCookie(&http.Cookie{Name:BrowserLoginBindingCookie,Value:binding})
 got,err:=BrowserLoginBindingFromRequest(req)
 if err!=nil||got!=binding{t.Fatalf("got %q %v",got,err)}
 ClearBrowserLoginBindingCookie(w)
 if !strings.Contains(w.Header().Values("Set-Cookie")[1],"Max-Age=0"){t.Fatal("not cleared")}
}

func TestBrowserLoginBindingRejectsMissingAndDuplicate(t *testing.T){
 req:=httptest.NewRequest("GET","https://domains.example/callback",nil)
 if _,err:=BrowserLoginBindingFromRequest(req);err==nil{t.Fatal("accepted missing cookie")}
 binding,err:=NewBrowserLoginBinding();if err!=nil{t.Fatal(err)}
 req.Header.Set("Cookie",BrowserLoginBindingCookie+"="+binding+"; "+BrowserLoginBindingCookie+"="+binding)
 if _,err:=BrowserLoginBindingFromRequest(req);err==nil{t.Fatal("accepted duplicate cookie")}
 if err:=SetBrowserLoginBindingCookie(httptest.NewRecorder(),binding,11*time.Minute);err==nil{t.Fatal("accepted excessive TTL")}
}
