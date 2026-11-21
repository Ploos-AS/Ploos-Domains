package main

import (
 "errors"
 "testing"
)

func setValidBrowserOIDCEnv(t *testing.T){
 t.Helper()
 values:=map[string]string{
  "ISSUER":"https://issuer.example",
  "AUTHORIZATION_ENDPOINT":"https://issuer.example/authorize",
  "TOKEN_ENDPOINT":"https://issuer.example/token",
  "CLIENT_ID":"domains-client",
  "REDIRECT_URI":"https://domains.example/login/callback",
  "PUBLIC_HOST":"domains.example",
  "TRUSTED_PROXY_IPS":"127.0.0.1,::1",
  "LOGIN_TTL":"5m",
  "SESSION_TTL":"1h",
 }
 for key,value:=range values{t.Setenv("PLOOS_BROWSER_OIDC_"+key,value)}
}
func TestLoadBrowserOIDCConfig(t *testing.T){
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 if cfg.PublicHost!="domains.example"||len(cfg.ProxyAddresses)!=2{t.Fatalf("unexpected config: %+v",cfg)}
}
func TestLoadBrowserOIDCConfigRejectsInvalid(t *testing.T){
 cases:=[]struct{key,value string}{
  {"ISSUER","http://issuer.example"},
  {"TOKEN_ENDPOINT","https://evil.example/token?x=1"},
  {"REDIRECT_URI","https://evil.example/callback"},
  {"TRUSTED_PROXY_IPS","0.0.0.0"},
  {"TRUSTED_PROXY_IPS","127.0.0.1, "},
  {"LOGIN_TTL","15m"},
  {"SESSION_TTL","garbage"},
  {"CLIENT_ID",""},
 }
 for _,tc:=range cases{t.Run(tc.key+"="+tc.value,func(t *testing.T){
  setValidBrowserOIDCEnv(t)
  t.Setenv("PLOOS_BROWSER_OIDC_"+tc.key,tc.value)
  if _,err:=loadBrowserOIDCConfig();!errors.Is(err,errBrowserOIDCConfig){t.Fatalf("accepted invalid config: %v",err)}
 })}
}
