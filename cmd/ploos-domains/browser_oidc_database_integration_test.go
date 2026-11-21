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

// CI's PostgreSQL role may own the schema. A privileged test role must fail
// closed; least-privilege role provisioning is tested separately at deploy.
func TestBrowserOIDCDatabaseRejectsSchemaOwner(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL")}
 db,err:=sql.Open("postgres",dsn);if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 var superuser bool
 if err:=db.QueryRowContext(ctx,"SELECT rolsuper FROM pg_roles WHERE rolname=current_user").Scan(&superuser);err!=nil{t.Fatal(err)}
 if !superuser{t.Skip("test connection is not a superuser; role-specific checks belong to deployment tests")}
 if err:=checkBrowserOIDCDatabase(ctx,db);!errors.Is(err,errBrowserOIDCDatabase){t.Fatalf("superuser accepted: %v",err)}
}
