package domains

import (
 "context"
 "errors"
 "testing"
)

func TestMultiTLDRegistry(t *testing.T) {
 r,err:=NewTLDRegistry([]TLDPolicy{{TLD:"no",Enabled:true},{TLD:"com",Enabled:true},{TLD:"net",Enabled:false}})
 if err!=nil { t.Fatal(err) }
 for _,tc:=range []struct{in,want string}{{"Example.NO","example.no"},{"Brand.COM.","brand.com"}} {
  got,err:=r.CanonicalQuery(tc.in)
  if err!=nil || got!=tc.want { t.Errorf("%q: got %q %v",tc.in,got,err) }
 }
 for _,in:=range []string{"brand.net","brand.org"} {
  if _,err:=r.CanonicalQuery(in);!errors.Is(err,ErrUnsupportedTLD){t.Errorf("%q: expected unsupported TLD, got %v",in,err)}
 }
 for _,in:=range []string{"x.co.uk","a..no","https://a.no","føø.no"," a.no"} {
  if _,err:=r.CanonicalQuery(in);err==nil{t.Errorf("accepted %q",in)}
 }
 for _,policies:=range [][]TLDPolicy{
  {{TLD:"no"},{TLD:"no"}},{{TLD:"NO"}},{{TLD:"co.uk"}},{{TLD:""}},{{TLD:"-bad"}},
 } {
  if _,err:=NewTLDRegistry(policies);err==nil{t.Errorf("accepted policies %+v",policies)}
 }
}
func TestMultiTLDSearch(t *testing.T) {
 r,err:=NewTLDRegistry([]TLDPolicy{{TLD:"com",Enabled:true}})
 if err!=nil { t.Fatal(err) }
 p:=&stubAvailability{result:SearchResult{Domain:"example.com",Status:Unknown,Source:"mock"}}
 got,err:=SearchDomain(context.Background(),r,p,"EXAMPLE.COM")
 if err!=nil || got.Status!=Unknown { t.Fatalf("got %+v err %v",got,err) }
 if _,err:=SearchDomain(context.Background(),r,p,"example.org");!errors.Is(err,ErrUnsupportedTLD){t.Fatalf("unsupported: %v",err)}
 if p.calls!=1 { t.Fatal("unsupported TLD reached provider") }
 p.result=SearchResult{Domain:"wrong.com",Status:Available,Source:"mock"}
 if _,err:=SearchDomain(context.Background(),r,p,"example.com");!errors.Is(err,ErrAvailabilityUnavailable){t.Fatalf("mismatched response: %v",err)}
 if _,err:=SearchDomain(context.Background(),r,nil,"example.com");!errors.Is(err,ErrAvailabilityUnavailable){t.Fatalf("missing provider: %v",err)}
}
