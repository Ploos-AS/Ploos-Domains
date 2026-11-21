package server

import (
 "crypto/sha256"
 "net/http/httptest"
 "strings"
 "testing"
 "time"
)

func TestBrowserSessionTokenUniqueAndDigest(t *testing.T){
 a,ha,err:=NewBrowserSessionToken();if err!=nil{t.Fatal(err)}
 b,hb,err:=NewBrowserSessionToken();if err!=nil{t.Fatal(err)}
 if a==b||ha==hb||len(a)<40{t.Fatal("session tokens are not independent")}
 if ha!=sha256.Sum256([]byte(a)){t.Fatal("digest mismatch")}
}
func TestBrowserSessionCookieSecurity(t *testing.T){
 w:=httptest.NewRecorder()
 if err:=SetBrowserSessionCookie(w,"opaque-token",BrowserSessionConfig{TTL:time.Hour});err!=nil{t.Fatal(err)}
 cookie:=w.Header().Get("Set-Cookie")
 for _,part:=range []string{"__Host-ploos_session=","Secure","HttpOnly","SameSite=Lax","Path=/"}{
  if !strings.Contains(cookie,part){t.Fatalf("cookie missing %s: %s",part,cookie)}
 }
 if strings.Contains(cookie,"Domain="){t.Fatal("host-only cookie must not set Domain")}
 ClearBrowserSessionCookie(w)
 if !strings.Contains(w.Header().Values("Set-Cookie")[1],"Max-Age=0"){t.Fatal("cookie not cleared")}
}
func TestBrowserSessionRejectsBadTTL(t *testing.T){
 for _,ttl:=range []time.Duration{0,time.Second,25*time.Hour}{
  if (BrowserSessionConfig{TTL:ttl}).Validate()==nil{t.Fatalf("accepted ttl %s",ttl)}
 }
}
