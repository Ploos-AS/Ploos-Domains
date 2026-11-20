package domeneshop

import (
 "context"
 "errors"
 "testing"

 "github.com/Ploos-AS/Ploos-Domains/internal/domains"
)

func TestAdapterFailsClosed(t *testing.T) {
 catalog,err:=domains.NewTLDCatalog([]domains.TLDConfiguration{
  {TLD:"no",Enabled:true,Operations:[]domains.TLDOperation{domains.OperationAvailability}},
  {TLD:"com",Enabled:false,Operations:[]domains.TLDOperation{domains.OperationAvailability}},
 })
 if err!=nil { t.Fatal(err) }
 client:=NewClient(catalog)
 for _,tc:=range []struct{name string;want error}{
  {"example.no",ErrNotConfigured},
  {"example.com",ErrCapabilityUnavailable},
  {"example.org",ErrCapabilityUnavailable},
  {"foo.bar.no",domains.ErrInvalidSearch},
  {"EXAMPLE.NO",domains.ErrInvalidSearch},
 } {
  result,err:=client.CheckAvailability(context.Background(),tc.name)
  if !errors.Is(err,tc.want) { t.Errorf("%q: expected %v, got %v",tc.name,tc.want,err) }
  if result.Status==domains.Available { t.Errorf("%q incorrectly available",tc.name) }
 }
 var nilClient *Client
 if _,err:=nilClient.CheckAvailability(context.Background(),"example.no");!errors.Is(err,ErrNotConfigured) {t.Fatalf("nil client: %v",err)}
}
