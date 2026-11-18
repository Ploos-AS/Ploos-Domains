package server

import (
 "net/http"
 "net/http/httptest"
 "testing"
)

func TestAuditMountWithoutDependencies(t *testing.T){
 handler:=NewHandlerWithAudit(nil,nil,nil)
 req:=httptest.NewRequest(http.MethodGet,"/api/v1/organizations/example/audit-events",nil)
 rr:=httptest.NewRecorder()
 handler.ServeHTTP(rr,req)
 if rr.Code!=http.StatusNotFound{t.Fatalf("route unexpectedly mounted: %d",rr.Code)}
 req=httptest.NewRequest(http.MethodGet,"/healthz",nil)
 rr=httptest.NewRecorder()
 handler.ServeHTTP(rr,req)
 if rr.Code!=http.StatusOK{t.Fatalf("healthz status=%d",rr.Code)}
}
