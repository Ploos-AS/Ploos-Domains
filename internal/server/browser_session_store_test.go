package server

import (
 "context"
 "errors"
 "testing"
 "time"
)

func TestBrowserSessionStoreFailClosed(t *testing.T){
 s:=BrowserSessionStore{}
 if err:=s.Create(context.Background(),"user","token",time.Now(),time.Hour);!errors.Is(err,ErrBrowserSessionUnavailable){t.Fatalf("create: %v",err)}
 if _,err:=s.Resolve(context.Background(),"token");!errors.Is(err,ErrBrowserSessionUnavailable){t.Fatalf("resolve: %v",err)}
 if err:=s.Revoke(context.Background(),"token");!errors.Is(err,ErrBrowserSessionUnavailable){t.Fatalf("revoke: %v",err)}
}
