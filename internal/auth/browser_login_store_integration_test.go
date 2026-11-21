//go:build integration

package auth

import (
 "context"
 "crypto/sha256"
 "database/sql"
 "errors"
 "os"
 "testing"
 "time"

 _ "github.com/lib/pq"
)

func TestBrowserLoginStoreSingleUseAndBinding(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL for disposable migrated PostgreSQL")}
 db,err:=sql.Open("postgres",dsn);if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second);defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 store:=BrowserLoginStore{DB:db}
 tx,err:=NewBrowserLoginTransaction();if err!=nil{t.Fatal(err)}
 if err:=store.Create(ctx,tx,"browser-A",time.Minute);err!=nil{t.Fatal(err)}
 t.Cleanup(func(){state:=sha256.Sum256([]byte(tx.State));_,_=db.ExecContext(context.Background(),"DELETE FROM oidc_login_transactions WHERE state_sha256=$1",state[:])})
 if _,err:=store.Consume(ctx,tx.State,"browser-B");!errors.Is(err,ErrLoginTransactionExpired){t.Fatalf("wrong browser: %v",err)}
 got,err:=store.Consume(ctx,tx.State,"browser-A");if err!=nil{t.Fatal(err)}
 if got.Nonce!=tx.Nonce||got.PKCEVerifier!=tx.PKCEVerifier{t.Fatal("wrong transaction secrets")}
 if _,err:=store.Consume(ctx,tx.State,"browser-A");!errors.Is(err,ErrLoginTransactionExpired){t.Fatalf("replay accepted: %v",err)}
}
