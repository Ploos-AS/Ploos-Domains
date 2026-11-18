package auth

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
)

var (
 ErrMembershipForbidden = errors.New("not authorized to grant membership")
 ErrMembershipExists = errors.New("organization membership already exists")
)

// MembershipManager grants only organization-scoped memberships. Owner role
// creation/transfer and unit-scoped membership are separate privileged flows.
type MembershipManager struct{ DB *sql.DB }

func (m MembershipManager) Grant(ctx context.Context,actorID,organizationID,targetID,role string) error {
 if m.DB==nil || actorID=="" || organizationID=="" || targetID=="" ||
  (role!="admin" && role!="viewer") || actorID==targetID {
  return ErrInvalidIdentity
 }
 tx,err:=m.DB.BeginTx(ctx,nil)
 if err!=nil{return fmt.Errorf("begin membership grant: %w",err)}
 defer func(){_ = tx.Rollback()}()
 var actorRole string
 // Lock the actor's membership while checking permissions, preventing
 // concurrent revocation or demotion before the grant commits.
 err=tx.QueryRowContext(ctx,`
 SELECT role FROM memberships
 WHERE user_id=$1::uuid AND organization_id=$2::uuid
 AND unit_id IS NULL AND role IN ('owner','admin')
 FOR SHARE`,actorID,organizationID).Scan(&actorRole)
 if errors.Is(err,sql.ErrNoRows){return ErrMembershipForbidden}
 if err!=nil{return fmt.Errorf("authorize membership grant: %w",err)}
 if role=="admin" && actorRole!="owner"{return ErrMembershipForbidden}
 id,err:=randomUUID()
 if err!=nil{return fmt.Errorf("generate membership ID: %w",err)}
 // ON CONFLICT matches the partial unique index and never modifies an
 // existing membership (including a privileged one).
 var created string
 err=tx.QueryRowContext(ctx,`
 INSERT INTO memberships(id,user_id,organization_id,unit_id,role)
 VALUES($1,$2::uuid,$3::uuid,NULL,$4)
 ON CONFLICT (organization_id,user_id) WHERE unit_id IS NULL DO NOTHING
 RETURNING id::text`,id,targetID,organizationID,role).Scan(&created)
 if errors.Is(err,sql.ErrNoRows){return ErrMembershipExists}
 if err!=nil{return fmt.Errorf("insert membership: %w",err)}
 if err:=tx.Commit();err!=nil{return fmt.Errorf("commit membership grant: %w",err)}
 return nil
}
