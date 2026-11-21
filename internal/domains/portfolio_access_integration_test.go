//go:build integration

package domains

import (
 "context"
 "crypto/rand"
 "database/sql"
 "encoding/hex"
 "fmt"
 "os"
 "testing"
 "time"

 _ "github.com/lib/pq"
)

func portfolioUUID(t *testing.T) string {
 t.Helper()
 b:=make([]byte,16)
 if _,err:=rand.Read(b);err!=nil{t.Fatal(err)}
 b[6]=(b[6]&0x0f)|0x40
 b[8]=(b[8]&0x3f)|0x80
 s:=hex.EncodeToString(b)
 return fmt.Sprintf("%s-%s-%s-%s-%s",s[:8],s[8:12],s[12:16],s[16:20],s[20:])
}

func TestPortfolioPostgresOrganizationAuthorization(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL for migrated test PostgreSQL")}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second)
 defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 orgA,orgB,user,member,unit:=portfolioUUID(t),portfolioUUID(t),portfolioUUID(t),portfolioUUID(t),portfolioUUID(t)
 for _,org:=range []string{orgA,orgB}{
  if _,err:=db.ExecContext(ctx,"INSERT INTO organizations(id,name) VALUES($1,'Portfolio test')",org);err!=nil{t.Fatal(err)}
  t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM organizations WHERE id=$1",org)})
 }
 if _,err:=db.ExecContext(ctx,"INSERT INTO users(id,external_subject) VALUES($1,$2)",user,"portfolio-test-"+user);err!=nil{t.Fatal(err)}
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",user)})
 store:=PortfolioPostgresRepository{DB:db}
 for _,e:=range []PortfolioEntry{
  {OrganizationID:orgA,Domain:"a-"+orgA[:8]+".example.no",Registrar:"external",Status:PortfolioUnknown},
  {OrganizationID:orgB,Domain:"b-"+orgB[:8]+".example.com",Registrar:"external",Status:PortfolioUnknown},
 }{
  if err:=store.Put(ctx,e);err!=nil{t.Fatal(err)}
 }
 reader:=PortfolioReader{DB:db}
 check:=func(org string,want int){
  t.Helper()
  got,err:=reader.ListForUser(ctx,user,org)
  if err!=nil{t.Fatal(err)}
  if len(got)!=want{t.Fatalf("org=%s got=%d want=%d",org,len(got),want)}
  for _,e:=range got {if e.OrganizationID!=org{t.Fatal("cross-organization disclosure")}}
 }
 check(orgA,0);check(orgB,0)
 if _,err:=db.ExecContext(ctx,"INSERT INTO memberships(id,user_id,organization_id,role) VALUES($1,$2,$3,'viewer')",member,user,orgA);err!=nil{t.Fatal(err)}
 check(orgA,1);check(orgB,0)
 for _,role:=range []string{"admin","owner","viewer"}{
  if _,err:=db.ExecContext(ctx,"UPDATE memberships SET role=$2 WHERE id=$1",member,role);err!=nil{t.Fatal(err)}
  check(orgA,1)
 }
 if _,err:=db.ExecContext(ctx,"INSERT INTO units(id,organization_id,name) VALUES($1,$2,'Test unit')",unit,orgA);err!=nil{t.Fatal(err)}
 if _,err:=db.ExecContext(ctx,"UPDATE memberships SET unit_id=$2,role='viewer' WHERE id=$1",member,unit);err!=nil{t.Fatal(err)}
 check(orgA,0)
 if _,err:=db.ExecContext(ctx,"UPDATE memberships SET unit_id=NULL,role='viewer' WHERE id=$1",member);err!=nil{t.Fatal(err)}
 check(orgA,1)
 if _,err:=db.ExecContext(ctx,"DELETE FROM memberships WHERE id=$1",member);err!=nil{t.Fatal(err)}
 check(orgA,0)
}
