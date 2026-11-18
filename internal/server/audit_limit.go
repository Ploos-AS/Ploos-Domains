package server

import (
 "net/http"
 "sync"
 "time"
)

// auditLimiter applies a bounded, process-local token bucket to the audit
// route before OIDC verification. It deliberately does not trust X-Forwarded-For.
// For multi-replica deployments use an additional shared gateway limiter.
type auditLimiter struct {
 mu sync.Mutex
 tokens float64
 capacity float64
 refillPerSecond float64
 last time.Time
 clock func()time.Time
}

func newAuditLimiter(capacity int,perMinute int,clock func()time.Time)*auditLimiter{
 return &auditLimiter{tokens:float64(capacity),capacity:float64(capacity),
  refillPerSecond:float64(perMinute)/60,last:clock(),clock:clock}
}
func (l *auditLimiter) allow()bool{
 l.mu.Lock()
 defer l.mu.Unlock()
 now:=l.clock()
 elapsed:=now.Sub(l.last).Seconds()
 if elapsed>0{
  l.tokens+=elapsed*l.refillPerSecond
  if l.tokens>l.capacity{l.tokens=l.capacity}
  l.last=now
 }
 if l.tokens<1{return false}
 l.tokens--
 return true
}
func (l *auditLimiter) wrap(next http.Handler)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("X-Content-Type-Options","nosniff")
  w.Header().Set("Referrer-Policy","no-referrer")
  if !l.allow(){
   w.Header().Set("Retry-After","60")
   http.Error(w,"too many requests",http.StatusTooManyRequests)
   return
  }
  next.ServeHTTP(w,r)
 })
}
