package main

import (
 "context"
 "database/sql"
 "errors"
 "os"
 "path/filepath"
 "testing"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

type startupDiscoveryStub struct{err error}
func (s startupDiscoveryStub) Validate(context.Context,server.OIDCDeploymentConfig)error{return s.err}

func TestBrowserOIDCPrepareHandlersFailsClosed(t *testing.T){
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 db:=&sql.DB{}
 cases:=[]struct{
  name string
  db *sql.DB
  verifier server.OIDCBrowserTokenVerifier
  discovery server.OIDCDiscoveryValidator
 }{
  {"missing database",nil,callbackVerifierStub{},startupDiscoveryStub{}},
  {"missing verifier",db,nil,startupDiscoveryStub{}},
  {"missing discovery",db,callbackVerifierStub{},nil},
 }
 for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){
  handlers,err:=prepareBrowserOIDCLoginHandlers(context.Background(),cfg,tc.db,tc.verifier,tc.discovery)
  if !errors.Is(err,errBrowserOIDCPrepareHandlers)||handlers.Start!=nil||handlers.Callback!=nil{
   t.Fatalf("incomplete handlers accepted: %v",err)
  }
 })}
}

func TestBrowserOIDCPrepareHandlersValidComposition(t *testing.T){
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 path:=filepath.Join(t.TempDir(),"secret")
 if err:=os.WriteFile(path,[]byte("test-only-secret"),0600);err!=nil{t.Fatal(err)}
 t.Setenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE",path)
 handlers,err:=prepareBrowserOIDCLoginHandlers(context.Background(),cfg,&sql.DB{},callbackVerifierStub{},startupDiscoveryStub{})
 if err!=nil{t.Fatal(err)}
 if handlers.Start==nil||handlers.Callback==nil{t.Fatal("missing unmounted handlers")}
}

func TestBrowserOIDCPrepareHandlersRejectsDiscoveryFailure(t *testing.T){
 setValidBrowserOIDCEnv(t)
 cfg,err:=loadBrowserOIDCConfig()
 if err!=nil{t.Fatal(err)}
 path:=filepath.Join(t.TempDir(),"secret")
 if err:=os.WriteFile(path,[]byte("test-only-secret"),0600);err!=nil{t.Fatal(err)}
 t.Setenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE",path)
 handlers,err:=prepareBrowserOIDCLoginHandlers(context.Background(),cfg,&sql.DB{},callbackVerifierStub{},startupDiscoveryStub{err:errors.New("discovery failed")})
 if !errors.Is(err,errBrowserOIDCPrepareHandlers)||handlers.Start!=nil||handlers.Callback!=nil{
  t.Fatalf("discovery failure accepted: %v",err)
 }
}
