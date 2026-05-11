FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/auth-service ./cmd

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /out/auth-service /app/auth-service

EXPOSE 8001

CMD ["/app/auth-service"]
