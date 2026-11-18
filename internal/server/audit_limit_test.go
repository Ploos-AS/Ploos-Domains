package server

import (
 "net/http"
 "net/http/httptest"
 "testing"
 "time"
)

func TestAuditLimiterExhaustionAndRefill(t *testing.T){
 now:=time.Date(2026,11,18,12,0,0,0,time.UTC)
 limiter:=newAuditLimiter(2,60,func()time.Time{return now})
 called:=0
 handler:=limiter.wrap(http.HandlerFunc(func(w http.ResponseWriter,_ *http.Request){called++;w.WriteHeader(http.StatusOK)}))
 request:=func() *httptest.ResponseRecorder{
  rr:=httptest.NewRecorder()
  handler.ServeHTTP(rr,httptest.NewRequest(http.MethodGet,"/audit",nil))
  if rr.Header().Get("Cache-Control")!="no-store"||rr.Header().Get("X-Content-Type-Options")!="nosniff"||rr.Header().Get("Referrer-Policy")!="no-referrer"{t.Fatal("missing security headers")}
  return rr
 }
 if request().Code!=200||request().Code!=200{t.Fatal("initial capacity denied")}
 limited:=request()
 if limited.Code!=429||limited.Header().Get("Retry-After")==""{t.Fatalf("rate limit=%d",limited.Code)}
 if called!=2{t.Fatalf("downstream called %d times",called)}
 now=now.Add(2*time.Second)
 if request().Code!=200{t.Fatal("limiter did not refill")}
}
