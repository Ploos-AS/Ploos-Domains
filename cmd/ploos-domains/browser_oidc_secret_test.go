package main

import (
 "errors"
 "os"
 "path/filepath"
 "testing"
)

func TestBrowserOIDCClientSecretFile(t *testing.T){
 dir:=t.TempDir()
 path:=filepath.Join(dir,"secret")
 if err:=os.WriteFile(path,[]byte("test-secret-value"),0600);err!=nil{t.Fatal(err)}
 got,err:=loadBrowserOIDCClientSecret(path)
 if err!=nil||got!="test-secret-value"{t.Fatalf("valid secret rejected: %v",err)}
 if err:=os.Chmod(path,0644);err!=nil{t.Fatal(err)}
 if _,err:=loadBrowserOIDCClientSecret(path);!errors.Is(err,errBrowserOIDCSecret){t.Fatalf("world-readable secret accepted: %v",err)}
 if err:=os.Chmod(path,0600);err!=nil{t.Fatal(err)}
 link:=filepath.Join(dir,"link")
 if err:=os.Symlink(path,link);err!=nil{t.Fatal(err)}
 if _,err:=loadBrowserOIDCClientSecret(link);!errors.Is(err,errBrowserOIDCSecret){t.Fatalf("symlink accepted: %v",err)}
 if _,err:=loadBrowserOIDCClientSecret("relative");!errors.Is(err,errBrowserOIDCSecret){t.Fatalf("relative path accepted: %v",err)}
 if err:=os.WriteFile(path,[]byte("secret\n"),0600);err!=nil{t.Fatal(err)}
 if _,err:=loadBrowserOIDCClientSecret(path);!errors.Is(err,errBrowserOIDCSecret){t.Fatalf("newline accepted: %v",err)}
}
