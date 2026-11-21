package main

import (
 "context"
 "database/sql"
 "errors"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

var errBrowserOIDCPrepareHandlers=errors.New("browser OIDC login handlers are not ready")

// prepareBrowserOIDCLoginHandlers composes already-provisioned, trusted
// dependencies. It neither opens a database nor mounts HTTP routes.
// The caller owns db and must keep it open for the lifetime of the handlers.
func prepareBrowserOIDCLoginHandlers(ctx context.Context,cfg server.OIDCDeploymentConfig,db *sql.DB,verifier server.OIDCBrowserTokenVerifier,discovery server.OIDCDiscoveryValidator)(server.OIDCLoginHandlers,error){
 if ctx==nil||db==nil||verifier==nil||discovery==nil||cfg.Validate()!=nil{
  return server.OIDCLoginHandlers{},errBrowserOIDCPrepareHandlers
 }
 transactions,callback,err:=prepareBrowserOIDCPostgresLoginStores(db,verifier)
 if err!=nil{return server.OIDCLoginHandlers{},errBrowserOIDCPrepareHandlers}
 callback,err=prepareBrowserOIDCCallbackDependencies(cfg,callback)
 if err!=nil{return server.OIDCLoginHandlers{},errBrowserOIDCPrepareHandlers}
 handlers,err:=server.PrepareOIDCLoginHandlers(ctx,cfg,discovery,transactions,callback)
 if err!=nil{return server.OIDCLoginHandlers{},errBrowserOIDCPrepareHandlers}
 return handlers,nil
}
