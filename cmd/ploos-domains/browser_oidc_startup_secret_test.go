package main

import (
 "context"
 "errors"
 "os"
 "path/filepath"
 "testing"
)

func TestBrowserOIDCOptInRequiresProtectedSecretFile(t *testing.T){
 setValidBrowserOIDCEnv(t)
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED","true")
 t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
 t.Setenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE","")
 h,closeDB,err:=buildHandler(context.Background())
 if !errors.Is(err,errBrowserOIDCSecret)||h!=nil||closeDB!=nil{
  t.Fatalf("missing secret accepted: handler=%t closeDB=%t err=%v",h!=nil,closeDB!=nil,err)
 }
}

func TestBrowserOIDCOptInValidSecretProceedsToDiscovery(t *testing.T){
 setValidBrowserOIDCEnv(t)
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED","true")
 t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
 path:=filepath.Join(t.TempDir(),"oidc-secret")
 if err:=os.WriteFile(path,[]byte("ci-test-secret"),0600);err!=nil{t.Fatal(err)}
 t.Setenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE",path)
 h,closeDB,err:=buildHandler(context.Background())
 if !errors.Is(err,errBrowserOIDCDiscovery)||h!=nil||closeDB!=nil{
  t.Fatalf("startup did not reach discovery: handler=%t closeDB=%t err=%v",h!=nil,closeDB!=nil,err)
 }
}

func TestBrowserOIDCDisabledIgnoresMissingSecret(t *testing.T){
 t.Setenv("PLOOS_BROWSER_OIDC_LOGIN_ENABLED","false")
 t.Setenv("PLOOS_BROWSER_OIDC_CLIENT_SECRET_FILE","/does/not/exist")
 t.Setenv("PLOOS_AUDIT_API_ENABLED","false")
 h,closeDB,err:=buildHandler(context.Background())
 if err!=nil||h==nil||closeDB==nil{t.Fatalf("disabled login rejected: %v",err)}
 if err:=closeDB();err!=nil{t.Fatal(err)}
}
