package domains

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "strings"
)

var ErrInvalidPortfolioPage = errors.New("invalid portfolio page request")

type PortfolioPage struct {
 Entries []PortfolioEntry `json:"domains"`
 NextCursor string `json:"next_cursor,omitempty"`
}

// ListPageForUser enforces membership in the same SQL statement as the read.
// Cursor is the last canonical FQDN, scoped to the supplied organization.
// This primitive requires a trusted internal user ID, never a request-supplied ID.
func (r PortfolioReader) ListPageForUser(ctx context.Context,userID,organizationID string,limit int,cursor string)(PortfolioPage,error){
 if r.DB==nil{return PortfolioPage{},ErrPortfolioStoreUnavailable}
 if ctx==nil||strings.TrimSpace(userID)==""||userID!=strings.TrimSpace(userID)||
  strings.TrimSpace(organizationID)==""||organizationID!=strings.TrimSpace(organizationID)||
  limit<1||limit>100{return PortfolioPage{},ErrInvalidPortfolioPage}
 if cursor!="" {
  name,err:=CanonicalName(cursor)
  if err!=nil||name!=cursor||!strings.Contains(cursor,"."){return PortfolioPage{},ErrInvalidPortfolioPage}
 }
 rows,err:=r.DB.QueryContext(ctx,`SELECT p.organization_id::text,p.fqdn,p.registrar,p.status,p.expires_at,p.last_verified_at
 FROM domain_portfolio p
 WHERE p.organization_id=$2::uuid AND p.fqdn > $3
 AND EXISTS (SELECT 1 FROM memberships m
  WHERE m.organization_id=p.organization_id AND m.user_id=$1::uuid
  AND m.unit_id IS NULL AND m.role IN ('owner','admin','viewer'))
 ORDER BY p.fqdn LIMIT $4`,userID,organizationID,cursor,limit+1)
 if err!=nil{return PortfolioPage{},fmt.Errorf("list portfolio page: %w",err)}
 defer rows.Close()
 entries:=make([]PortfolioEntry,0,limit+1)
 for rows.Next(){
  var e PortfolioEntry
  var status string
  var expires,verified sql.NullTime
  if err:=rows.Scan(&e.OrganizationID,&e.Domain,&e.Registrar,&status,&expires,&verified);err!=nil{return PortfolioPage{},fmt.Errorf("scan portfolio page: %w",err)}
  e.Status=PortfolioStatus(status)
  if expires.Valid{v:=expires.Time;e.ExpiresAt=&v}
  if verified.Valid{v:=verified.Time;e.LastVerifiedAt=&v}
  if err:=ValidatePortfolioEntry(e);err!=nil{return PortfolioPage{},err}
  entries=append(entries,e)
 }
 if err:=rows.Err();err!=nil{return PortfolioPage{},fmt.Errorf("iterate portfolio page: %w",err)}
 page:=PortfolioPage{Entries:entries}
 if len(entries)>limit {
  page.Entries=entries[:limit]
  page.NextCursor=entries[limit-1].Domain
 }
 return page,nil
}
