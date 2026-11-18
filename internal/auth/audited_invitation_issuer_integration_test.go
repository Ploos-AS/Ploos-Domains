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

func TestAuditedInvitationIssuancePostgres(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL")}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{t.Fatal(err)}
 defer db.Close()
 ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second)
 defer cancel()
 org,_:=randomUUID()
 owner,_:=randomUUID()
 viewer,_:=randomUUID()
 for _,uid:=range []string{owner,viewer}{
  if _,err:=db.ExecContext(ctx,"INSERT INTO users(id,external_subject) VALUES($1,$2)",uid,"audited-issuance:"+uid);err!=nil{t.Fatal(err)}
  defer db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",uid)
 }
 if _,err:=db.ExecContext(ctx,"INSERT INTO organizations(id,name) VALUES($1,'Audited issuance')",org);err!=nil{t.Fatal(err)}
 defer func(){
  _,_=db.ExecContext(context.Background(),"DELETE FROM audit_events WHERE organization_id=$1",org)
  _,_=db.ExecContext(context.Background(),"DELETE FROM organizations WHERE id=$1",org)
 }()
 for _,item:=range []struct{id,role string}{{owner,"owner"},{viewer,"viewer"}}{
  mid,_:=randomUUID()
  if _,err:=db.ExecContext(ctx,"INSERT INTO memberships(id,user_id,organization_id,role) VALUES($1,$2,$3,$4)",mid,item.id,org,item.role);err!=nil{t.Fatal(err)}
 }
 now:=time.Now().UTC()
 s:=AuditedInvitationIssuer{DB:db,Clock:func()time.Time{return now}}
 _,err=s.Issue(ctx,viewer,org,"https://id.example.org","subject",time.Hour)
 if !errors.Is(err,ErrInvitationForbidden){t.Fatalf("viewer issue: %v",err)}
 token,err:=s.Issue(ctx,owner,org,"https://id.example.org","subject",time.Hour)
 if err!=nil{t.Fatal(err)}
 digest:=sha256.Sum256([]byte(token))
 var invitationCount,auditCount int
 if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM enrollment_invitations WHERE organization_id=$1 AND token_sha256=$2",org,digest[:]).Scan(&invitationCount);err!=nil{t.Fatal(err)}
 if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM audit_events WHERE organization_id=$1 AND actor_user_id=$2 AND event_type='invitation.issued'",org,owner).Scan(&auditCount);err!=nil{t.Fatal(err)}
 if invitationCount!=1||auditCount!=1{t.Fatalf("invitations=%d audit=%d",invitationCount,auditCount)}
}
