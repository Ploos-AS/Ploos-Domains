//go:build integration

package auth

import (
 "context"
 "database/sql"
 "errors"
 "os"
 "testing"
 "time"

 _ "github.com/lib/pq"
)

func TestOrganizationAuditQueryPostgres(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL")}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second)
 defer cancel()
 org1,_:=randomUUID()
 org2,_:=randomUUID()
 owner,_:=randomUUID()
 viewer,_:=randomUUID()
 for _,uid:=range []string{owner,viewer}{
  if _,err:=db.ExecContext(ctx,"INSERT INTO users(id,external_subject) VALUES($1,$2)",uid,"audit-query:"+uid);err!=nil{t.Fatal(err)}
  defer db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",uid)
 }
 for _,org:=range []string{org1,org2}{
  if _,err:=db.ExecContext(ctx,"INSERT INTO organizations(id,name) VALUES($1,'Audit query test')",org);err!=nil{t.Fatal(err)}
  defer func(id string){
   _,_=db.ExecContext(context.Background(),"DELETE FROM audit_events WHERE organization_id=$1",id)
   _,_=db.ExecContext(context.Background(),"DELETE FROM organizations WHERE id=$1",id)
  }(org)
 }
 for _,tc:=range []struct{uid,org,role string}{{owner,org1,"owner"},{viewer,org1,"viewer"}}{
  id,_:=randomUUID()
  if _,err:=db.ExecContext(ctx,"INSERT INTO memberships(id,user_id,organization_id,role) VALUES($1,$2,$3,$4)",id,tc.uid,tc.org,tc.role);err!=nil{t.Fatal(err)}
 }
 now:=time.Now().UTC()
 for i:=0;i<3;i++{
  id,_:=randomUUID()
  if _,err:=db.ExecContext(ctx,"INSERT INTO audit_events(id,organization_id,event_type,occurred_at) VALUES($1,$2,'membership.granted',$3)",id,org1,now.Add(time.Duration(i)*time.Second));err!=nil{t.Fatal(err)}
 }
 id,_:=randomUUID()
 if _,err:=db.ExecContext(ctx,"INSERT INTO audit_events(id,organization_id,event_type,occurred_at) VALUES($1,$2,'membership.granted',$3)",id,org2,now);err!=nil{t.Fatal(err)}
 q:=AuditQuery{DB:db}
 _,err=q.List(ctx,viewer,org1,2,nil)
 if !errors.Is(err,ErrAuditAccessDenied){t.Fatalf("viewer err=%v",err)}
 _,err=q.List(ctx,owner,org2,2,nil)
 if !errors.Is(err,ErrAuditAccessDenied){t.Fatalf("other org err=%v",err)}
 first,err:=q.List(ctx,owner,org1,2,nil)
 if err!=nil{t.Fatal(err)}
 if len(first.Events)!=2||first.NextCursor==nil{t.Fatalf("first page=%+v",first)}
 second,err:=q.List(ctx,owner,org1,2,first.NextCursor)
 if err!=nil{t.Fatal(err)}
 if len(second.Events)!=1||second.NextCursor!=nil{t.Fatalf("second page=%+v",second)}
 for _,event:=range append(first.Events,second.Events...){
  if event.OrganizationID!=org1{t.Fatal("cross-organization audit leak")}
 }
}
