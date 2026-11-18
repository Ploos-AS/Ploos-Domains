package main

import (
 "context"
 "errors"
 "testing"
)

func TestAuditPrivilegesRejectNilDatabase(t *testing.T){
 err:=checkAuditPrivileges(context.Background(),nil)
 if !errors.Is(err,errAuditPrivileges){t.Fatalf("expected unsafe privilege error, got %v",err)}
}
