// Package auth defines a fail-closed boundary for verified OIDC identities.
// JWT validation is performed by a separate, trusted implementation of Verifier.
package auth

import (
 "context"
 "errors"
 "net/http"
 "strings"
)

var (
 ErrUnauthenticated = errors.New("unauthenticated")
 ErrInvalidIdentity = errors.New("invalid verified identity")
)

// Identity is only constructed from a successfully verified token.
// Subject must be interpreted together with Issuer to avoid cross-provider collisions.
type Identity struct {
 Issuer string
 Subject string
 Audience string
}

func (i Identity) Validate() error {
 if i.Issuer=="" || i.Subject=="" || i.Audience=="" {return ErrInvalidIdentity}
 return nil
}

// Verifier MUST verify signature, trusted issuer, intended audience, expiry,
// not-before and allowed signing algorithm before returning an Identity.
// Never implement this by merely decoding a JWT payload.
type Verifier interface {
 VerifyToken(ctx context.Context, rawToken string) (Identity,error)
}

type contextKey struct{}

func FromContext(ctx context.Context)(Identity,bool){
 identity,ok:=ctx.Value(contextKey{}).(Identity)
 return identity,ok
}

// Middleware accepts only a single Bearer token and injects a verified identity.
// Do not configure this with a verifier that skips cryptographic validation.
func Middleware(verifier Verifier,next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if verifier==nil || next==nil {http.Error(w,"authentication unavailable",http.StatusServiceUnavailable);return}
  values:=r.Header.Values("Authorization")
  if len(values)!=1 {http.Error(w,"unauthorized",http.StatusUnauthorized);return}
  parts:=strings.Fields(values[0])
  if len(parts)!=2 || !strings.EqualFold(parts[0],"Bearer") || parts[1]=="" {
   http.Error(w,"unauthorized",http.StatusUnauthorized);return
  }
  identity,err:=verifier.VerifyToken(r.Context(),parts[1])
  if err!=nil || identity.Validate()!=nil {
   http.Error(w,"unauthorized",http.StatusUnauthorized);return
  }
  next.ServeHTTP(w,r.WithContext(context.WithValue(r.Context(),contextKey{},identity)))
 })
}
