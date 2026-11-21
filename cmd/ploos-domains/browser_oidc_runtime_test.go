package main

import (
 "context"
 "database/sql"
 "sync"
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

func TestBrowserOIDCRuntimeCloseConcurrentAndRepeated(t *testing.T){
 db,err:=sql.Open("postgres","postgres://unused:unused@localhost:5432/unused")
 if err!=nil{t.Fatal(err)}
 runtime:=&browserOIDCRuntime{DB:db}
 var wg sync.WaitGroup
 for i:=0;i<16;i++{
  wg.Add(1)
  go func(){
   defer wg.Done()
   if err:=runtime.Close();err!=nil{t.Errorf("close: %v",err)}
  }()
 }
 wg.Wait()
 if err:=runtime.Close();err!=nil{t.Fatal(err)}
 if err:=db.PingContext(context.Background());err==nil{
  t.Fatal("runtime database still usable after close")
 }
}

func TestBrowserOIDCRuntimeRejectsCanceledContextBeforeNetwork(t *testing.T){
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 ctx,cancel:=context.WithCancel(context.Background())
 cancel()
 runtime,err:=initializeBrowserOIDCRuntime(ctx,cfg,"postgres://unused")
 if runtime!=nil||!errors.Is(err,errBrowserOIDCRuntime){
  t.Fatalf("canceled initialization accepted: %v",err)
 }
}
