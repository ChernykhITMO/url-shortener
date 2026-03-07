FROM golang:1.25-alpine AS app-builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/url-shortener ./cmd/app

FROM golang:1.25-alpine AS goose-builder
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.0

FROM alpine:3.20 AS app
WORKDIR /app

COPY --from=app-builder /out/url-shortener /app/url-shortener
COPY config/local.yaml /app/config/local.yaml

EXPOSE 8080
CMD ["/app/url-shortener", "--storage=postgres", "--config=/app/config/local.yaml"]

FROM alpine:3.20 AS migrate
WORKDIR /migrations

COPY --from=goose-builder /go/bin/goose /usr/local/bin/goose
COPY migrations/migrations ./migrations

CMD ["goose", "-dir", "/migrations/migrations", "postgres", "postgres://postgres:postgres@db:5432/postgres?sslmode=disable", "up"]