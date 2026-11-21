//go:build integration

package main

import (
 "context"
 "database/sql"
 "errors"
 "os"
 "testing"
 "time"

 _ "github.com/lib/pq"
)

func TestBrowserOIDCRestrictedDatabaseRole(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_BROWSER_OIDC_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_BROWSER_OIDC_DATABASE_URL for a disposable least-privilege PostgreSQL role")}
 db,err:=sql.Open("postgres",dsn);if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 if err:=checkBrowserOIDCDatabase(ctx,db);err!=nil{t.Fatalf("valid restricted role rejected: %v",err)}
 // The same session, with one required privilege revoked, must fail closed.
 // This is a disposable CI database. The test must not modify production roles.
 adminDSN:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if adminDSN==""{t.Skip("set PLOOS_TEST_DATABASE_URL to test privilege revocation")}
 admin,err:=sql.Open("postgres",adminDSN);if err!=nil{t.Fatal(err)}
 defer admin.Close()
 if _,err:=admin.ExecContext(ctx,"REVOKE INSERT ON public.browser_sessions FROM ploos_browser_oidc_runtime_ci");err!=nil{t.Fatal(err)}
 defer func(){_,_=admin.ExecContext(context.Background(),"GRANT INSERT ON public.browser_sessions TO ploos_browser_oidc_runtime_ci")}()
 // Clear pooled sessions to avoid cached privilege state across connections.
 db.SetMaxIdleConns(0)
 if err:=checkBrowserOIDCDatabase(ctx,db);!errors.Is(err,errBrowserOIDCDatabase){
  t.Fatalf("missing INSERT privilege accepted: %v",err)
 }
}
