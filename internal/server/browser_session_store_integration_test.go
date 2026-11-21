//go:build integration

package server

import (
 "context"
 "database/sql"
 "errors"
 "os"
 "testing"
 "time"

 _ "github.com/lib/pq"
)

func TestBrowserSessionStorePostgresLifecycle(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL for migrated disposable PostgreSQL")}
 db,err:=sql.Open("postgres",dsn);if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second);defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 var userID string
 err=db.QueryRowContext(ctx,`INSERT INTO users(id,external_subject) VALUES(gen_random_uuid(),$1) RETURNING id::text`,"browser-session-test-"+time.Now().Format("20060102150405.000000000")).Scan(&userID)
 if err!=nil{t.Fatal(err)}
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",userID)})
 store:=BrowserSessionStore{DB:db}
 first,_,err:=NewBrowserSessionToken();if err!=nil{t.Fatal(err)}
 second,_,err:=NewBrowserSessionToken();if err!=nil{t.Fatal(err)}
 now:=time.Now().UTC()
 for _,token:=range []string{first,second}{
  if err:=store.Create(ctx,userID,token,now,time.Hour);err!=nil{t.Fatal(err)}
 }
 if got,err:=store.Resolve(ctx,first);err!=nil||got!=userID{t.Fatalf("resolve: %s %v",got,err)}
 var storedRaw int
 if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM browser_sessions WHERE token_sha256=$1",[]byte(first)).Scan(&storedRaw);err!=nil{t.Fatal(err)}
 if storedRaw!=0{t.Fatal("raw token was persisted")}
 if err:=store.Revoke(ctx,first);err!=nil{t.Fatal(err)}
 if _,err:=store.Resolve(ctx,first);!errors.Is(err,ErrBrowserSessionInvalid){t.Fatalf("revoked session: %v",err)}
 if got,err:=store.Resolve(ctx,second);err!=nil||got!=userID{t.Fatalf("other session: %s %v",got,err)}
 if _,err:=db.ExecContext(ctx,"UPDATE browser_sessions SET expires_at=now()-interval '1 second' WHERE user_id=$1",userID);err!=nil{t.Fatal(err)}
 if _,err:=store.Resolve(ctx,second);!errors.Is(err,ErrBrowserSessionInvalid){t.Fatalf("expired session: %v",err)}
}
