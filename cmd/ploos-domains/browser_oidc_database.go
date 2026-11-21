package main

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
)

var errBrowserOIDCDatabase=errors.New("browser OIDC database not ready")

// checkBrowserOIDCDatabase checks required migrated tables and minimal
// privileges for the future browser login runtime role. It does not mutate
// data or grant permissions. Use a dedicated least-privilege database role.
func checkBrowserOIDCDatabase(ctx context.Context,db *sql.DB)error{
 if ctx==nil||db==nil{return errBrowserOIDCDatabase}
 var superuser bool
 if err:=db.QueryRowContext(ctx,"SELECT rolsuper FROM pg_roles WHERE rolname=current_user").Scan(&superuser);err!=nil{
  return fmt.Errorf("%w: role lookup: %v",errBrowserOIDCDatabase,err)
 }
 if superuser{return errBrowserOIDCDatabase}
 for _,item:=range []struct{table string;privileges []string}{
  {"public.users",[]string{"SELECT"}},
  {"public.user_oidc_identities",[]string{"SELECT"}},
  {"public.oidc_login_transactions",[]string{"SELECT","INSERT","UPDATE","DELETE"}},
  {"public.browser_sessions",[]string{"SELECT","INSERT","UPDATE","DELETE"}},
 }{
  var owns bool
  if err:=db.QueryRowContext(ctx,`SELECT pg_has_role(current_user,c.relowner,'MEMBER') FROM pg_class c WHERE c.oid=$1::regclass`,item.table).Scan(&owns);err!=nil{
   return fmt.Errorf("%w: missing table %s: %v",errBrowserOIDCDatabase,item.table,err)
  }
  if owns{return fmt.Errorf("%w: application role owns %s",errBrowserOIDCDatabase,item.table)}
  for _,privilege:=range item.privileges{
   var allowed bool
   if err:=db.QueryRowContext(ctx,"SELECT has_table_privilege(current_user,$1,$2)",item.table,privilege).Scan(&allowed);err!=nil{
    return fmt.Errorf("%w: privilege check %s %s: %v",errBrowserOIDCDatabase,item.table,privilege,err)
   }
   if !allowed{return fmt.Errorf("%w: missing %s on %s",errBrowserOIDCDatabase,privilege,item.table)}
  }
 }
 return nil
}
