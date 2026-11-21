package server

import (
 "net"
 "net/http"
 "net/netip"
 "strings"
)

// TrustedHTTPSBoundary checks HTTPS using direct TLS or headers received
// exclusively from a configured, trusted reverse-proxy IP address.
// Never trust forwarded headers from arbitrary clients.
type TrustedHTTPSBoundary struct {
 ProxyAddresses []netip.Addr
 PublicHost string
}

func (b TrustedHTTPSBoundary) IsHTTPS(r *http.Request)bool{
 if r==nil||b.PublicHost==""||r.Host!=b.PublicHost{return false}
 if r.TLS!=nil{return true}
 host,_,err:=net.SplitHostPort(r.RemoteAddr)
 if err!=nil{return false}
 remote,err:=netip.ParseAddr(host)
 if err!=nil{return false}
 trusted:=false
 for _,address:=range b.ProxyAddresses{
  if remote.Unmap()==address.Unmap(){trusted=true;break}
 }
 if !trusted{return false}
 // Reject ambiguous forwarded values, including comma-separated chains.
 proto:=r.Header.Values("X-Forwarded-Proto")
 forwardedHost:=r.Header.Values("X-Forwarded-Host")
 if len(proto)!=1||len(forwardedHost)!=1{return false}
 return strings.EqualFold(proto[0],"https")&&forwardedHost[0]==b.PublicHost
}
