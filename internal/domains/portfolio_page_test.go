package domains

import (
 "context"
 "errors"
 "testing"
)

func TestPortfolioPageFailsClosed(t *testing.T){
 r:=PortfolioReader{}
 _,err:=r.ListPageForUser(context.Background(),"user","org",50,"")
 if !errors.Is(err,ErrPortfolioStoreUnavailable){t.Fatalf("got %v",err)}
}
