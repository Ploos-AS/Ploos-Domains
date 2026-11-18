//go:build integration

package auth

import (
 "context"
 "database/sql"
 "crypto/sha256"
 "errors"
 "os"
 "sync"
 "testing"
 "time"

 _ "github.com/lib/pq"
)

func TestAtomicInvitationMembershipPostgres(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL")}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{t.Fatal(err)}
 defer db.Close()
 db.SetMaxOpenConns(8)
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second)
 defer cancel()
 org,_:=randomUUID()
 if _,err:=db.ExecContext(ctx,"INSERT INTO organizations(id,name) VALUES($1,'Atomic redemption test')",org);err!=nil{t.Fatal(err)}
 defer func(){
  _,_=db.ExecContext(context.Background(),"DELETE FROM audit_events WHERE organization_id=$1",org)
  _,_=db.ExecContext(context.Background(),"DELETE FROM organizations WHERE id=$1",org)
 }()
 now:=time.Now().UTC()
 provisioner:=UserProvisioner{DB:db}
 setup:=func(t *testing.T)(Identity,string,string,[32]byte){
  t.Helper()
  subject,_:=randomUUID()
  identity:=Identity{Issuer:"https://test.example.org",Subject:"atomic-"+subject,Audience:"ploos-domains"}
  uid,err:=provisioner.Provision(ctx,identity)
  if err!=nil{t.Fatal(err)}
  t.Cleanup(func(){
   _,_=db.ExecContext(context.Background(),"DELETE FROM enrollment_invitations WHERE consumed_by=$1",uid)
   _,_=db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",uid)
  })
  token,digest,err:=NewInvitationToken()
  if err!=nil{t.Fatal(err)}
  id,_:=randomUUID()
  _,err=db.ExecContext(ctx,`INSERT INTO enrollment_invitations
  (id,token_sha256,organization_id,invited_issuer,invited_subject,expires_at)
  VALUES($1,$2,$3,$4,$5,$6)`,id,digest[:],org,identity.Issuer,identity.Subject,now.Add(time.Hour))
  if err!=nil{t.Fatal(err)}
  return identity,uid,token,digest
 }
 service:=InvitationMembership{DB:db}
 t.Run("atomic success and replay",func(t *testing.T){
  identity,uid,token,_:=setup(t)
  got,err:=service.Redeem(ctx,token,identity,uid,now)
  if err!=nil||got!=org{t.Fatalf("org=%s err=%v",got,err)}
  var role string
  if err:=db.QueryRowContext(ctx,"SELECT role FROM memberships WHERE organization_id=$1 AND user_id=$2 AND unit_id IS NULL",org,uid).Scan(&role);err!=nil{t.Fatal(err)}
  if role!="viewer"{t.Fatalf("role=%s",role)}
  var events int
  if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM audit_events WHERE organization_id=$1 AND actor_user_id=$2 AND event_type=$3",org,uid,string(AuditInvitationRedeemed)).Scan(&events);err!=nil{t.Fatal(err)}
  if events!=1{t.Fatalf("redemption audit events=%d",events)}
  _,err=service.Redeem(ctx,token,identity,uid,now)
  if !errors.Is(err,ErrInvitationUnavailable){t.Fatalf("replay=%v",err)}
 })
 t.Run("existing membership does not consume",func(t *testing.T){
  identity,uid,token,digest:=setup(t)
  id,_:=randomUUID()
  if _,err:=db.ExecContext(ctx,"INSERT INTO memberships(id,user_id,organization_id,role) VALUES($1,$2,$3,'admin')",id,uid,org);err!=nil{t.Fatal(err)}
  _,err:=service.Redeem(ctx,token,identity,uid,now)
  if !errors.Is(err,ErrMembershipExists){t.Fatalf("duplicate=%v",err)}
  var consumed bool
  if err:=db.QueryRowContext(ctx,"SELECT consumed_at IS NOT NULL FROM enrollment_invitations WHERE token_sha256=$1",digest[:]).Scan(&consumed);err!=nil{t.Fatal(err)}
  if consumed{t.Fatal("failed grant consumed invitation")}
  var role string
  if err:=db.QueryRowContext(ctx,"SELECT role FROM memberships WHERE id=$1",id).Scan(&role);err!=nil{t.Fatal(err)}
  if role!="admin"{t.Fatalf("role overwritten: %s",role)}
 })
 t.Run("wrong identity",func(t *testing.T){
  identity,uid,token,_:=setup(t)
  identity.Subject="not-the-invited-subject"
  _,err:=service.Redeem(ctx,token,identity,uid,now)
  if !errors.Is(err,ErrInvitationUnavailable){t.Fatalf("wrong identity=%v",err)}
 })
 t.Run("concurrent",func(t *testing.T){
  identity,uid,token,_:=setup(t)
  var wg sync.WaitGroup
  ch:=make(chan error,6)
  for i:=0;i<6;i++{
   wg.Add(1)
   go func(){defer wg.Done();_,err:=service.Redeem(ctx,token,identity,uid,now);ch<-err}()
  }
  wg.Wait();close(ch)
  successes,denials:=0,0
  for err:=range ch{
   if err==nil{successes++}else if errors.Is(err,ErrInvitationUnavailable){denials++}else{t.Fatalf("concurrent error=%v",err)}
  }
  if successes!=1||denials!=5{t.Fatalf("successes=%d denied=%d",successes,denials)}
  digest:=sha256.Sum256([]byte(token))
  var count int
  if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM enrollment_invitations WHERE token_sha256=$1 AND consumed_at IS NOT NULL",digest[:]).Scan(&count);err!=nil{t.Fatal(err)}
  if count!=1{t.Fatalf("consumed=%d",count)}
 })
}
