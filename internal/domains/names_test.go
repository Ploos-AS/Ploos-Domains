package domains

import "testing"

func TestCanonicalName(t *testing.T) {
 good := map[string]string{
  "SIRDAL.SENIORNETT.NO.":"sirdal.seniornett.no",
  "a-b.example.org":"a-b.example.org",
  "x.example":"x.example",
 }
 for input,want := range good {
  got,err:=CanonicalName(input)
  if err!=nil || got!=want { t.Errorf("%q: got %q, %v; want %q",input,got,err,want) }
 }
 bad:=[]string{"", ".", "a..example.org", "-x.example.org", "x-.example.org", "a_b.example.org", "*.example.org", "https://example.org", "sirdål.no", "foo.example.org.."}
 for _,input:=range bad {
  if got,err:=CanonicalName(input);err==nil { t.Errorf("accepted %q as %q",input,got) }
 }
}

func TestAllocationName(t *testing.T) {
 for _,name:=range []string{"sirdal.seniornett.no","KURS.SIRDAL.SENIORNETT.NO."} {
  if _,err:=AllocationName(name,"seniornett.no");err!=nil { t.Errorf("%q rejected: %v",name,err) }
 }
 for _,name:=range []string{"seniornett.no","evil-seniornett.no","sirdal.seniornett.no.evil.org","other.no"} {
  if _,err:=AllocationName(name,"seniornett.no");err==nil { t.Errorf("accepted %q",name) }
 }
}
