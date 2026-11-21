package main

import (
 "context"
 "errors"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

var errBrowserOIDCDiscovery=errors.New("browser OIDC discovery startup check failed")

// checkBrowserOIDCStartupDiscovery bounds the provider metadata request.
// Production startup passes server.OIDCDiscoveryClient with default TLS checks.
func checkBrowserOIDCStartupDiscovery(ctx context.Context,cfg server.OIDCDeploymentConfig,validator server.OIDCDiscoveryValidator)error{
 if ctx==nil||validator==nil||cfg.Validate()!=nil{return errBrowserOIDCDiscovery}
 probeCtx,cancel:=context.WithTimeout(ctx,10*time.Second)
 defer cancel()
 if err:=validator.Validate(probeCtx,cfg);err!=nil{return errBrowserOIDCDiscovery}
 return nil
}
