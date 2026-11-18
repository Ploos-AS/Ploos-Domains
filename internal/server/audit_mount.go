package server

import (
 "net/http"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

// NewHandlerWithAudit preserves the public health and API foundation routes,
// while mounting the audit endpoint behind mandatory OIDC middleware.
// Never call with nil dependencies.
func NewHandlerWithAudit(verifier auth.Verifier,resolver AuditIdentityResolver,lister AuditLister)http.Handler{
 base:=NewHandler()
 if verifier==nil||resolver==nil||lister==nil{return base}
 mux:=http.NewServeMux()
 limiter:=newAuditLimiter(60,60,time.Now)
 mux.Handle(AuditRoute,limiter.wrap(auth.Middleware(verifier,NewAuditHandler(resolver,lister))))
 mux.Handle("/",base)
 return mux
}
