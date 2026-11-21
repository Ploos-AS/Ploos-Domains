package server

import (
 "errors"
 "net/netip"
 "testing"
 "time"
)

func validOIDCDeploymentConfig()OIDCDeploymentConfig{
 return OIDCDeploymentConfig{
  Issuer:"https://issuer.example",AuthorizationEndpoint:"https://issuer.example/authorize",
  TokenEndpoint:"https://issuer.example/token",ClientID:"domains-client",
  RedirectURI:"https://domains.example/login/callback",PublicHost:"domains.example",
  ProxyAddresses:[]netip.Addr{netip.MustParseAddr("127.0.0.1")},
  LoginTTL:5*time.Minute,SessionTTL:time.Hour,
 }
}
func TestOIDCDeploymentConfigValid(t *testing.T){
 cfg:=validOIDCDeploymentConfig()
 if err:=cfg.Validate();err!=nil{t.Fatal(err)}
 start,err:=cfg.LoginStartConfig()
 if err!=nil||start.HTTPSBoundary==nil||start.HTTPSBoundary.PublicHost!="domains.example"{
  t.Fatalf("invalid derived config: %+v %v",start,err)
 }
 cfg.ProxyAddresses[0]=netip.MustParseAddr("192.0.2.10")
 if start.HTTPSBoundary.ProxyAddresses[0]!=netip.MustParseAddr("127.0.0.1"){t.Fatal("proxy slice was not copied")}
}
func TestOIDCDeploymentConfigRejectsUnsafeValues(t *testing.T){
 tests:=[]struct{name string;mutate func(*OIDCDeploymentConfig)}{
  {"http issuer",func(c *OIDCDeploymentConfig){c.Issuer="http://issuer.example"}},
  {"http authorize",func(c *OIDCDeploymentConfig){c.AuthorizationEndpoint="http://issuer.example/authorize"}},
  {"http token",func(c *OIDCDeploymentConfig){c.TokenEndpoint="http://issuer.example/token"}},
  {"userinfo",func(c *OIDCDeploymentConfig){c.TokenEndpoint="https://user:pass@issuer.example/token"}},
  {"query token",func(c *OIDCDeploymentConfig){c.TokenEndpoint="https://issuer.example/token?secret=x"}},
  {"wrong callback host",func(c *OIDCDeploymentConfig){c.RedirectURI="https://evil.example/callback"}},
  {"http callback",func(c *OIDCDeploymentConfig){c.RedirectURI="http://domains.example/callback"}},
  {"empty client",func(c *OIDCDeploymentConfig){c.ClientID=""}},
  {"empty public host",func(c *OIDCDeploymentConfig){c.PublicHost=""}},
  {"invalid proxy",func(c *OIDCDeploymentConfig){c.ProxyAddresses=[]netip.Addr{{}}}},
  {"long login ttl",func(c *OIDCDeploymentConfig){c.LoginTTL=11*time.Minute}},
  {"missing session ttl",func(c *OIDCDeploymentConfig){c.SessionTTL=0}},
 }
 for _,tc:=range tests{t.Run(tc.name,func(t *testing.T){
  cfg:=validOIDCDeploymentConfig();tc.mutate(&cfg)
  if err:=cfg.Validate();!errors.Is(err,ErrInvalidOIDCDeploymentConfig){t.Fatalf("accepted unsafe config: %v",err)}
  if _,err:=cfg.LoginStartConfig();!errors.Is(err,ErrInvalidOIDCDeploymentConfig){t.Fatalf("constructed handler config: %v",err)}
 })}
}
