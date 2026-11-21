package main

import (
 "context"
 "errors"
 "testing"
)

func TestBrowserOIDCOptInRejectsInvalidConfiguration(t *testing.T){
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED","true")
 t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_TTL","invalid")
 handler,closeDB,err:=buildHandler(context.Background())
 if !errors.Is(err,errBrowserOIDCConfig)||handler!=nil||closeDB!=nil{
  t.Fatalf("invalid config not rejected: handler=%t closeDB=%t err=%v",handler!=nil,closeDB!=nil,err)
 }
}

func TestBrowserOIDCOptInValidConfigStillFailsClosed(t *testing.T){
 setValidBrowserOIDCEnv(t)
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED","true")
 t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
 handler,closeDB,err:=buildHandler(context.Background())
 if !errors.Is(err,errBrowserLoginNotConfigured)||handler!=nil||closeDB!=nil{
  t.Fatalf("incomplete OIDC wiring enabled: handler=%t closeDB=%t err=%v",handler!=nil,closeDB!=nil,err)
 }
}

func TestBrowserOIDCDisabledIgnoresInvalidBrowserConfiguration(t *testing.T){
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED","false")
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_TTL","invalid")
 t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
 handler,closeDB,err:=buildHandler(context.Background())
 if err!=nil||handler==nil||closeDB==nil{t.Fatalf("disabled login blocked normal startup: %v",err)}
 if err:=closeDB();err!=nil{t.Fatal(err)}
}
