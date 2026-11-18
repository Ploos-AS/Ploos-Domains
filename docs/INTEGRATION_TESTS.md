# PostgreSQL integration tests

Run against a **disposable** PostgreSQL 16+ database with migrations 000001–000003 applied.

```sh
export PLOOS_TEST_DATABASE_URL='postgres://ploos:password@localhost:5432/ploos_domains_test?sslmode=disable'
go test -mod=mod -race -tags=integration -count=1 ./internal/verification
```

The integration suite tests the actual Go `CompleteVerification` transaction against PostgreSQL for success, replay rejection, expiry, missing challenges, and concurrent attempts. It creates isolated fixture organizations and cleans them up.

The PostgreSQL driver `github.com/lib/pq` is currently test-only. CI resolves the module using `-mod=mod`; a committed `go.sum` should be added once dependency resolution is available in the development environment.

This suite **does not** yet test authenticated HTTP requests, live DNS TXT lookups, or organization authorization through a real database-backed authorizer. Those remain required before exposing a public API.
