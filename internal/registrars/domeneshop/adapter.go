// Package domeneshop defines a fail-closed adapter boundary for a future
// authorized Domeneshop reseller integration. No HTTP/API operations exist yet.
package domeneshop

import (
 "context"
 "errors"

 "github.com/Ploos-AS/Ploos-Domains/internal/domains"
)

var ErrNotConfigured = errors.New("Domeneshop registrar integration is not configured")
var ErrCapabilityUnavailable = errors.New("Domeneshop registrar capability is not enabled")

// Client is deliberately inert until the commercial agreement, credentials,
// API endpoints and permissions are verified and implemented.
type Client struct {
 catalog domains.TLDCatalog
}

// NewClient takes a validated, operator-supplied catalogue. It does not
// assume any TLD is supported by Domeneshop by default.
func NewClient(catalog domains.TLDCatalog) *Client {
 return &Client{catalog:catalog}
}

// CheckAvailability implements domains.AvailabilityProvider but deliberately
// refuses to issue any network call or report availability until a verified
// Domeneshop API implementation is supplied.
func (c *Client) CheckAvailability(_ context.Context, name string)(domains.SearchResult,error){
 if c==nil { return domains.SearchResult{},ErrNotConfigured }
 canonical,err:=domains.CanonicalName(name)
 if err!=nil || canonical!=name { return domains.SearchResult{},domains.ErrInvalidSearch }
 // Capability checks are for the current simple single-label TLD model.
 labels:=splitLabels(name)
 if len(labels)!=2 { return domains.SearchResult{},domains.ErrInvalidSearch }
 if !c.catalog.Supports(labels[1],domains.OperationAvailability) {
  return domains.SearchResult{},ErrCapabilityUnavailable
 }
 return domains.SearchResult{},ErrNotConfigured
}
