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

func TestInvitationConsumptionPostgres(t *testing.T) {
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL")}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{t.Fatal(err)}
 defer db.Close()
 db.SetMaxOpenConns(8)
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second)
 defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 now:=time.Now().UTC()
 issuer:="https://test.example.org"
 orgID,err:=randomUUID()
 if err!=nil{t.Fatal(err)}
 if _,err:=db.ExecContext(ctx,"INSERT INTO organizations(id,name) VALUES($1,'Enrollment integration')",orgID);err!=nil{t.Fatal(err)}
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM organizations WHERE id=$1",orgID)})
 p:=UserProvisioner{DB:db}
 makeIdentity:=func(t *testing.T) (Identity,string){
  t.Helper()
  subject,err:=randomUUID()
  if err!=nil{t.Fatal(err)}
  identity:=Identity{Issuer:issuer,Subject:"enrollment-"+subject,Audience:"ploos-domains"}
  userID,err:=p.Provision(ctx,identity)
  if err!=nil{t.Fatal(err)}
  t.Cleanup(func(){
   _,_=db.ExecContext(context.Background(),"DELETE FROM enrollment_invitations WHERE consumed_by=$1",userID)
   _,_=db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",userID)
  })
  return identity,userID
 }
 makeInvite:=func(t *testing.T,identity Identity,expires time.Time) string {
  t.Helper()
  token,digest,err:=NewInvitationToken()
  if err!=nil{t.Fatal(err)}
  invitationID,err:=randomUUID()
  if err!=nil{t.Fatal(err)}
  _,err=db.ExecContext(ctx,`INSERT INTO enrollment_invitations
  (id,token_sha256,organization_id,invited_issuer,invited_subject,created_at,expires_at)
  VALUES($1,$2,$3,$4,$5,$6,$7)`,invitationID,digest[:],orgID,identity.Issuer,identity.Subject,now.Add(-2*time.Hour),expires)
  if err!=nil{t.Fatal(err)}
  return token
 }
 store:=InvitationStore{DB:db}
 t.Run("single use",func(t *testing.T){
  identity,userID:=makeIdentity(t)
  token:=makeInvite(t,identity,now.Add(time.Hour))
  got,err:=store.ConsumeInvitation(ctx,token,identity,userID,now)
  if err!=nil||got!=orgID{t.Fatalf("org=%q err=%v",got,err)}
  _,err=store.ConsumeInvitation(ctx,token,identity,userID,now)
  if !errors.Is(err,ErrInvitationUnavailable){t.Fatalf("replay: %v",err)}
 })
 t.Run("expired",func(t *testing.T){
  identity,userID:=makeIdentity(t)
  token:=makeInvite(t,identity,now.Add(-time.Minute))
  _,err:=store.ConsumeInvitation(ctx,token,identity,userID,now)
  if !errors.Is(err,ErrInvitationUnavailable){t.Fatalf("expired: %v",err)}
 })
 t.Run("wrong identity and user",func(t *testing.T){
  identity,userID:=makeIdentity(t)
  other,otherID:=makeIdentity(t)
  token:=makeInvite(t,identity,now.Add(time.Hour))
  _,err:=store.ConsumeInvitation(ctx,token,other,otherID,now)
  if !errors.Is(err,ErrInvitationUnavailable){t.Fatalf("other identity: %v",err)}
  _,err=store.ConsumeInvitation(ctx,token,identity,otherID,now)
  if !errors.Is(err,ErrInvitationUnavailable){t.Fatalf("other user: %v",err)}
  if _,err:=store.ConsumeInvitation(ctx,token,identity,userID,now);err!=nil{t.Fatal(err)}
 })
 t.Run("tampered token",func(t *testing.T){
  identity,userID:=makeIdentity(t)
  token:=makeInvite(t,identity,now.Add(time.Hour))
  altered:="0"+token[1:]
  if altered==token{altered="1"+token[1:]}
  _,err:=store.ConsumeInvitation(ctx,altered,identity,userID,now)
  if !errors.Is(err,ErrInvitationUnavailable){t.Fatalf("tampered: %v",err)}
 })
 t.Run("concurrent redemption",func(t *testing.T){
  identity,userID:=makeIdentity(t)
  token:=makeInvite(t,identity,now.Add(time.Hour))
  var wg sync.WaitGroup
  results:=make(chan error,6)
  for i:=0;i<6;i++{
   wg.Add(1)
   go func(){defer wg.Done();_,err:=store.ConsumeInvitation(ctx,token,identity,userID,now);results<-err}()
  }
  wg.Wait();close(results)
  successes,denied:=0,0
  for err:=range results {
   if err==nil{successes++}else if errors.Is(err,ErrInvitationUnavailable){denied++}else{t.Fatalf("unexpected error: %v",err)}
  }
  if successes!=1||denied!=5{t.Fatalf("success=%d denied=%d",successes,denied)}
  var consumed int
  digest:=sha256.Sum256([]byte(token))
  if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM enrollment_invitations WHERE token_sha256=$1 AND consumed_at IS NOT NULL",digest[:]).Scan(&consumed);err!=nil{t.Fatal(err)}
  if consumed!=1{t.Fatalf("consumed=%d",consumed)}
  var memberships int
  if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM memberships WHERE user_id=$1",userID).Scan(&memberships);err!=nil{t.Fatal(err)}
  if memberships!=0{t.Fatalf("unexpected memberships=%d",memberships)}
 })
}
