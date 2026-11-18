//go:build integration

package verification

import (
 "context"
 "database/sql"
 "os"
 "testing"
 "time"
)

func TestPostgresAuthorizer(t *testing.T) {
 db:=openIntegrationDB(t)
 ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second)
 defer cancel()
 domain,_,_:=setupVerification(t,db,time.Now().Add(time.Hour))
 otherDomain,_,_:=setupVerification(t,db,time.Now().Add(time.Hour))
 actor:=newUUID(t)
 _,err:=db.ExecContext(ctx,"INSERT INTO users(id,external_subject) VALUES($1,$2)",actor,"integration-"+actor)
 if err!=nil {t.Fatal(err)}
 t.Cleanup(func(){_,_=db.ExecContext(context.Background(),"DELETE FROM users WHERE id=$1",actor)})
 var orgID,otherOrgID string
 if err:=db.QueryRowContext(ctx,"SELECT organization_id::text FROM domains WHERE id=$1",domain).Scan(&orgID);err!=nil{t.Fatal(err)}
 if err:=db.QueryRowContext(ctx,"SELECT organization_id::text FROM domains WHERE id=$1",otherDomain).Scan(&otherOrgID);err!=nil{t.Fatal(err)}
 auth:=PostgresAuthorizer{DB:db}
 check:=func(o,d string,want bool){
  t.Helper()
  got,err:=auth.CanVerifyDomain(ctx,actor,o,d)
  if err!=nil||got!=want {t.Fatalf("org=%s domain=%s got=%v err=%v want=%v",o,d,got,err,want)}
 }
 check(orgID,domain,false)
 membership:=newUUID(t)
 _,err=db.ExecContext(ctx,`INSERT INTO memberships(id,user_id,organization_id,role) VALUES($1,$2,$3,'viewer')`,membership,actor,orgID)
 if err!=nil{t.Fatal(err)}
 check(orgID,domain,false)
 unit:=newUUID(t)
 _,err=db.ExecContext(ctx,"INSERT INTO units(id,organization_id,name) VALUES($1,$2,'Test unit')",unit,orgID)
 if err!=nil{t.Fatal(err)}
 _,err=db.ExecContext(ctx,"UPDATE memberships SET unit_id=$2,role='unit_admin' WHERE id=$1",membership,unit)
 if err!=nil{t.Fatal(err)}
 check(orgID,domain,false)
 _,err=db.ExecContext(ctx,"UPDATE memberships SET unit_id=NULL,role='admin' WHERE id=$1",membership)
 if err!=nil{t.Fatal(err)}
 check(orgID,domain,true)
 check(otherOrgID,otherDomain,false)
 check(otherOrgID,domain,false)
 check(orgID,otherDomain,false)
 _,err=db.ExecContext(ctx,"UPDATE memberships SET role='owner' WHERE id=$1",membership)
 if err!=nil{t.Fatal(err)}
 check(orgID,domain,true)
 _,err=db.ExecContext(ctx,"UPDATE memberships SET role='viewer' WHERE id=$1",membership)
 if err!=nil{t.Fatal(err)}
 check(orgID,domain,false)
}

func openIntegrationDB(t *testing.T)*sql.DB{
 t.Helper()
 dsn:=os.Getenv("PLOOS_TEST_DATABASE_URL")
 if dsn==""{t.Skip("set PLOOS_TEST_DATABASE_URL")}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{t.Fatal(err)}
 t.Cleanup(func(){_ = db.Close()})
 if err:=db.Ping();err!=nil{t.Fatal(err)}
 return db
}
