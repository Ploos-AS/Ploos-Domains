package domains

import (
 "context"
 "errors"
 "fmt"
 "strings"
)

var ErrUnsupportedTLD = errors.New("unsupported top-level domain")

// TLDPolicy is an explicitly configured provider capability, not a claim that
// any particular registrar currently permits registration for this TLD.
type TLDPolicy struct {
 TLD string
 Enabled bool
}

// TLDRegistry is a fail-closed list of extensions approved for search.
// Populate it only after confirming the chosen registrar's actual capabilities.
type TLDRegistry struct { policies map[string]TLDPolicy }

func NewTLDRegistry(policies []TLDPolicy) (TLDRegistry, error) {
 registry := TLDRegistry{policies:make(map[string]TLDPolicy,len(policies))}
 for _,p:=range policies {
  tld:=strings.ToLower(p.TLD)
  if tld=="" || strings.Contains(tld,".") { return TLDRegistry{},ErrUnsupportedTLD }
  canonical,err:=CanonicalName("example."+tld)
  if err!=nil || canonical!="example."+tld || tld!=p.TLD { return TLDRegistry{},ErrUnsupportedTLD }
  if _,exists:=registry.policies[tld];exists { return TLDRegistry{},ErrUnsupportedTLD }
  registry.policies[tld]=p
 }
 return registry,nil
}

// CanonicalQuery accepts only a single second-level ASCII name below an
// explicitly enabled TLD. This intentionally excludes Unicode IDNs and
// multi-label public suffixes until dedicated policy support is implemented.
func (r TLDRegistry) CanonicalQuery(query string) (string,error) {
 if query!=strings.TrimSpace(query) { return "",ErrInvalidSearch }
 name,err:=CanonicalName(query)
 if err!=nil { return "",fmt.Errorf("%w: %v",ErrInvalidSearch,err) }
 labels:=strings.Split(name,".")
 if len(labels)!=2 { return "",ErrInvalidSearch }
 policy,ok:=r.policies[labels[1]]
 if !ok || !policy.Enabled { return "",ErrUnsupportedTLD }
 return name,nil
}

// SearchDomain is a provider-neutral, fail-closed multi-TLD availability
// lookup. Provider results are advisory and never a domain reservation.
func SearchDomain(ctx context.Context,registry TLDRegistry,provider AvailabilityProvider,query string)(SearchResult,error){
 name,err:=registry.CanonicalQuery(query)
 if err!=nil { return SearchResult{},err }
 if provider==nil { return SearchResult{},ErrAvailabilityUnavailable }
 result,err:=provider.CheckAvailability(ctx,name)
 if err!=nil { return SearchResult{},fmt.Errorf("%w: %v",ErrAvailabilityUnavailable,err) }
 if result.Domain!=name || result.Source=="" { return SearchResult{},ErrAvailabilityUnavailable }
 switch result.Status {
 case Available,Unavailable,Unknown: return result,nil
 default: return SearchResult{},ErrAvailabilityUnavailable
 }
}
