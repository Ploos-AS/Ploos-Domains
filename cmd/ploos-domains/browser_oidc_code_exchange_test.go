package main

import (
 "errors"
 "os"
 "path/filepath"
 "testing"
)

func TestBrowserOIDCCodeExchangerFromSecretFile(t *testing.T){
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 path:=filepath.Join(t.TempDir(),"client-secret")
 if err:=os.WriteFile(path,[]byte("test-only-confidential-secret"),0600);err!=nil{t.Fatal(err)}
 t.Setenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE",path)
 exchanger,err:=loadBrowserOIDCCodeExchanger(cfg)
 if err!=nil{t.Fatal(err)}
 got:=exchanger.Config
 if got.ClientSecret!="test-only-confidential-secret"||got.ClientID!=cfg.ClientID||got.TokenEndpoint!=cfg.TokenEndpoint||got.RedirectURI!=cfg.RedirectURI{
  t.Fatal("exchange configuration does not match trusted settings")
 }
 if err:=os.Chmod(path,0644);err!=nil{t.Fatal(err)}
 if _,err:=loadBrowserOIDCCodeExchanger(cfg);!errors.Is(err,errBrowserOIDCExchangeConfig){
  t.Fatalf("insecure secret accepted: %v",err)
 }
}
