package main

import (
 "context"
 "errors"
 "os"
 "testing"
)

func TestBrowserOIDCLoginStartupFailsClosed(t *testing.T){
 for _,value:=range []string{"true","1","TRUE","enabled","invalid"}{
  t.Run(value,func(t *testing.T){
   t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED",value)
   t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
   handler,closeDB,err:=buildHandler(context.Background())
   if !errors.Is(err,errBrowserLoginNotConfigured)||handler!=nil||closeDB!=nil{
    t.Fatalf("unsafe opt-in: handlerPresent=%t closeDBPresent=%t err=%v",handler!=nil,closeDB!=nil,err)
   }
  })
 }
}
func TestBrowserOIDCLoginDisabledByDefault(t *testing.T){
 for _,value:=range []string{"","false"}{
  t.Run("disabled-"+value,func(t *testing.T){
   t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED",value)
   t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
   handler,closeDB,err:=buildHandler(context.Background())
   if err!=nil||handler==nil||closeDB==nil{t.Fatalf("disabled startup: %v",err)}
   if err:=closeDB();err!=nil{t.Fatal(err)}
  })
 }
 _=os.ErrNotExist
}
