package main

import (
 "errors"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

var errBrowserOIDCCallbackWiring=errors.New("browser OIDC callback dependencies are incomplete")

// prepareBrowserOIDCCallbackDependencies binds a confidential OIDC code
// exchanger to separately provisioned transaction, verifier, identity and
// session dependencies. No handler is mounted by this helper.
func prepareBrowserOIDCCallbackDependencies(cfg server.OIDCDeploymentConfig,dependencies server.OIDCCallbackDependencies)(server.OIDCCallbackDependencies,error){
 if cfg.Validate()!=nil||dependencies.Transactions==nil||dependencies.Verifier==nil||
  dependencies.Users==nil||dependencies.Sessions==nil||dependencies.Exchanger!=nil{
  return server.OIDCCallbackDependencies{},errBrowserOIDCCallbackWiring
 }
 exchanger,err:=loadBrowserOIDCCodeExchanger(cfg)
 if err!=nil{return server.OIDCCallbackDependencies{},errBrowserOIDCCallbackWiring}
 dependencies.Exchanger=exchanger
 dependencies.SessionTTL=cfg.SessionTTL
 return dependencies,nil
}
