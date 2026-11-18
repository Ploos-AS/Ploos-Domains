FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/ploos-domains ./cmd/ploos-domains

FROM alpine:3.21
RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app
COPY --from=build /out/ploos-domains /usr/local/bin/ploos-domains
USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/ploos-domains"]
