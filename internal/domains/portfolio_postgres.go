package domains

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "strings"
 "time"
)

var ErrPortfolioStoreUnavailable = errors.New("domain portfolio store unavailable")

// PortfolioPostgresRepository is a persistence primitive. Caller must
// authorize access to organizationID before calling any method.
type PortfolioPostgresRepository struct { DB *sql.DB }

// List returns only rows belonging to the supplied organization. No registrar
// calls or DNS ownership checks occur here.
func (r PortfolioPostgresRepository) List(ctx context.Context,organizationID string)([]PortfolioEntry,error){
 if r.DB==nil {return nil,ErrPortfolioStoreUnavailable}
 if ctx==nil || strings.TrimSpace(organizationID)=="" || organizationID!=strings.TrimSpace(organizationID) {return nil,ErrInvalidPortfolioEntry}
 rows,err:=r.DB.QueryContext(ctx,`SELECT organization_id::text,fqdn,registrar,status,expires_at,last_verified_at
 FROM domain_portfolio WHERE organization_id=$1::uuid ORDER BY fqdn`,organizationID)
 if err!=nil {return nil,fmt.Errorf("list portfolio: %w",err)}
 defer rows.Close()
 entries:=make([]PortfolioEntry,0)
 for rows.Next(){
  var e PortfolioEntry
  var status string
  var expires,verified sql.NullTime
  if err:=rows.Scan(&e.OrganizationID,&e.Domain,&e.Registrar,&status,&expires,&verified);err!=nil{return nil,fmt.Errorf("scan portfolio: %w",err)}
  e.Status=PortfolioStatus(status)
  if expires.Valid {v:=expires.Time;e.ExpiresAt=&v}
  if verified.Valid {v:=verified.Time;e.LastVerifiedAt=&v}
  if err:=ValidatePortfolioEntry(e);err!=nil{return nil,err}
  entries=append(entries,e)
 }
 if err:=rows.Err();err!=nil{return nil,fmt.Errorf("iterate portfolio: %w",err)}
 return entries,nil
}

// Put records locally known metadata only. A caller must authorize the
// organization and verify registrar-derived values before marking active.
func (r PortfolioPostgresRepository) Put(ctx context.Context,e PortfolioEntry)error{
 if r.DB==nil {return ErrPortfolioStoreUnavailable}
 if ctx==nil {return ErrInvalidPortfolioEntry}
 if err:=ValidatePortfolioEntry(e);err!=nil{return err}
 _,err:=r.DB.ExecContext(ctx,`INSERT INTO domain_portfolio
 (organization_id,fqdn,registrar,status,expires_at,last_verified_at)
 VALUES ($1::uuid,$2,$3,$4,$5,$6)
 ON CONFLICT (organization_id,fqdn) DO UPDATE SET
 registrar=EXCLUDED.registrar,status=EXCLUDED.status,
 expires_at=EXCLUDED.expires_at,last_verified_at=EXCLUDED.last_verified_at,
 updated_at=now()`,e.OrganizationID,e.Domain,e.Registrar,string(e.Status),nullableTime(e.ExpiresAt),nullableTime(e.LastVerifiedAt))
 if err!=nil{return fmt.Errorf("store portfolio: %w",err)}
 return nil
}
func nullableTime(t *time.Time)any{
 if t==nil{return nil}
 return *t
}
