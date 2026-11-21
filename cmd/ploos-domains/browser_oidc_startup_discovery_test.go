package main

import (
 "context"
 "errors"
 "testing"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

type discoveryProbe struct{called int;fail bool}
func (p *discoveryProbe) Validate(_ context.Context,_ server.OIDCDeploymentConfig)error{
 p.called++
 if p.fail{return errors.New("mismatched issuer")}
 return nil
}
func TestBrowserOIDCStartupDiscovery(t *testing.T){
 cfg:=validBrowserOIDCStartupConfig(t)
 ok:=&discoveryProbe{}
 if err:=checkBrowserOIDCStartupDiscovery(context.Background(),cfg,ok);err!=nil||ok.called!=1{
  t.Fatalf("valid discovery rejected: %v calls=%d",err,ok.called)
 }
 bad:=&discoveryProbe{fail:true}
 if err:=checkBrowserOIDCStartupDiscovery(context.Background(),cfg,bad);!errors.Is(err,errBrowserOIDCDiscovery)||bad.called!=1{
  t.Fatalf("mismatched discovery accepted: %v calls=%d",err,bad.called)
 }
 if err:=checkBrowserOIDCStartupDiscovery(context.Background(),cfg,nil);!errors.Is(err,errBrowserOIDCDiscovery){
  t.Fatalf("nil validator accepted: %v",err)
 }
 if err:=checkBrowserOIDCStartupDiscovery(nil,cfg,ok);!errors.Is(err,errBrowserOIDCDiscovery){
  t.Fatalf("nil context accepted: %v",err)
 }
}
func validBrowserOIDCStartupConfig(t *testing.T)server.OIDCDeploymentConfig{
 t.Helper()
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 return cfg
}
