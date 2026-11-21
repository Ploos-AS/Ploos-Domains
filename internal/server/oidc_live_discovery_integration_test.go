//go:build integration

package server

import (
 "context"
 "os"
 "testing"
 "time"
)

// Manual opt-in provider smoke test. It never exchanges codes or uses secrets.
// Normal CI skips this test unless the explicit live-test gate is set.
func TestBrowserOIDCLiveProviderDiscovery(t *testing.T){
 if os.Getenv("PLOOS_TEST_LIVE_OIDC_ENABLED")!="true"{
  t.Skip("live OIDC provider test requires explicit opt-in")
 }
 cfg:=OIDCDeploymentConfig{
  Issuer:os.Getenv("PLOOS_TEST_LIVE_OIDC_ISSUER"),
  AuthorizationEndpoint:os.Getenv("PLOOS_TEST_LIVE_OIDC_AUTHORIZATION_ENDPOINT"),
  TokenEndpoint:os.Getenv("PLOOS_TEST_LIVE_OIDC_TOKEN_ENDPOINT"),
  ClientID:os.Getenv("PLOOS_TEST_LIVE_OIDC_CLIENT_ID"),
  RedirectURI:os.Getenv("PLOOS_TEST_LIVE_OIDC_REDIRECT_URI"),
  PublicHost:os.Getenv("PLOOS_TEST_LIVE_OIDC_PUBLIC_HOST"),
  LoginTTL:5*time.Minute,
  SessionTTL:time.Hour,
 }
 if err:=cfg.Validate();err!=nil{t.Fatalf("invalid live OIDC test configuration: %v",err)}
 ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second)
 defer cancel()
 if err:= (OIDCDiscoveryClient{}).Validate(ctx,cfg);err!=nil{
  t.Fatalf("live provider discovery validation failed: %v",err)
 }
}
