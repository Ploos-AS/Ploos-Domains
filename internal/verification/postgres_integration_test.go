//go:build integration

package verification

import (
 "context"
 "crypto/rand"
 "crypto/sha256"
 "encoding/hex"
 "errors"
 "fmt"
 "os"
 "sync"
 "testing"
 "time"

 "database/sql"
 _ "github.com/lib/pq"
)

func newUUID(t *testing.T) string {
 t.Helper()
 b:=make([]byte,16)
 if _,err:=rand.Read(b);err!=nil{t.Fatal(err)}
 b[6]=(b[6]&0x0f)|0x40
 b[8]=(b[8]&0x3f)|0x80
 s:=hex.EncodeToString(b)
 return fmt.Sprintf("%s-%s-%s-%s-%s",s[:8],s[8:12],s[12:16],s[16:20],s[20:])
}

func setupVerification(t *testing.T,db *sql.DB,expiresAt time.Time)(string,string,string) {
 t.Helper()
 org,domain,challenge:=newUUID(t),newUUID(t),newUUID(t)
 fqdn:="test-"+org[:8]+".example.org"
 token:="ploos-domain-verification="+challenge
 digest:=sha256.Sum256([]byte(token))
 ctx:=context.Background()
 if _,err:=db.ExecContext(ctx,"INSERT INTO organizations(id,name) VALUES($1,'Integration Test')",org);err!=nil{t.Fatal(err)}
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM organizations WHERE id=$1",org)})
 if _,err:=db.ExecContext(ctx,"INSERT INTO domains(id,organization_id,fqdn) VALUES($1,$2,$3)",domain,org,fqdn);err!=nil{t.Fatal(err)}
 if _,err:=db.ExecContext(ctx,`INSERT INTO domain_verification_challenges
 (id,domain_id,token_sha256,record_name,created_at,expires_at)
 VALUES($1,$2,$3,$4,$5,$6)`,challenge,domain,digest[:],"_ploos-verify."+fqdn,expiresAt.Add(-2*time.Hour),expiresAt);err!=nil{t.Fatal(err)}
 return domain,challenge,token
}

func TestPostgresVerificationLifecycle(t *testing.T){
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL for PostgreSQL integration tests")}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{t.Fatal(err)}
 defer db.Close()
 db.SetMaxOpenConns(8)
 ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second)
 defer cancel()
 if err:=db.PingContext(ctx);err!=nil{t.Fatal(err)}
 now:=time.Now().UTC()
 t.Run("success and replay",func(t *testing.T){
  domain,challenge,_:=setupVerification(t,db,now.Add(time.Hour))
  if err:=CompleteVerification(ctx,db,domain,challenge,now);err!=nil{t.Fatal(err)}
  if err:=CompleteVerification(ctx,db,domain,challenge,now);!errors.Is(err,ErrDomainAlreadyVerified){t.Fatalf("replay got %v",err)}
  var verified,consumed bool
  if err:=db.QueryRowContext(ctx,`SELECT d.verified_at IS NOT NULL,c.consumed_at IS NOT NULL
  FROM domains d JOIN domain_verification_challenges c ON c.domain_id=d.id WHERE d.id=$1`,domain).Scan(&verified,&consumed);err!=nil{t.Fatal(err)}
  if !verified||!consumed{t.Fatal("verification and consumption not atomic")}
 })
 t.Run("expired",func(t *testing.T){
  domain,challenge,_:=setupVerification(t,db,now.Add(-time.Minute))
  if err:=CompleteVerification(ctx,db,domain,challenge,now);!errors.Is(err,ErrChallengeUnavailable){t.Fatalf("got %v",err)}
 })
 t.Run("wrong challenge",func(t *testing.T){
  domain,_,_:=setupVerification(t,db,now.Add(time.Hour))
  if err:=CompleteVerification(ctx,db,domain,newUUID(t),now);!errors.Is(err,ErrChallengeUnavailable){t.Fatalf("got %v",err)}
 })
 t.Run("concurrent consumption",func(t *testing.T){
  domain,challenge,_:=setupVerification(t,db,now.Add(time.Hour))
  var wg sync.WaitGroup
  errs:=make(chan error,2)
  for i:=0;i<2;i++ {wg.Add(1);go func(){defer wg.Done();errs<-CompleteVerification(context.Background(),db,domain,challenge,now)}()}
  wg.Wait();close(errs)
  success,rejected:=0,0
  for err:=range errs {
   if err==nil {success++} else if errors.Is(err,ErrDomainAlreadyVerified)||errors.Is(err,ErrChallengeUnavailable){rejected++} else {t.Fatalf("unexpected concurrent error: %v",err)}
  }
  if success!=1||rejected!=1{t.Fatalf("success=%d rejected=%d",success,rejected)}
 })
}
