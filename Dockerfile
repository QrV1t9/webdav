FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o webdav ./cmd/

RUN go build -o argon ./cmd/argon

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/webdav .
COPY --from=builder /app/argon .
COPY config/prod.yaml /app/config/prod.yaml

ENV CONFIG_PATH=/app/config/prod.yaml

CMD ["./webdav"]