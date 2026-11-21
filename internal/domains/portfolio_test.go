package domains

import (
 "errors"
 "testing"
 "time"
)

func TestPortfolioIsolationAndOrdering(t *testing.T){
 expiry:=time.Date(2027,1,2,0,0,0,0,time.UTC)
 entries:=[]PortfolioEntry{
  {OrganizationID:"org-a",Domain:"z.no",Registrar:"external",Status:PortfolioUnknown,ExpiresAt:&expiry},
  {OrganizationID:"org-b",Domain:"secret.com",Registrar:"external",Status:PortfolioActive},
  {OrganizationID:"org-a",Domain:"a.com",Registrar:"domeneshop",Status:PortfolioPending},
 }
 result,err:=FilterPortfolio(entries,"org-a")
 if err!=nil {t.Fatal(err)}
 if len(result)!=2 || result[0].Domain!="a.com" || result[1].Domain!="z.no" {t.Fatalf("unexpected results %+v",result)}
 *result[1].ExpiresAt=expiry.Add(time.Hour)
 if !entries[0].ExpiresAt.Equal(expiry) {t.Fatal("returned pointer aliases original")}
 if _,err:=FilterPortfolio(entries," ");!errors.Is(err,ErrInvalidPortfolioEntry){t.Fatal(err)}
}
func TestPortfolioValidationFailsClosed(t *testing.T){
 base:=PortfolioEntry{OrganizationID:"org-a",Domain:"example.no",Registrar:"external",Status:PortfolioUnknown}
 if err:=ValidatePortfolioEntry(base);err!=nil {t.Fatal(err)}
 for _,change:=range []func(*PortfolioEntry){
  func(e *PortfolioEntry){e.OrganizationID=""},
  func(e *PortfolioEntry){e.Domain="EXAMPLE.NO"},
  func(e *PortfolioEntry){e.Domain="example"},
  func(e *PortfolioEntry){e.Registrar=""},
  func(e *PortfolioEntry){e.Status="registered"},
 }{
  e:=base;change(&e)
  if err:=ValidatePortfolioEntry(e);!errors.Is(err,ErrInvalidPortfolioEntry){t.Errorf("accepted %+v: %v",e,err)}
 }
}
