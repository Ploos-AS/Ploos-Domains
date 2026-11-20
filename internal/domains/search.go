package domains

import (
 "context"
 "errors"
 "fmt"
 "strings"
)

var ErrInvalidSearch = errors.New("invalid domain search")
var ErrAvailabilityUnavailable = errors.New("domain availability provider unavailable")

// Availability is deliberately not a boolean: unknown must never be sold as available.
type Availability string
const (
 Available Availability = "available"
 Unavailable Availability = "unavailable"
 Unknown Availability = "unknown"
)

type SearchResult struct {
 Domain string `json:"domain"`
 Status Availability `json:"status"`
 // Source identifies the registrar/reseller adapter used for the lookup.
 Source string `json:"source"`
}

// AvailabilityProvider must use an authorized registrar/reseller availability
// service. DNS resolution, RDAP absence and HTTP reachability are NOT proof
// that a domain can be registered.
type AvailabilityProvider interface {
 CheckAvailability(ctx context.Context, domain string) (SearchResult, error)
}

// CanonicalNOQuery validates a full ASCII .no domain name, not a subdomain.
// Unicode IDNs require a separate standards-compliant IDNA conversion layer.
func CanonicalNOQuery(query string) (string, error) {
 if query != strings.TrimSpace(query) { return "", ErrInvalidSearch }
 name, err := CanonicalName(query)
 if err != nil { return "", fmt.Errorf("%w: %v", ErrInvalidSearch, err) }
 parts := strings.Split(name, ".")
 if len(parts) != 2 || parts[1] != "no" || len(parts[0]) == 0 { return "", ErrInvalidSearch }
 return name, nil
}

func SearchNO(ctx context.Context, provider AvailabilityProvider, query string) (SearchResult, error) {
 domain, err := CanonicalNOQuery(query)
 if err != nil { return SearchResult{}, err }
 if provider == nil { return SearchResult{}, ErrAvailabilityUnavailable }
 result, err := provider.CheckAvailability(ctx, domain)
 if err != nil { return SearchResult{}, fmt.Errorf("%w: %v", ErrAvailabilityUnavailable, err) }
 if result.Domain != domain || result.Source == "" { return SearchResult{}, ErrAvailabilityUnavailable }
 switch result.Status {
 case Available, Unavailable, Unknown:
  return result, nil
 default:
  return SearchResult{}, ErrAvailabilityUnavailable
 }
}
