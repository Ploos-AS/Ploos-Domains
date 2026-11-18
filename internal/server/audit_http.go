package server

import (
 "context"
 "encoding/base64"
 "encoding/json"
 "errors"
 "net/http"
 "strconv"
 "strings"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/auth"
)

type AuditIdentityResolver interface {
 ResolveUserID(context.Context,auth.Identity)(string,error)
}
type AuditLister interface {
 List(context.Context,string,string,int,*auth.AuditCursor)(auth.AuditPage,error)
}

type auditCursorWire struct {
 At time.Time `json:"at"`
 ID string `json:"id"`
}
type auditEventJSON struct {
 ID string `json:"id"`
 OrganizationID string `json:"organization_id"`
 ActorUserID *string `json:"actor_user_id"`
 TargetUserID *string `json:"target_user_id"`
 Type auth.AuditEventType `json:"type"`
 OccurredAt time.Time `json:"occurred_at"`
}
type auditPageJSON struct {
 Events []auditEventJSON `json:"events"`
 NextCursor string `json:"next_cursor,omitempty"`
}

func decodeAuditCursor(value string)(*auth.AuditCursor,error){
 if value==""{return nil,nil}
 if len(value)>512{return nil,auth.ErrInvalidAuditQuery}
 raw,err:=base64.RawURLEncoding.DecodeString(value)
 if err!=nil{return nil,auth.ErrInvalidAuditQuery}
 var wire auditCursorWire
 if err:=json.Unmarshal(raw,&wire);err!=nil||wire.ID==""||wire.At.IsZero(){return nil,auth.ErrInvalidAuditQuery}
 return &auth.AuditCursor{OccurredAt:wire.At,ID:wire.ID},nil
}
func encodeAuditCursor(c *auth.AuditCursor)string{
 if c==nil{return ""}
 raw,_:=json.Marshal(auditCursorWire{At:c.OccurredAt,ID:c.ID})
 return base64.RawURLEncoding.EncodeToString(raw)
}
func auditNullable(v interface{ }) *string{
 // This helper is intentionally not used for identity input.
 return nil
}

// NewAuditHandler returns an authenticated read-only handler.
// Mount under GET /api/v1/organizations/{organizationID}/audit-events
// and wrap it with auth.Middleware configured with a trusted verifier.
func NewAuditHandler(resolver AuditIdentityResolver,lister AuditLister)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Method!=http.MethodGet{w.Header().Set("Allow","GET");http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
  if resolver==nil||lister==nil{http.Error(w,"audit unavailable",http.StatusServiceUnavailable);return}
  identity,ok:=auth.FromContext(r.Context())
  if !ok||identity.Validate()!=nil{http.Error(w,"unauthorized",http.StatusUnauthorized);return}
  organizationID:=r.PathValue("organizationID")
  if organizationID==""{http.Error(w,"invalid organization",http.StatusBadRequest);return}
  values:=r.URL.Query()
  for key:=range values{
   if key!="limit"&&key!="cursor"{http.Error(w,"invalid query",http.StatusBadRequest);return}
   if len(values[key])!=1{http.Error(w,"invalid query",http.StatusBadRequest);return}
  }
  limit:=50
  if raw:=values.Get("limit");raw!=""{
   n,err:=strconv.Atoi(raw)
   if err!=nil||n<1||n>100{http.Error(w,"invalid limit",http.StatusBadRequest);return}
   limit=n
  }
  cursor,err:=decodeAuditCursor(values.Get("cursor"))
  if err!=nil{http.Error(w,"invalid cursor",http.StatusBadRequest);return}
  actorID,err:=resolver.ResolveUserID(r.Context(),identity)
  if errors.Is(err,auth.ErrIdentityNotLinked){http.Error(w,"forbidden",http.StatusForbidden);return}
  if err!=nil{http.Error(w,"audit unavailable",http.StatusServiceUnavailable);return}
  page,err:=lister.List(r.Context(),actorID,organizationID,limit,cursor)
  if errors.Is(err,auth.ErrAuditAccessDenied){http.Error(w,"forbidden",http.StatusForbidden);return}
  if errors.Is(err,auth.ErrInvalidAuditQuery){http.Error(w,"invalid query",http.StatusBadRequest);return}
  if err!=nil{http.Error(w,"audit unavailable",http.StatusServiceUnavailable);return}
  out:=auditPageJSON{Events:make([]auditEventJSON,0,len(page.Events)),NextCursor:encodeAuditCursor(page.NextCursor)}
  for _,event:=range page.Events{
   row:=auditEventJSON{ID:event.ID,OrganizationID:event.OrganizationID,Type:event.Type,OccurredAt:event.OccurredAt}
   if event.ActorUserID.Valid{v:=event.ActorUserID.String;row.ActorUserID=&v}
   if event.TargetUserID.Valid{v:=event.TargetUserID.String;row.TargetUserID=&v}
   out.Events=append(out.Events,row)
  }
  w.Header().Set("X-Content-Type-Options","nosniff")
  jsonResponse(w,http.StatusOK,out)
 })
}

// AuditRoute is an explicit path pattern for trusted server wiring.
const AuditRoute = "GET /api/v1/organizations/{organizationID}/audit-events"

var _ = strings.TrimSpace
