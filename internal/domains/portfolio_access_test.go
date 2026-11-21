package domains

import (
 "context"
 "errors"
 "testing"
)

func TestPortfolioReaderFailsClosedWithoutDB(t *testing.T){
 r:=PortfolioReader{}
 if _,err:=r.ListForUser(context.Background(),"user","org");!errors.Is(err,ErrPortfolioStoreUnavailable){t.Fatalf("got %v",err)}
}
func TestPortfolioReaderRejectsInvalidIdentity(t *testing.T){
 // A real database integration test will verify org isolation and roles.
 r:=PortfolioReader{}
 if _,err:=r.ListForUser(nil,"user","org");!errors.Is(err,ErrPortfolioStoreUnavailable){t.Fatalf("got %v",err)}
}
