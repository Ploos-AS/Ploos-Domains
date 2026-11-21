package domains

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "strings"
)

var ErrPortfolioAccessDenied = errors.New("portfolio access denied")

// PortfolioReader exposes only organization-authorized reads. Both membership
// and portfolio filters are evaluated by PostgreSQL in the same statement.
// Never expose PortfolioPostgresRepository.List directly to customer callers.
type PortfolioReader struct { DB *sql.DB }

func (r PortfolioReader) ListForUser(ctx context.Context,userID,organizationID string)([]PortfolioEntry,error){
 if r.DB==nil {return nil,ErrPortfolioStoreUnavailable}
 if ctx==nil || strings.TrimSpace(userID)=="" || userID!=strings.TrimSpace(userID) ||
  strings.TrimSpace(organizationID)=="" || organizationID!=strings.TrimSpace(organizationID) {return nil,ErrInvalidPortfolioEntry}
 // An organization-level membership is required. Unit-only memberships do
 // not grant access to the entire organization's domain portfolio.
 rows,err:=r.DB.QueryContext(ctx,`SELECT p.organization_id::text,p.fqdn,p.registrar,p.status,p.expires_at,p.last_verified_at
 FROM domain_portfolio p
 WHERE p.organization_id=$2::uuid
 AND EXISTS (SELECT 1 FROM memberships m
             WHERE m.organization_id=p.organization_id
             AND m.user_id=$1::uuid AND m.unit_id IS NULL
             AND m.role IN ('owner','admin','viewer'))
 ORDER BY p.fqdn`,userID,organizationID)
 if err!=nil {return nil,fmt.Errorf("authorized portfolio read: %w",err)}
 defer rows.Close()
 entries:=make([]PortfolioEntry,0)
 for rows.Next(){
  var e PortfolioEntry
  var status string
  var expires,verified sql.NullTime
  if err:=rows.Scan(&e.OrganizationID,&e.Domain,&e.Registrar,&status,&expires,&verified);err!=nil{return nil,fmt.Errorf("scan authorized portfolio: %w",err)}
  e.Status=PortfolioStatus(status)
  if expires.Valid {v:=expires.Time;e.ExpiresAt=&v}
  if verified.Valid {v:=verified.Time;e.LastVerifiedAt=&v}
  if err:=ValidatePortfolioEntry(e);err!=nil{return nil,err}
  entries=append(entries,e)
 }
 if err:=rows.Err();err!=nil{return nil,fmt.Errorf("iterate authorized portfolio: %w",err)}
 return entries,nil
}
