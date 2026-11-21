package main

import (
 "context"
 "database/sql"
 "errors"
 "testing"
)

func TestBrowserOIDCDatabaseRejectsNil(t *testing.T){
 if err:=checkBrowserOIDCDatabase(context.Background(),nil);!errors.Is(err,errBrowserOIDCDatabase){t.Fatalf("nil db: %v",err)}
 if err:=checkBrowserOIDCDatabase(nil,&sql.DB{});!errors.Is(err,errBrowserOIDCDatabase){t.Fatalf("nil ctx: %v",err)}
}
