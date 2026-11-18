package auth

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

var ErrAuditAccessDenied = errors.New("audit access denied")
var ErrInvalidAuditQuery = errors.New("invalid audit query")

type AuditRow struct {
 ID string
 OrganizationID string
 ActorUserID sql.NullString
 TargetUserID sql.NullString
 Type AuditEventType
 OccurredAt time.Time
}

type AuditPage struct {
 Events []AuditRow
 NextCursor *AuditCursor
}

type AuditCursor struct {
 OccurredAt time.Time
 ID string
}

type AuditQuery struct { DB *sql.DB }

// List returns at most limit events for one organization. The caller must
// derive actorID from authenticated OIDC identity, not request parameters.
// The actor's owner/admin membership is locked until the query finishes.
func (q AuditQuery) List(ctx context.Context,actorID,organizationID string,limit int,cursor *AuditCursor)(AuditPage,error){
 if q.DB==nil||actorID==""||organizationID==""||limit<1||limit>100{
  return AuditPage{},ErrInvalidAuditQuery
 }
 if cursor!=nil&&(cursor.ID==""||cursor.OccurredAt.IsZero()){
  return AuditPage{},ErrInvalidAuditQuery
 }
 tx,err:=q.DB.BeginTx(ctx,&sql.TxOptions{Isolation:sql.LevelReadCommitted})
 if err!=nil{return AuditPage{},fmt.Errorf("begin audit query: %w",err)}
 defer func(){_ = tx.Rollback()}()
 var role string
 err=tx.QueryRowContext(ctx,`
 SELECT role FROM memberships
 WHERE organization_id=$1::uuid AND user_id=$2::uuid
 AND unit_id IS NULL AND role IN ('owner','admin')
 FOR SHARE
 `,organizationID,actorID).Scan(&role)
 if errors.Is(err,sql.ErrNoRows){return AuditPage{},ErrAuditAccessDenied}
 if err!=nil{return AuditPage{},fmt.Errorf("authorize audit query: %w",err)}
 // Use a composite cursor (timestamp, UUID) for stable deterministic order.
 // Fetch one extra row to detect whether a subsequent page exists.
 var rows *sql.Rows
 if cursor==nil{
  rows,err=tx.QueryContext(ctx,`
  SELECT id::text,organization_id::text,actor_user_id::text,target_user_id::text,event_type,occurred_at
  FROM audit_events WHERE organization_id=$1::uuid
  ORDER BY occurred_at DESC,id DESC LIMIT $2
  `,organizationID,limit+1)
 }else{
  rows,err=tx.QueryContext(ctx,`
  SELECT id::text,organization_id::text,actor_user_id::text,target_user_id::text,event_type,occurred_at
  FROM audit_events WHERE organization_id=$1::uuid
  AND (occurred_at,id)<($2,$3::uuid)
  ORDER BY occurred_at DESC,id DESC LIMIT $4
  `,organizationID,cursor.OccurredAt,cursor.ID,limit+1)
 }
 if err!=nil{return AuditPage{},fmt.Errorf("query audit events: %w",err)}
 page:=AuditPage{Events:make([]AuditRow,0,limit)}
 for rows.Next(){
  var event AuditRow
  if err:=rows.Scan(&event.ID,&event.OrganizationID,&event.ActorUserID,&event.TargetUserID,&event.Type,&event.OccurredAt);err!=nil{
   _ = rows.Close()
   return AuditPage{},fmt.Errorf("scan audit event: %w",err)
  }
  page.Events=append(page.Events,event)
 }
 if err:=rows.Err();err!=nil{_ = rows.Close();return AuditPage{},fmt.Errorf("iterate audit events: %w",err)}
 if err:=rows.Close();err!=nil{return AuditPage{},fmt.Errorf("close audit rows: %w",err)}
 if len(page.Events)>limit{
  page.Events=page.Events[:limit]
  last:=page.Events[len(page.Events)-1]
  page.NextCursor=&AuditCursor{OccurredAt:last.OccurredAt,ID:last.ID}
 }
 if err:=tx.Commit();err!=nil{return AuditPage{},fmt.Errorf("finish audit query: %w",err)}
 return page,nil
}
