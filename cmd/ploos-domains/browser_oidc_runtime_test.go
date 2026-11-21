package main

import (
 "context"
 "errors"
 "testing"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

func TestBrowserOIDCRuntimeFailsClosedOnInvalidInputs(t *testing.T){
 cases:=[]struct{name string;ctx context.Context;cfg server.OIDCDeploymentConfig;dsn string}{
  {"nil context",nil,server.OIDCDeploymentConfig{},""},
  {"empty config",context.Background(),server.OIDCDeploymentConfig{},"postgres://unused"},
  {"empty DSN",context.Background(),server.OIDCDeploymentConfig{},""},
 }
 for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){
  runtime,err:=initializeBrowserOIDCRuntime(tc.ctx,tc.cfg,tc.dsn)
  if runtime!=nil||!errors.Is(err,errBrowserOIDCRuntime){t.Fatalf("unsafe runtime initialized: %v",err)}
 })}
}

func TestBrowserOIDCRuntimeCloseNil(t *testing.T){
 var runtime *browserOIDCRuntime
 if err:=runtime.Close();err!=nil{t.Fatal(err)}
}
