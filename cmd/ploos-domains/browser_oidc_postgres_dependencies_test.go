package main

import (
 "database/sql"
 "errors"
 "testing"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

func TestBrowserOIDCPostgresDependenciesRejectMissingInputs(t *testing.T){
 if _,err:=prepareBrowserOIDCPostgresDependencies(nil,callbackVerifierStub{});!errors.Is(err,errBrowserOIDCPostgresDependencies){
  t.Fatalf("nil database accepted: %v",err)
 }
 if _,err:=prepareBrowserOIDCPostgresDependencies(&sql.DB{},nil);!errors.Is(err,errBrowserOIDCPostgresDependencies){
  t.Fatalf("nil verifier accepted: %v",err)
 }
}

func TestBrowserOIDCPostgresDependenciesUseOneDatabase(t *testing.T){
 db:=&sql.DB{}
 writer,callback,err:=prepareBrowserOIDCPostgresLoginStores(db,callbackVerifierStub{})
 if err!=nil{t.Fatal(err)}
 start,ok:=writer.(auth.BrowserLoginStore)
 if !ok||start.DB!=db{t.Fatal("start writer not bound to runtime database")}
 tx,ok:=callback.Transactions.(auth.BrowserLoginStore)
 if !ok||tx.DB!=db{t.Fatal("callback consumer not bound to runtime database")}
 users,ok:=callback.Users.(auth.IdentityRepository)
 if !ok||users.DB!=db{t.Fatal("identity resolver not bound to runtime database")}
 sessions,ok:=callback.Sessions.(server.BrowserSessionStore)
 if !ok||sessions.DB!=db{t.Fatal("session store not bound to runtime database")}
 if callback.Verifier==nil||callback.Exchanger!=nil{
  t.Fatal("verifier missing or exchanger wired before secret validation")
 }
}
