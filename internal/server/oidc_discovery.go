package server

import (
 "context"
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "net/url"
 "strings"
 "time"
)

var ErrOIDCDiscoveryMismatch=errors.New("OIDC discovery metadata mismatch")

// OIDCDiscoveryClient obtains provider metadata from the fixed configured
// issuer, never from browser input. The HTTP client must validate TLS.
type OIDCDiscoveryClient struct{HTTPClient *http.Client}

type oidcDiscoveryMetadata struct{
 Issuer string `json:"issuer"`
 AuthorizationEndpoint string `json:"authorization_endpoint"`
 TokenEndpoint string `json:"token_endpoint"`
 JWKSURI string `json:"jwks_uri"`
}

func (d OIDCDiscoveryClient) Validate(ctx context.Context,c OIDCDeploymentConfig)error{
 if ctx==nil||c.Validate()!=nil{return ErrOIDCDiscoveryMismatch}
 // This implementation supports issuers without a path, where the standard
 // well-known discovery URL is unambiguous.
 issuerURL,err:=url.Parse(c.Issuer)
 if err!=nil||issuerURL.Path!=""||issuerURL.RawPath!=""{
  return ErrOIDCDiscoveryMismatch
 }
 discoveryURL:=c.Issuer+"/.well-known/openid-configuration"
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,discoveryURL,nil)
 if err!=nil{return ErrOIDCDiscoveryMismatch}
 req.Header.Set("Accept","application/json")
 client:=d.HTTPClient
 if client==nil{client=&http.Client{Timeout:10*time.Second}}
 safe:=*client
 safe.CheckRedirect=func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}
 resp,err:=safe.Do(req)
 if err!=nil{return ErrOIDCDiscoveryMismatch}
 defer resp.Body.Close()
 if resp.StatusCode!=http.StatusOK||!strings.HasPrefix(strings.ToLower(resp.Header.Get("Content-Type")),"application/json"){
  return ErrOIDCDiscoveryMismatch
 }
 raw,err:=io.ReadAll(io.LimitReader(resp.Body,65537))
 if err!=nil||len(raw)>65536{return ErrOIDCDiscoveryMismatch}
 var meta oidcDiscoveryMetadata
 if err:=json.Unmarshal(raw,&meta);err!=nil{return ErrOIDCDiscoveryMismatch}
 if meta.Issuer!=c.Issuer||meta.AuthorizationEndpoint!=c.AuthorizationEndpoint||meta.TokenEndpoint!=c.TokenEndpoint{
  return ErrOIDCDiscoveryMismatch
 }
 jwks,err:=url.Parse(meta.JWKSURI)
 if err!=nil||jwks.Scheme!="https"||jwks.Host==""||jwks.User!=nil||jwks.Fragment!=""{
  return ErrOIDCDiscoveryMismatch
 }
 return nil
}
