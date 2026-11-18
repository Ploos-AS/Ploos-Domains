package verification

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
)

// PostgresAuthorizer allows only organization-scoped owners and admins
// to verify a domain belonging to their organization.
// actorID must be a trusted, authenticated internal user UUID.
type PostgresAuthorizer struct { DB *sql.DB }

func (a PostgresAuthorizer) CanVerifyDomain(ctx context.Context, actorID, organizationID, domainID string) (bool,error) {
 if a.DB==nil || actorID=="" || organizationID=="" || domainID=="" {
  return false,ErrInvalid
 }
 var allowed bool
 err:=a.DB.QueryRowContext(ctx,`
 SELECT EXISTS (
  SELECT 1 FROM domains d
  JOIN memberships m ON m.organization_id=d.organization_id
  WHERE d.id=$1::uuid AND d.organization_id=$2::uuid
    AND m.user_id=$3::uuid AND m.unit_id IS NULL
    AND m.role IN ('owner','admin')
 )`,domainID,organizationID,actorID).Scan(&allowed)
 if errors.Is(err,sql.ErrNoRows) { return false,nil }
 if err!=nil { return false,fmt.Errorf("authorize domain verification: %w",err) }
 return allowed,nil
}

var _ DomainAuthorizer = PostgresAuthorizer{}
