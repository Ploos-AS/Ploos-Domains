package domains

import (
 "context"
 "errors"
 "testing"
)

func TestPortfolioPostgresRequiresDatabase(t *testing.T){
 repo:=PortfolioPostgresRepository{}
 if _,err:=repo.List(context.Background(),"org");!errors.Is(err,ErrPortfolioStoreUnavailable){t.Fatalf("list: %v",err)}
 if err:=repo.Put(context.Background(),PortfolioEntry{});!errors.Is(err,ErrPortfolioStoreUnavailable){t.Fatalf("put: %v",err)}
}
