package main

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "testing"
)

func TestAuditDisabledByDefault(t *testing.T){
 t.Setenv("PLOOS_AUDIT_API_ENABLED","")
 handler,closeFn,err:=buildHandler(context.Background())
 if err!=nil{t.Fatal(err)}
 defer closeFn()
 req:=httptest.NewRequest(http.MethodGet,"/api/v1/organizations/example/audit-events",nil)
 rr:=httptest.NewRecorder()
 handler.ServeHTTP(rr,req)
 if rr.Code!=http.StatusNotFound{t.Fatalf("disabled route returned %d",rr.Code)}
}
func TestAuditEnabledFailsClosedWithoutConfiguration(t *testing.T){
 t.Setenv("PLOOS_AUDIT_API_ENABLED","true")
 t.Setenv("PLOOS_OIDC_ISSUER","")
 t.Setenv("PLOOS_OIDC_AUDIENCE","")
 t.Setenv("PLOOS_DATABASE_URL","")
 _,_,err:=buildHandler(context.Background())
 if !errors.Is(err,errAuditConfig){t.Fatalf("expected invalid config, got %v",err)}
}
