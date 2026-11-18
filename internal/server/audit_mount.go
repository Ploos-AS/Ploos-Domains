package server

import (
 "net/http"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

// NewHandlerWithAudit preserves the public health and API foundation routes,
// while mounting the audit endpoint behind mandatory OIDC middleware.
// Never call with nil dependencies.
func NewHandlerWithAudit(verifier auth.Verifier,resolver AuditIdentityResolver,lister AuditLister)http.Handler{
 base:=NewHandler()
 if verifier==nil||resolver==nil||lister==nil{return base}
 mux:=http.NewServeMux()
 mux.Handle(AuditRoute,auth.Middleware(verifier,NewAuditHandler(resolver,lister)))
 mux.Handle("/",base)
 return mux
}
