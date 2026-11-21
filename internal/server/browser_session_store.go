package server

import (
 "context"
 "crypto/sha256"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

var ErrBrowserSessionUnavailable = errors.New("browser session unavailable")
var ErrBrowserSessionInvalid = errors.New("invalid browser session")

// BrowserSessionStore persists only token digests. Callers must establish
// authenticated user identity before creating a session.
type BrowserSessionStore struct { DB *sql.DB }

func (s BrowserSessionStore) Create(ctx context.Context,userID,rawToken string,createdAt time.Time,ttl time.Duration)error{
 if s.DB==nil{return ErrBrowserSessionUnavailable}
 if ctx==nil||userID==""||rawToken==""||createdAt.IsZero()||(BrowserSessionConfig{TTL:ttl}).Validate()!=nil{return ErrBrowserSessionInvalid}
 digest:=sha256.Sum256([]byte(rawToken))
 _,err:=s.DB.ExecContext(ctx,`INSERT INTO browser_sessions(token_sha256,user_id,created_at,expires_at)
 VALUES($1,$2::uuid,$3,$4)`,digest[:],userID,createdAt,createdAt.Add(ttl))
 if err!=nil{return fmt.Errorf("create browser session: %w",err)}
 return nil
}

// Resolve rejects expired or revoked sessions using database time. The raw
// cookie token is never returned or written to the database.
func (s BrowserSessionStore) Resolve(ctx context.Context,rawToken string)(string,error){
 if s.DB==nil{return "",ErrBrowserSessionUnavailable}
 if ctx==nil||rawToken==""{return "",ErrBrowserSessionInvalid}
 digest:=sha256.Sum256([]byte(rawToken))
 var userID string
 err:=s.DB.QueryRowContext(ctx,`SELECT user_id::text FROM browser_sessions
 WHERE token_sha256=$1 AND revoked_at IS NULL AND expires_at>now()`,digest[:]).Scan(&userID)
 if errors.Is(err,sql.ErrNoRows){return "",ErrBrowserSessionInvalid}
 if err!=nil{return "",fmt.Errorf("resolve browser session: %w",err)}
 return userID,nil
}

// Revoke invalidates one specific session, leaving other user sessions intact.
func (s BrowserSessionStore) Revoke(ctx context.Context,rawToken string)error{
 if s.DB==nil{return ErrBrowserSessionUnavailable}
 if ctx==nil||rawToken==""{return ErrBrowserSessionInvalid}
 digest:=sha256.Sum256([]byte(rawToken))
 _,err:=s.DB.ExecContext(ctx,`UPDATE browser_sessions SET revoked_at=now()
 WHERE token_sha256=$1 AND revoked_at IS NULL`,digest[:])
 if err!=nil{return fmt.Errorf("revoke browser session: %w",err)}
 return nil
}
