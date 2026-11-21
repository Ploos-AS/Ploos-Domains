package domains

import (
 "errors"
 "sort"
 "strings"
 "time"
)

var ErrInvalidPortfolioEntry = errors.New("invalid domain portfolio entry")

// PortfolioStatus describes Ploos' knowledge of a domain, not a registry
// guarantee. Unknown is the safe default for imported/unverified records.
type PortfolioStatus string
const (
 PortfolioUnknown PortfolioStatus = "unknown"
 PortfolioActive PortfolioStatus = "active"
 PortfolioPending PortfolioStatus = "pending"
 PortfolioExpired PortfolioStatus = "expired"
 PortfolioTransferred PortfolioStatus = "transferred"
)

type PortfolioEntry struct {
 OrganizationID string `json:"organization_id"`
 Domain string `json:"domain"`
 Registrar string `json:"registrar"`
 Status PortfolioStatus `json:"status"`
 ExpiresAt *time.Time `json:"expires_at,omitempty"`
 LastVerifiedAt *time.Time `json:"last_verified_at,omitempty"`
}

// ValidatePortfolioEntry validates a locally tracked domain. It does not
// imply domain ownership, registration, or registrar API confirmation.
func ValidatePortfolioEntry(e PortfolioEntry) error {
 if strings.TrimSpace(e.OrganizationID)=="" || e.OrganizationID!=strings.TrimSpace(e.OrganizationID) ||
  strings.TrimSpace(e.Registrar)=="" || e.Registrar!=strings.TrimSpace(e.Registrar) {return ErrInvalidPortfolioEntry}
 name,err:=CanonicalName(e.Domain)
 if err!=nil || name!=e.Domain || !strings.Contains(name,".") {return ErrInvalidPortfolioEntry}
 switch e.Status {
 case PortfolioUnknown,PortfolioActive,PortfolioPending,PortfolioExpired,PortfolioTransferred:
 default: return ErrInvalidPortfolioEntry
 }
 if e.ExpiresAt!=nil && e.ExpiresAt.IsZero() {return ErrInvalidPortfolioEntry}
 if e.LastVerifiedAt!=nil && e.LastVerifiedAt.IsZero() {return ErrInvalidPortfolioEntry}
 return nil
}

// FilterPortfolio returns a stable, domain-sorted copy for exactly one
// organization. Caller authorization MUST be checked before calling it.
func FilterPortfolio(entries []PortfolioEntry,organizationID string)([]PortfolioEntry,error){
 if strings.TrimSpace(organizationID)=="" || organizationID!=strings.TrimSpace(organizationID) {return nil,ErrInvalidPortfolioEntry}
 result:=make([]PortfolioEntry,0)
 for _,e:=range entries {
  if err:=ValidatePortfolioEntry(e);err!=nil {return nil,err}
  if e.OrganizationID==organizationID {
   copy:=e
   if e.ExpiresAt!=nil {v:=*e.ExpiresAt;copy.ExpiresAt=&v}
   if e.LastVerifiedAt!=nil {v:=*e.LastVerifiedAt;copy.LastVerifiedAt=&v}
   result=append(result,copy)
  }
 }
 sort.Slice(result,func(i,j int)bool{return result[i].Domain<result[j].Domain})
 return result,nil
}
