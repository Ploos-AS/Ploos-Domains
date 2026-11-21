package main

import (
 "database/sql"
 "errors"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

var errBrowserOIDCPostgresDependencies=errors.New("browser OIDC PostgreSQL callback dependencies unavailable")

// prepareBrowserOIDCPostgresDependencies binds all callback persistence
// operations to the same checked, least-privilege runtime connection.
// Verifier is supplied separately and must validate signed OIDC ID tokens.
func prepareBrowserOIDCPostgresDependencies(db *sql.DB,verifier server.OIDCBrowserTokenVerifier)(server.OIDCCallbackDependencies,error){
 if db==nil||verifier==nil{return server.OIDCCallbackDependencies{},errBrowserOIDCPostgresDependencies}
 return server.OIDCCallbackDependencies{
  Transactions:auth.BrowserLoginStore{DB:db},
  Verifier:verifier,
  Users:auth.IdentityRepository{DB:db},
  Sessions:server.BrowserSessionStore{DB:db},
 },nil
}

// prepareBrowserOIDCPostgresLoginStores provides both transaction writer
// and callback dependencies without mounting any HTTP routes.
func prepareBrowserOIDCPostgresLoginStores(db *sql.DB,verifier server.OIDCBrowserTokenVerifier)(server.OIDCLoginTransactionWriter,server.OIDCCallbackDependencies,error){
 callback,err:=prepareBrowserOIDCPostgresDependencies(db,verifier)
 if err!=nil{return nil,server.OIDCCallbackDependencies{},err}
 return auth.BrowserLoginStore{DB:db},callback,nil
}
