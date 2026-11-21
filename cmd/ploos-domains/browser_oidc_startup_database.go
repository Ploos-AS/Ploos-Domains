package main

import (
 "context"
 "database/sql"
 "strings"
 "time"
)

// checkBrowserOIDCStartupDatabase checks a dedicated browser-login DSN.
// The connection is always closed, including when readiness fails.
func checkBrowserOIDCStartupDatabase(ctx context.Context,dsn string)error{
 if ctx==nil||strings.TrimSpace(dsn)==""{return errBrowserOIDCDatabase}
 db,err:=sql.Open("postgres",dsn)
 if err!=nil{return errBrowserOIDCDatabase}
 defer db.Close()
 db.SetMaxOpenConns(2)
 db.SetMaxIdleConns(0)
 db.SetConnMaxLifetime(time.Minute)
 probeCtx,cancel:=context.WithTimeout(ctx,5*time.Second)
 defer cancel()
 if err:=db.PingContext(probeCtx);err!=nil{return errBrowserOIDCDatabase}
 return checkBrowserOIDCDatabase(probeCtx,db)
}
