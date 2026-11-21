//go:build integration

package server

import (
 "context"
 "encoding/json"
 "io"
 "net/http"
 "net/url"
 "os"
 "strings"
 "testing"
 "time"
)

// Manual smoke test of public JWKS availability and basic shape.
// It does not validate any signed token or enable login.
func TestBrowserOIDCLiveJWKSReadiness(t *testing.T){
 if os.Getenv("PLOOS_TEST_LIVE_OIDC_ENABLED")!="true"{
  t.Skip("live OIDC provider test requires explicit opt-in")
 }
 issuer:=os.Getenv("PLOOS_TEST_LIVE_OIDC_ISSUER")
 parsed,err:=url.Parse(issuer)
 if err!=nil||parsed.Scheme!="https"||parsed.Host==""||parsed.User!=nil||parsed.RawQuery!=""||parsed.Fragment!=""||parsed.Path!=""||strings.HasSuffix(issuer,"/"){
  t.Fatal("invalid HTTPS root issuer for live JWKS test")
 }
 ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second)
 defer cancel()
 client:=&http.Client{Timeout:10*time.Second,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}
 getJSON:=func(address string,limit int,into any){
  t.Helper()
  req,err:=http.NewRequestWithContext(ctx,http.MethodGet,address,nil)
  if err!=nil{t.Fatal(err)}
  req.Header.Set("Accept","application/json")
  resp,err:=client.Do(req)
  if err!=nil{t.Fatalf("provider endpoint unavailable: %v",err)}
  defer resp.Body.Close()
  if resp.StatusCode!=http.StatusOK||!strings.HasPrefix(strings.ToLower(resp.Header.Get("Content-Type")),"application/json"){
   t.Fatalf("provider returned invalid JSON response: status %d",resp.StatusCode)
  }
  body,err:=io.ReadAll(io.LimitReader(resp.Body,int64(limit)+1))
  if err!=nil||len(body)>limit{t.Fatal("provider JSON exceeds response limit")}
  if err:=json.Unmarshal(body,into);err!=nil{t.Fatalf("invalid provider JSON: %v",err)}
 }
 var discovery oidcDiscoveryMetadata
 getJSON(issuer+"/.well-known/openid-configuration",65536,&discovery)
 if discovery.Issuer!=issuer{t.Fatal("provider issuer mismatch")}
 jwksURL,err:=url.Parse(discovery.JWKSURI)
 if err!=nil||jwksURL.Scheme!="https"||jwksURL.Host==""||jwksURL.User!=nil||jwksURL.Fragment!=""||jwksURL.RawQuery!=""{
  t.Fatal("provider advertised invalid HTTPS JWKS URL")
 }
 var jwks struct{Keys []struct{
  Kty string `json:"kty"`
  Kid string `json:"kid"`
  Use string `json:"use"`
  Alg string `json:"alg"`
 } `json:"keys"`}
 getJSON(discovery.JWKSURI,131072,&jwks)
 if len(jwks.Keys)==0{t.Fatal("provider JWKS has no keys")}
 hasSigningKey:=false
 for _,key:=range jwks.Keys{
  if key.Kty!=""&&(key.Use==""||key.Use=="sig"){hasSigningKey=true;break}
 }
 if !hasSigningKey{t.Fatal("provider JWKS has no advertised signing key")}
}
