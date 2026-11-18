//go:build integration

package auth

import (
 "context"
 "database/sql"
 "os"
 "sync"
 "testing"
 "time"

 _ "github.com/lib/pq"
)

func TestProvisioningPostgres(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL")}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{t.Fatal(err)}
 defer db.Close()
 db.SetMaxOpenConns(8)
 ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second)
 defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 id,err:=randomUUID()
 if err!=nil{t.Fatal(err)}
 identity:=Identity{Issuer:"https://test.example.org",Subject:"integration-"+id,Audience:"ploos-domains"}
 p:=UserProvisioner{DB:db}
 var wg sync.WaitGroup
 results:=make(chan string,6)
 errs:=make(chan error,6)
 for i:=0;i<6;i++ {
  wg.Add(1)
  go func(){defer wg.Done();uid,e:=p.Provision(ctx,identity);results<-uid;errs<-e}()
 }
 wg.Wait();close(results);close(errs)
 var first string
 for e:=range errs {if e!=nil{t.Fatal(e)}}
 for uid:=range results {if first==""{first=uid};if uid!=first{t.Fatalf("multiple users: %s and %s",first,uid)}}
 if first==""{t.Fatal("missing user")}
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",first)})
 var count int
 if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM user_oidc_identities WHERE issuer=$1 AND subject=$2",identity.Issuer,identity.Subject).Scan(&count);err!=nil{t.Fatal(err)}
 if count!=1{t.Fatalf("identity links=%d",count)}
 if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM memberships WHERE user_id=$1",first).Scan(&count);err!=nil{t.Fatal(err)}
 if count!=0{t.Fatalf("unexpected memberships=%d",count)}
}
