package main

import (
 "errors"
 "net/netip"
 "os"
 "strings"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

var errBrowserOIDCConfig=errors.New("invalid browser OIDC environment configuration")

// loadBrowserOIDCConfig reads fixed operator-controlled values only.
// It does not enable login or read the client secret.
func loadBrowserOIDCConfig()(server.OIDCDeploymentConfig,error){
 get:=func(key string)string{return os.Getenv("PLOOS_BROWSER_OIDC_"+key)}
 loginTTL,err:=time.ParseDuration(get("LOGIN_TTL"))
 if err!=nil{return server.OIDCDeploymentConfig{},errBrowserOIDCConfig}
 sessionTTL,err:=time.ParseDuration(get("SESSION_TTL"))
 if err!=nil{return server.OIDCDeploymentConfig{},errBrowserOIDCConfig}
 var peers []netip.Addr
 rawPeers:=get("TRUSTED_PROXY_IPS")
 if rawPeers==""{return server.OIDCDeploymentConfig{},errBrowserOIDCConfig}
 for _,raw:=range strings.Split(rawPeers,","){
  if raw==""||strings.TrimSpace(raw)!=raw{return server.OIDCDeploymentConfig{},errBrowserOIDCConfig}
  ip,err:=netip.ParseAddr(raw)
  if err!=nil{return server.OIDCDeploymentConfig{},errBrowserOIDCConfig}
  peers=append(peers,ip)
 }
 cfg:=server.OIDCDeploymentConfig{
  Issuer:get("ISSUER"),AuthorizationEndpoint:get("AUTHORIZATION_ENDPOINT"),
  TokenEndpoint:get("TOKEN_ENDPOINT"),ClientID:get("CLIENT_ID"),
  RedirectURI:get("REDIRECT_URI"),PublicHost:get("PUBLIC_HOST"),
  ProxyAddresses:peers,LoginTTL:loginTTL,SessionTTL:sessionTTL,
 }
 if cfg.Validate()!=nil{return server.OIDCDeploymentConfig{},errBrowserOIDCConfig}
 return cfg,nil
}
