package server

import (
 "context"
 "errors"
 "net/http"
 "net/netip"
)

// OIDCDiscoveryValidator allows deterministic startup-gate testing.
type OIDCDiscoveryValidator interface{
 Validate(context.Context,OIDCDeploymentConfig)error
}

var ErrOIDCLoginNotReady=errors.New("OIDC login is not ready")

// OIDCLoginHandlers are constructed only after configuration and discovery
// validation. They are NOT automatically registered on any HTTP router.
type OIDCLoginHandlers struct{
 Start http.Handler
 Callback http.Handler
}

// PrepareOIDCLoginHandlers fails closed before constructing either handler.
// Call it at startup with a trusted discovery client and only register the
// returned handlers after all application startup checks succeed.
func PrepareOIDCLoginHandlers(ctx context.Context,cfg OIDCDeploymentConfig,discovery OIDCDiscoveryValidator,transactions OIDCLoginTransactionWriter,callback OIDCCallbackDependencies)(OIDCLoginHandlers,error){
 if ctx==nil||discovery==nil||transactions==nil||cfg.Validate()!=nil{
  return OIDCLoginHandlers{},ErrOIDCLoginNotReady
 }
 if callback.Transactions==nil||callback.Exchanger==nil||callback.Verifier==nil||
  callback.Users==nil||callback.Sessions==nil||callback.SessionTTL!=cfg.SessionTTL{
  return OIDCLoginHandlers{},ErrOIDCLoginNotReady
 }
 if err:=discovery.Validate(ctx,cfg);err!=nil{
  return OIDCLoginHandlers{},ErrOIDCLoginNotReady
 }
 startCfg,err:=cfg.LoginStartConfig()
 if err!=nil{return OIDCLoginHandlers{},ErrOIDCLoginNotReady}
 callback.HTTPSBoundary=&TrustedHTTPSBoundary{
  PublicHost:cfg.PublicHost,
  ProxyAddresses:append([]netip.Addr(nil),cfg.ProxyAddresses...),
 }
 return OIDCLoginHandlers{
  Start:NewOIDCLoginStartHandler(startCfg,transactions),
  Callback:NewOIDCCallbackHandler(callback),
 },nil
}
