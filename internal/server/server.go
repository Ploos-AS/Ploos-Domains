// Package server exposes the initial Ploos Domains HTTP API.
package server

import (
 "encoding/json"
 "net/http"
)

type statusResponse struct {
 Status string `json:"status"`
 Service string `json:"service"`
 Version string `json:"version"`
}

func jsonResponse(w http.ResponseWriter, code int, payload any) {
 w.Header().Set("Content-Type", "application/json; charset=utf-8")
 w.Header().Set("Cache-Control", "no-store")
 w.WriteHeader(code)
 _ = json.NewEncoder(w).Encode(payload)
}

// NewHandler builds an HTTP handler without external dependencies.
func NewHandler() http.Handler {
 mux := http.NewServeMux()
 mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
  jsonResponse(w, http.StatusOK, statusResponse{Status:"ok", Service:"ploos-domains", Version:"dev"})
 })
 mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
  // Readiness is process-only until database and external dependencies are wired.
  jsonResponse(w, http.StatusOK, statusResponse{Status:"ready", Service:"ploos-domains", Version:"dev"})
 })
 mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
  jsonResponse(w, http.StatusOK, map[string]string{"name":"Ploos Domains API","version":"v1","status":"foundation"})
 })
 return mux
}
