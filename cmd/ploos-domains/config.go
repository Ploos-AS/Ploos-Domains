package main

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "net/http"
 "os"
 "strings"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
 "github.com/Ploos-AS/Ploos-Domains/internal/server"
 _ "github.com/lib/pq"
)

var errAuditConfig = errors.New("invalid audit API configuration")

// buildHandler is fail-closed: explicit opt-in requires a valid OIDC issuer,
// audience and PostgreSQL DSN. Startup fails rather than silently exposing
// a partially authenticated audit endpoint.
func buildHandler(ctx context.Context)(http.Handler,func()error,error){
 if os.Getenv("PLOOS_AUDIT_API_ENABLED")!="true"{
  return server.NewHandler(),func()error{return nil},nil
 }
 issuer:=os.Getenv("PLOOS_OIDC_ISSUER")
 audience:=os.Getenv("PLOOS_OIDC_AUDIENCE")
 dsn:=os.Getenv("PLOOS_DATABASE_URL")
 if issuer==""||audience==""||strings.TrimSpace(dsn)==""{
  return nil,nil,errAuditConfig
 }
 if err:=auth.ValidateOIDCConfig(issuer,audience);err!=nil{return nil,nil,fmt.Errorf("%w: %v",errAuditConfig,err)}
 verifier,err:=auth.NewOIDCVerifier(ctx,issuer,audience)
 if err!=nil{return nil,nil,fmt.Errorf("initialize audit OIDC verifier: %w",err)}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{return nil,nil,fmt.Errorf("open audit database: %w",err)}
 db.SetMaxOpenConns(10)
 db.SetMaxIdleConns(5)
 db.SetConnMaxLifetime(30*time.Minute)
 pingCtx,cancel:=context.WithTimeout(ctx,5*time.Second)
 defer cancel()
 if err:=db.PingContext(pingCtx);err!=nil{_ = db.Close();return nil,nil,fmt.Errorf("audit database unavailable: %w",err)}
 return server.NewHandlerWithAudit(verifier,auth.IdentityRepository{DB:db},auth.AuditQuery{DB:db}),db.Close,nil
}
