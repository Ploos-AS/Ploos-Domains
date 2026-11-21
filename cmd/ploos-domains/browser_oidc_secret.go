package main

import (
 "errors"
 "io"
 "os"
 "path/filepath"
 "strings"
)

var errBrowserOIDCSecret=errors.New("invalid browser OIDC client secret file")

// loadBrowserOIDCClientSecret reads a dedicated local secret file. The path
// must be absolute, a regular file, not
// accessible by group/other. Symlinks are rejected to prevent path swaps.
func loadBrowserOIDCClientSecret(path string)(string,error){
 if path==""||!filepath.IsAbs(path)||filepath.Clean(path)!=path{return "",errBrowserOIDCSecret}
 info,err:=os.Lstat(path)
 if err!=nil||!info.Mode().IsRegular()||info.Mode().Perm()&0077!=0||info.Size()<1||info.Size()>4096{
  return "",errBrowserOIDCSecret
 }
 f,err:=os.Open(path)
 if err!=nil{return "",errBrowserOIDCSecret}
 defer f.Close()
 opened,err:=f.Stat()
 if err!=nil||!opened.Mode().IsRegular()||!os.SameFile(info,opened){return "",errBrowserOIDCSecret}
 // The file must not change between metadata validation and read.
 raw:=make([]byte,opened.Size())
 n,err:=io.ReadFull(f,raw)
 if err!=nil||int64(n)!=opened.Size(){return "",errBrowserOIDCSecret}
 secret:=string(raw)
 if strings.TrimSpace(secret)!=secret||strings.ContainsAny(secret,"\x00\r\n")||secret==""{
  return "",errBrowserOIDCSecret
 }
 return secret,nil
}
