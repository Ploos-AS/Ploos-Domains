# Local development

Requires Go 1.23 or newer.

```sh
go run ./cmd/ploos-domains
```

The default listener is `:8080`. Override with `PLOOS_LISTEN_ADDR=127.0.0.1:8081`.

```sh
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
curl http://localhost:8080/api/v1
go test ./...
go vet ./...
```

The current readiness endpoint checks only the running process. It must check PostgreSQL once persistence is introduced. No authentication, database, DNS mutations or registrar calls are implemented in this foundation PR. Do not expose it as a production management API.
