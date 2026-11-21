package auth

import (
 "context"
 "errors"
 "testing"
 "time"
)

func TestBrowserLoginStoreFailsClosed(t *testing.T){
 s:=BrowserLoginStore{}
 tx,err:=NewBrowserLoginTransaction();if err!=nil{t.Fatal(err)}
 if err:=s.Create(context.Background(),tx,"binding",5*time.Minute);!errors.Is(err,ErrLoginTransactionUnavailable){t.Fatalf("create: %v",err)}
 if _,err:=s.Consume(context.Background(),tx.State,"binding");!errors.Is(err,ErrLoginTransactionUnavailable){t.Fatalf("consume: %v",err)}
}
