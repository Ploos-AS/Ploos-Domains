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

func TestMembershipGrantsPostgres(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL")}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second)
 defer cancel()
 org,_:=randomUUID()
 owner,_:=randomUUID()
 admin,_:=randomUUID()
 viewer,_:=randomUUID()
 target,_:=randomUUID()
 for _,id:=range []string{owner,admin,viewer,target}{
  if _,err:=db.ExecContext(ctx,"INSERT INTO users(id,external_subject) VALUES($1,$2)",id,"membership-test:"+id);err!=nil{t.Fatal(err)}
  defer db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",id)
 }
 if _,err:=db.ExecContext(ctx,"INSERT INTO organizations(id,name) VALUES($1,'Membership integration')",org);err!=nil{t.Fatal(err)}
 defer db.ExecContext(context.Background(),"DELETE FROM organizations WHERE id=$1",org)
 for _,tc:=range []struct{id,role string}{{owner,"owner"},{admin,"admin"},{viewer,"viewer"}}{
  membershipID,_:=randomUUID()
  if _,err:=db.ExecContext(ctx,"INSERT INTO memberships(id,user_id,organization_id,role) VALUES($1,$2,$3,$4)",membershipID,tc.id,org,tc.role);err!=nil{t.Fatal(err)}
 }
 m:=MembershipManager{DB:db}
 if err:=m.Grant(ctx,admin,org,target,"admin");!errors.Is(err,ErrMembershipForbidden){t.Fatalf("admin escalation: %v",err)}
 if err:=m.Grant(ctx,viewer,org,target,"viewer");!errors.Is(err,ErrMembershipForbidden){t.Fatalf("viewer grant: %v",err)}
 if err:=m.Grant(ctx,owner,org,target,"owner");!errors.Is(err,ErrInvalidIdentity){t.Fatalf("owner grant: %v",err)}
 if err:=m.Grant(ctx,owner,org,target,"admin");err!=nil{t.Fatal(err)}
 if err:=m.Grant(ctx,owner,org,target,"viewer");!errors.Is(err,ErrMembershipExists){t.Fatalf("duplicate: %v",err)}
 var role string
 if err:=db.QueryRowContext(ctx,"SELECT role FROM memberships WHERE organization_id=$1 AND user_id=$2 AND unit_id IS NULL",org,target).Scan(&role);err!=nil{t.Fatal(err)}
 if role!="admin"{t.Fatalf("existing role overwritten: %s",role)}
}
