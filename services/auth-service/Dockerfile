FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/auth-service ./cmd

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /out/auth-service /app/auth-service
COPY .env /app/.env
COPY migrations /app/migrations
COPY --chmod=0755 docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

EXPOSE 8001

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/app/auth-service"]
