package main

import (
 "log"
 "net/http"
 "os"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/server"
)

func main() {
 addr := os.Getenv("PLOOS_LISTEN_ADDR")
 if addr == "" { addr = ":8080" }
 srv := &http.Server{
  Addr: addr,
  Handler: server.NewHandler(),
  ReadHeaderTimeout: 5 * time.Second,
  ReadTimeout: 15 * time.Second,
  WriteTimeout: 15 * time.Second,
  IdleTimeout: 60 * time.Second,
 }
 log.Printf("Ploos Domains listening on %s", addr)
 if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal(err) }
}
