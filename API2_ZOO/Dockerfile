# --- Builder Stage ---
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git gcc g++ libc-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .env

RUN go build -ldflags="-w -s" -o api-zoo .

# --- Runtime Stage ---
FROM alpine:3.19

RUN apk add --no-cache libstdc++

WORKDIR /app

COPY --from=builder /app/api-zoo .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/.env .env

EXPOSE 8080

CMD ["./api-zoo"]
