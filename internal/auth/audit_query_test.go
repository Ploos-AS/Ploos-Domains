package auth

import (
 "context"
 "errors"
 "testing"
)

func TestAuditQueryRejectsInvalidRequests(t *testing.T){
 q:=AuditQuery{}
 for _,limit:=range []int{0,-1,101}{
  _,err:=q.List(context.Background(),"actor","org",limit,nil)
  if !errors.Is(err,ErrInvalidAuditQuery){t.Fatalf("limit=%d err=%v",limit,err)}
 }
 _,err:=q.List(context.Background(),"actor","org",10,&AuditCursor{ID:"bad"})
 if !errors.Is(err,ErrInvalidAuditQuery){t.Fatal(err)}
}
