# --- Builder Stage ---
FROM golang:1.23-alpine AS builder

# Instalar dependencias necesarias para compilar (git, libc, etc.)
RUN apk add --no-cache git gcc g++ libc-dev

WORKDIR /app

# Copiar archivos de go y descargar módulos
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código
COPY . .

# Compilar el binario (optimizado)
RUN go build -ldflags="-w -s" -o api-zoo .

# --- Runtime Stage ---
FROM alpine:3.19

# Instalar las dependencias necesarias para ejecutar binario
RUN apk add --no-cache libstdc++

WORKDIR /app

# Copiar binario y recursos desde la etapa de construcción
COPY --from=builder /app/api-zoo .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/.env .env


EXPOSE 8080

CMD ["./api-zoo"]
