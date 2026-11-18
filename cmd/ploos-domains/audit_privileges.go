package main

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
)

var errAuditPrivileges = errors.New("audit database privileges are unsafe")

// checkAuditPrivileges runs as the configured application login. A schema
// owner, superuser, or role with audit mutation rights must not serve audit reads.
// PostgreSQL evaluates inherited role grants for has_table_privilege.
func checkAuditPrivileges(ctx context.Context,db *sql.DB)error{
 if db==nil{return errAuditPrivileges}
 var superuser,canRead,canUpdate,canDelete,canTruncate,ownsAudit bool
 err:=db.QueryRowContext(ctx,`
 SELECT r.rolsuper,
  has_table_privilege(current_user,'public.audit_events','SELECT'),
  has_table_privilege(current_user,'public.audit_events','UPDATE'),
  has_table_privilege(current_user,'public.audit_events','DELETE'),
  has_table_privilege(current_user,'public.audit_events','TRUNCATE'),
  pg_has_role(current_user,c.relowner,'MEMBER')
 FROM pg_roles r
 JOIN pg_class c ON c.oid='public.audit_events'::regclass
 WHERE r.rolname=current_user
 `).Scan(&superuser,&canRead,&canUpdate,&canDelete,&canTruncate,&ownsAudit)
 if err!=nil{return fmt.Errorf("check audit privileges: %w",err)}
 if superuser||!canRead||canUpdate||canDelete||canTruncate||ownsAudit{
  return errAuditPrivileges
 }
 for _,table:=range []string{"public.memberships","public.user_oidc_identities"}{
  var allowed bool
  if err:=db.QueryRowContext(ctx,"SELECT has_table_privilege(current_user,$1,'SELECT')",table).Scan(&allowed);err!=nil{
   return fmt.Errorf("check required read privilege: %w",err)
  }
  if !allowed{return errAuditPrivileges}
 }
 return nil
}
