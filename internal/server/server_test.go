package server

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"
)

func TestEndpoints(t *testing.T) {
 cases := []struct{ path string; status int; expected string }{
  {"/healthz", http.StatusOK, "ok"},
  {"/readyz", http.StatusOK, "ready"},
  {"/api/v1", http.StatusOK, "foundation"},
  {"/missing", http.StatusNotFound, ""},
 }
 for _, tc := range cases {
  t.Run(tc.path, func(t *testing.T) {
   rec := httptest.NewRecorder()
   NewHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
   if rec.Code != tc.status { t.Fatalf("status: got %d, want %d", rec.Code, tc.status) }
   if tc.expected != "" {
    if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" { t.Fatalf("content type: %q", got) }
    var data map[string]string
    if err := json.Unmarshal(rec.Body.Bytes(), &data); err != nil { t.Fatal(err) }
    if data["status"] != tc.expected { t.Fatalf("status payload: %q, want %q", data["status"], tc.expected) }
   }
  })
 }
}

func TestMethodNotAllowed(t *testing.T) {
 rec := httptest.NewRecorder()
 NewHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/healthz", nil))
 if rec.Code != http.StatusMethodNotAllowed { t.Fatalf("got %d", rec.Code) }
}
