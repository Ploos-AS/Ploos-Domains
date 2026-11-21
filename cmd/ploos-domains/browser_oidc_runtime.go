package main

import (
 "context"
 "database/sql"
 "errors"
 "strings"
 "sync"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

var errBrowserOIDCRuntime=errors.New("browser OIDC runtime initialization failed")

// browserOIDCRuntime owns a restricted PostgreSQL connection for the
// lifetime of its unmounted handlers. Call Close on shutdown.
type browserOIDCRuntime struct{
 Handlers server.OIDCLoginHandlers
 DB *sql.DB
 closeOnce sync.Once
 closeErr error
}
func (r *browserOIDCRuntime)Close()error{
 if r==nil{return nil}
 r.closeOnce.Do(func(){
  if r.DB!=nil{r.closeErr=r.DB.Close()}
 })
 return r.closeErr
}

// initializeBrowserOIDCRuntime validates provider metadata and DB privileges,
// initializes a signed ID-token verifier and constructs unmounted handlers.
// The caller must not mount these handlers until the full startup policy is
// separately reviewed and enabled.
func initializeBrowserOIDCRuntime(ctx context.Context,cfg server.OIDCDeploymentConfig,dsn string)(*browserOIDCRuntime,error){
 if ctx==nil||ctx.Err()!=nil||cfg.Validate()!=nil||strings.TrimSpace(dsn)==""{return nil,errBrowserOIDCRuntime}
 probeCtx,cancel:=context.WithTimeout(ctx,10*time.Second)
 defer cancel()
 discovery:=server.OIDCDiscoveryClient{}
 if err:=discovery.Validate(probeCtx,cfg);err!=nil{return nil,errBrowserOIDCRuntime}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{return nil,errBrowserOIDCRuntime}
 success:=false
 defer func(){if !success{_ =db.Close()}}()
 db.SetMaxOpenConns(5)
 db.SetMaxIdleConns(2)
 db.SetConnMaxLifetime(30*time.Minute)
 if err:=db.PingContext(probeCtx);err!=nil{return nil,errBrowserOIDCRuntime}
 if err:=checkBrowserOIDCDatabase(probeCtx,db);err!=nil{return nil,errBrowserOIDCRuntime}
 verifier,err:=auth.NewOIDCVerifier(probeCtx,cfg.Issuer,cfg.ClientID)
 if err!=nil{return nil,errBrowserOIDCRuntime}
 handlers,err:=prepareBrowserOIDCLoginHandlers(probeCtx,cfg,db,verifier,discovery)
 if err!=nil{return nil,errBrowserOIDCRuntime}
 success=true
 return &browserOIDCRuntime{Handlers:handlers,DB:db},nil
}
