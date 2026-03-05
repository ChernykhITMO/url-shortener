FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /url-shortener ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /url-shortener /app/url-shortener

RUN mkdir -p /app/config

COPY config/local.yaml /app/config/local.yaml

EXPOSE 8080

CMD ["/app/url-shortener", "--storage=postgres", "--config=/app/config/local.yaml"]

