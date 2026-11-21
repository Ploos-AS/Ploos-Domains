package main

import (
 "errors"
 "os"
 "path/filepath"
 "testing"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

func TestBrowserOIDCCallbackWiringRejectsMissingDependencies(t *testing.T){
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 if _,err:=prepareBrowserOIDCCallbackDependencies(cfg,server.OIDCCallbackDependencies{});!errors.Is(err,errBrowserOIDCCallbackWiring){
  t.Fatalf("incomplete callback accepted: %v",err)
 }
}

func TestBrowserOIDCCallbackWiringRejectsMissingSecret(t *testing.T){
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 t.Setenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE","")
 d:=server.OIDCCallbackDependencies{
  Transactions:callbackTransactionStub{},
  Verifier:callbackVerifierStub{},
  Users:callbackUserStub{},
  Sessions:callbackSessionStub{},
 }
 if _,err:=prepareBrowserOIDCCallbackDependencies(cfg,d);!errors.Is(err,errBrowserOIDCCallbackWiring){
  t.Fatalf("missing secret accepted: %v",err)
 }
}

func TestBrowserOIDCCallbackWiringSetsExchanger(t *testing.T){
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 path:=filepath.Join(t.TempDir(),"secret")
 if err:=os.WriteFile(path,[]byte("test-only-secret"),0600);err!=nil{t.Fatal(err)}
 t.Setenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE",path)
 d:=server.OIDCCallbackDependencies{
  Transactions:callbackTransactionStub{},
  Verifier:callbackVerifierStub{},
  Users:callbackUserStub{},
  Sessions:callbackSessionStub{},
 }
 wired,err:=prepareBrowserOIDCCallbackDependencies(cfg,d)
 if err!=nil{t.Fatal(err)}
 exchanger,ok:=wired.Exchanger.(server.OIDCCodeExchange)
 if !ok||exchanger.Config.ClientSecret!="test-only-secret"||wired.SessionTTL!=cfg.SessionTTL{
  t.Fatal("callback exchanger or TTL mismatch")
 }
}
