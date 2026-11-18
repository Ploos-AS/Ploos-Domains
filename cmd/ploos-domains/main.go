package main

import (
 "context"
 "log"
 "net/http"
 "os"
 "time"

)

func main() {
 addr := os.Getenv("PLOOS_LISTEN_ADDR")
 if addr == "" { addr = ":8080" }
 handler,closeDB,err:=buildHandler(context.Background())
 if err!=nil{log.Fatal(err)}
 defer closeDB()
 srv := &http.Server{
  Addr: addr,
  Handler: handler,
  ReadHeaderTimeout: 5 * time.Second,
  ReadTimeout: 15 * time.Second,
  WriteTimeout: 15 * time.Second,
  IdleTimeout: 60 * time.Second,
 }
 log.Printf("Ploos Domains listening on %s", addr)
 if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal(err) }
}
