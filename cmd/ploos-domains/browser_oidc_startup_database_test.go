package main

import (
 "context"
 "errors"
 "testing"
)

func TestBrowserOIDCStartupRejectsMissingDatabaseDSN(t *testing.T){
 setValidBrowserOIDCEnv(t)
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED","true")
 t.Setenv("PLOOS_BROWSER_OIDC_DATABASE_URL","")
 t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
 handler,closeDB,err:=buildHandler(context.Background())
 if !errors.Is(err,errBrowserOIDCDatabase)||handler!=nil||closeDB!=nil{
  t.Fatalf("missing DSN accepted: handler=%t closeDB=%t err=%v",handler!=nil,closeDB!=nil,err)
 }
}

func TestBrowserOIDCStartupDatabaseNilAndEmpty(t *testing.T){
 if err:=checkBrowserOIDCStartupDatabase(context.Background(),"");!errors.Is(err,errBrowserOIDCDatabase){t.Fatalf("empty DSN: %v",err)}
 if err:=checkBrowserOIDCStartupDatabase(nil,"postgres://invalid");!errors.Is(err,errBrowserOIDCDatabase){t.Fatalf("nil context: %v",err)}
}
