package server

import (
 "context"
 "errors"
 "net/http"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
 "github.com/Ploos-AS/Ploos-Domains/internal/domains"
)

const PortfolioRoute = "GET /api/v1/organizations/{organizationID}/portfolio"

type PortfolioLister interface {
 ListForUser(context.Context,string,string)([]domains.PortfolioEntry,error)
}

// NewPortfolioHandler requires auth.Middleware at the mounting boundary.
// The caller-supplied organization is never used as proof of authorization.
func NewPortfolioHandler(resolver AuditIdentityResolver,lister PortfolioLister)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("X-Content-Type-Options","nosniff")
  w.Header().Set("Cache-Control","no-store")
  if r.Method!=http.MethodGet {w.Header().Set("Allow","GET");http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
  if resolver==nil||lister==nil {http.Error(w,"portfolio unavailable",http.StatusServiceUnavailable);return}
  identity,ok:=auth.FromContext(r.Context())
  if !ok||identity.Validate()!=nil {http.Error(w,"unauthorized",http.StatusUnauthorized);return}
  org:=r.PathValue("organizationID")
  if org=="" {http.Error(w,"invalid organization",http.StatusBadRequest);return}
  if len(r.URL.Query())!=0 {http.Error(w,"unsupported query",http.StatusBadRequest);return}
  user,err:=resolver.ResolveUserID(r.Context(),identity)
  if errors.Is(err,auth.ErrIdentityNotLinked){http.Error(w,"forbidden",http.StatusForbidden);return}
  if err!=nil {http.Error(w,"portfolio unavailable",http.StatusServiceUnavailable);return}
  entries,err:=lister.ListForUser(r.Context(),user,org)
  if errors.Is(err,domains.ErrInvalidPortfolioEntry){http.Error(w,"invalid organization",http.StatusBadRequest);return}
  if err!=nil {http.Error(w,"portfolio unavailable",http.StatusServiceUnavailable);return}
  if entries==nil {entries=[]domains.PortfolioEntry{}}
  jsonResponse(w,http.StatusOK,map[string]any{"domains":entries})
 })
}

// NewHandlerWithPortfolio mounts the portfolio route only when all trusted
// dependencies are supplied. No token validation is bypassed.
func NewHandlerWithPortfolio(verifier auth.Verifier,resolver AuditIdentityResolver,lister PortfolioLister)http.Handler{
 base:=NewHandler()
 if verifier==nil||resolver==nil||lister==nil{return base}
 mux:=http.NewServeMux()
 limiter:=newAuditLimiter(60,60,time.Now)
 mux.Handle("GET /api/v1/organizations/{organizationID}/portfolio",limiter.wrap(auth.Middleware(verifier,NewPortfolioHandler(resolver,lister))))
 mux.Handle("/api/v1/organizations/{organizationID}/portfolio",limiter.wrap(auth.Middleware(verifier,NewPortfolioHandler(resolver,lister))))
 mux.Handle("/",base)
 return mux
}
