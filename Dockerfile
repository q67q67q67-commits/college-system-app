# API для колледжа
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api ./cmd/api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /api /api
COPY web /app/web
WORKDIR /app
ENV NC_ADDR=:8080
ENV NC_DB_URL=postgres://nc:nc_dev@postgres:5432/nc_db?sslmode=disable
EXPOSE 8080
CMD ["/api"]
