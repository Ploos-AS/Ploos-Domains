package domains

import (
 "context"
 "errors"
 "testing"
)

type stubAvailability struct { result SearchResult; err error; calls int }
func (s *stubAvailability) CheckAvailability(_ context.Context, _ string) (SearchResult,error) {
 s.calls++
 return s.result,s.err
}
func TestCanonicalNOQuery(t *testing.T) {
 for _, tc := range []struct{ input, want string }{
  {"Example.NO", "example.no"}, {"test-name.no.", "test-name.no"},
 } {
  got,err:=CanonicalNOQuery(tc.input)
  if err!=nil || got!=tc.want { t.Fatalf("%q: got %q, %v",tc.input,got,err) }
 }
 for _,input:=range []string{"", "foo.com", "foo.bar.no", "foo..no", "-foo.no", "foo.no.evil", "https://foo.no", "foo.no/path", "føø.no", " foo.no", "foo.no "} {
  if _,err:=CanonicalNOQuery(input);err==nil { t.Errorf("accepted %q",input) }
 }
}
func TestSearchNO(t *testing.T) {
 ctx:=context.Background()
 p:=&stubAvailability{result:SearchResult{Domain:"example.no",Status:Available,Source:"test-registrar"}}
 result,err:=SearchNO(ctx,p,"EXAMPLE.NO")
 if err!=nil || result.Status!=Available || p.calls!=1 { t.Fatalf("result=%+v err=%v calls=%d",result,err,p.calls) }
 if _,err:=SearchNO(ctx,p,"example.com");!errors.Is(err,ErrInvalidSearch) { t.Fatalf("invalid search: %v",err) }
 if p.calls!=1 { t.Fatal("invalid search reached provider") }
 for _,bad:=range []SearchResult{
  {Domain:"other.no",Status:Available,Source:"test"},
  {Domain:"example.no",Status:"maybe",Source:"test"},
  {Domain:"example.no",Status:Available},
 } {
  p.result=bad
  if _,err:=SearchNO(ctx,p,"example.no");!errors.Is(err,ErrAvailabilityUnavailable) { t.Errorf("accepted bad result %+v: %v",bad,err) }
 }
 p.err=errors.New("outage")
 if _,err:=SearchNO(ctx,p,"example.no");!errors.Is(err,ErrAvailabilityUnavailable) { t.Fatalf("outage not surfaced: %v",err) }
 if _,err:=SearchNO(ctx,nil,"example.no");!errors.Is(err,ErrAvailabilityUnavailable) { t.Fatalf("nil provider: %v",err) }
}
