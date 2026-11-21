package main

import (
 "errors"
 "os"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

var errBrowserOIDCExchangeConfig=errors.New("browser OIDC code exchange is not configured")

// loadBrowserOIDCCodeExchanger binds the validated, operator-controlled
// provider settings to the protected confidential client secret. It does not
// make a token request or register a callback route.
func loadBrowserOIDCCodeExchanger(cfg server.OIDCDeploymentConfig)(server.OIDCCodeExchange,error){
 if cfg.Validate()!=nil{return server.OIDCCodeExchange{},errBrowserOIDCExchangeConfig}
 secret,err:=loadBrowserOIDCClientSecret(os.Getenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE"))
 if err!=nil{return server.OIDCCodeExchange{},errBrowserOIDCExchangeConfig}
 return server.OIDCCodeExchange{Config:server.OIDCCodeExchangeConfig{
  TokenEndpoint:cfg.TokenEndpoint,
  ClientID:cfg.ClientID,
  ClientSecret:secret,
  RedirectURI:cfg.RedirectURI,
 }},nil
}
