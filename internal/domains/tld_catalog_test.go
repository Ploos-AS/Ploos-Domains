package domains

import "testing"

func TestTLDCatalogCapabilities(t *testing.T){
 c,err:=NewTLDCatalog([]TLDConfiguration{
  {TLD:"no",Enabled:true,Operations:[]TLDOperation{OperationAvailability,OperationRegistration}},
  {TLD:"com",Enabled:true,Operations:[]TLDOperation{OperationAvailability}},
  {TLD:"net",Enabled:false,Operations:[]TLDOperation{OperationAvailability,OperationRegistration}},
 })
 if err!=nil {t.Fatal(err)}
 if !c.Supports("no",OperationRegistration) || c.Supports("com",OperationRegistration) || c.Supports("net",OperationRegistration) || c.Supports("org",OperationAvailability) {t.Fatal("unexpected capability")}
 r,err:=c.SearchRegistry()
 if err!=nil {t.Fatal(err)}
 for _,domain:=range []string{"example.no","example.com"} {if _,err:=r.CanonicalQuery(domain);err!=nil {t.Errorf("%s: %v",domain,err)}}
 for _,domain:=range []string{"example.net","example.org"} {if _,err:=r.CanonicalQuery(domain);err==nil {t.Errorf("unexpected support: %s",domain)}}
 cfg,ok:=c.Policy("no")
 if !ok {t.Fatal("missing policy")}
 cfg.Operations[0]=OperationRenewal
 if !c.Supports("no",OperationAvailability){t.Fatal("mutated catalog via returned slice")}
}
func TestTLDCatalogRejectsInvalidPolicies(t *testing.T){
 bad:=[][]TLDConfiguration{
  {{TLD:"NO"}},{{TLD:"co.uk"}},{{TLD:"-bad"}},{{TLD:"no"},{TLD:"no"}},
  {{TLD:"no",Operations:[]TLDOperation{"unsupported"}}},
  {{TLD:"no",Operations:[]TLDOperation{OperationAvailability,OperationAvailability}}},
 }
 for _,cfg:=range bad {if _,err:=NewTLDCatalog(cfg);err==nil{t.Errorf("accepted %+v",cfg)}}
}
