package server

import (
 "crypto/tls"
 "net/http/httptest"
 "net/netip"
 "testing"
)

func TestTrustedHTTPSBoundary(t *testing.T){
 b:=TrustedHTTPSBoundary{ProxyAddresses:[]netip.Addr{netip.MustParseAddr("127.0.0.1")},PublicHost:"domains.example"}
 cases:=[]struct{name,remote,proto,host string;tls bool;want bool}{
  {"direct TLS","203.0.113.4:4000","","",true,true},
  {"trusted proxy","127.0.0.1:4000","https","domains.example",false,true},
  {"untrusted spoof","203.0.113.4:4000","https","domains.example",false,false},
  {"insecure proxy","127.0.0.1:4000","http","domains.example",false,false},
  {"wrong forwarded host","127.0.0.1:4000","https","attacker.example",false,false},
  {"chained proto","127.0.0.1:4000","https,http","domains.example",false,false},
  {"missing host","127.0.0.1:4000","https","",false,false},
 }
 for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){
  r:=httptest.NewRequest("GET","http://domains.example/login",nil)
  r.RemoteAddr=tc.remote
  if tc.tls{r.TLS=&tls.ConnectionState{}}
  if tc.proto!=""{r.Header.Set("X-Forwarded-Proto",tc.proto)}
  if tc.host!=""{r.Header.Set("X-Forwarded-Host",tc.host)}
  if got:=b.IsHTTPS(r);got!=tc.want{t.Fatalf("got %v want %v",got,tc.want)}
 })}
}
func TestTrustedHTTPSBoundaryRejectsDuplicateHeaders(t *testing.T){
 b:=TrustedHTTPSBoundary{ProxyAddresses:[]netip.Addr{netip.MustParseAddr("127.0.0.1")},PublicHost:"domains.example"}
 r:=httptest.NewRequest("GET","http://domains.example/login",nil)
 r.RemoteAddr="127.0.0.1:1000"
 r.Header.Add("X-Forwarded-Proto","https")
 r.Header.Add("X-Forwarded-Proto","https")
 r.Header.Set("X-Forwarded-Host","domains.example")
 if b.IsHTTPS(r){t.Fatal("accepted duplicate forwarded proto")}
}
