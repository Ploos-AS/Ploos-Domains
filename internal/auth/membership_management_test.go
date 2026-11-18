package auth

import (
 "context"
 "errors"
 "testing"
)

func TestMembershipManagerRejectsUnsafeRequests(t *testing.T){
 m:=MembershipManager{}
 for _,role:=range []string{"owner","unit_admin","superadmin",""}{
  if _,ok:=interface{}(m.Grant(context.Background(),"actor","org","target",role)).(error);!ok{t.Fatal("expected error")}
  if err:=m.Grant(context.Background(),"actor","org","target",role);!errors.Is(err,ErrInvalidIdentity){t.Fatalf("role %q: %v",role,err)}
 }
 if err:=m.Grant(context.Background(),"actor","org","actor","viewer");!errors.Is(err,ErrInvalidIdentity){t.Fatal(err)}
}
