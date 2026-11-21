package server

import "net/http"

// oidcRequestIsHTTPS keeps direct TLS support for existing unmounted
// handler tests. Production proxy deployments must configure a boundary.
func oidcRequestIsHTTPS(r *http.Request,boundary *TrustedHTTPSBoundary)bool{
 if r==nil{return false}
 if boundary!=nil{return boundary.IsHTTPS(r)}
 return r.TLS!=nil
}
