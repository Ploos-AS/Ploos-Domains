package server

import (
 "errors"
 "net/netip"
 "net/url"
 "strings"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

var ErrInvalidOIDCDeploymentConfig = errors.New("invalid OIDC deployment configuration")

// OIDCDeploymentConfig is validated before any login endpoints are mounted.
// It contains no secret values; secrets are supplied separately at runtime.
type OIDCDeploymentConfig struct {
 Issuer string
 AuthorizationEndpoint string
 TokenEndpoint string
 ClientID string
 RedirectURI string
 PublicHost string
 ProxyAddresses []netip.Addr
 LoginTTL time.Duration
 SessionTTL time.Duration
}

func (c OIDCDeploymentConfig) Validate()error{
 if auth.ValidateOIDCConfig(c.Issuer,c.ClientID)!=nil||
  c.PublicHost==""||strings.ContainsAny(c.PublicHost,"/:@?# 	\r\n")||
  c.LoginTTL<time.Minute||c.LoginTTL>10*time.Minute||
  (BrowserSessionConfig{TTL:c.SessionTTL}).Validate()!=nil{
  return ErrInvalidOIDCDeploymentConfig
 }
 for _,raw:=range []string{c.AuthorizationEndpoint,c.TokenEndpoint,c.RedirectURI}{
  u,err:=url.Parse(raw)
  if err!=nil||u.Scheme!="https"||u.Host==""||u.User!=nil||
   u.Fragment!=""||u.RawQuery!=""||u.Opaque!=""||
   strings.ContainsAny(u.Host," 	\r\n") {
   return ErrInvalidOIDCDeploymentConfig
  }
 }
 callback,_:=url.Parse(c.RedirectURI)
 if callback.Host!=c.PublicHost||callback.Path==""{
  return ErrInvalidOIDCDeploymentConfig
 }
 for _,ip:=range c.ProxyAddresses{
  if !ip.IsValid()||ip.IsUnspecified()||ip.IsMulticast(){
   return ErrInvalidOIDCDeploymentConfig
  }
 }
 return nil
}

// LoginStartConfig constructs only the unmounted handler configuration.
func (c OIDCDeploymentConfig) LoginStartConfig()(OIDCLoginStartConfig,error){
 if err:=c.Validate();err!=nil{return OIDCLoginStartConfig{},err}
 return OIDCLoginStartConfig{
  AuthorizationEndpoint:c.AuthorizationEndpoint,ClientID:c.ClientID,
  RedirectURI:c.RedirectURI,TTL:c.LoginTTL,
  HTTPSBoundary:&TrustedHTTPSBoundary{PublicHost:c.PublicHost,ProxyAddresses:append([]netip.Addr(nil),c.ProxyAddresses...)},
 },nil
}
