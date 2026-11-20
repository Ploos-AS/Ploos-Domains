package domeneshop

import (
 "errors"
 "net/url"
 "strings"
 "time"
)

var ErrInvalidConfiguration = errors.New("invalid Domeneshop client configuration")

// APIConfig describes connection parameters supplied only after official
// reseller API documentation and authorization have been verified.
// This does not enable any API operations.
type APIConfig struct {
 BaseURL string
 Timeout time.Duration
}

// Validate rejects insecure and ambiguous endpoints. The URL is not an
// assumed Domeneshop endpoint; it must be supplied by the operator.
func (c APIConfig) Validate() error {
 u,err:=url.Parse(c.BaseURL)
 if err!=nil || u==nil || u.Scheme!="https" || u.Hostname()=="" ||
  u.User!=nil || u.RawQuery!="" || u.Fragment!="" ||
  u.Opaque!="" || u.ForceQuery || strings.TrimSpace(c.BaseURL)!=c.BaseURL ||
  u.Port()!="" && u.Port()!="443" || c.Timeout<=0 || c.Timeout>30*time.Second {
  return ErrInvalidConfiguration
 }
 if u.Path!="" && u.Path!="/" { return ErrInvalidConfiguration }
 return nil
}
