# ===== STAGE 1: BUILDER =====
FROM golang:1.24-alpine AS builder

# Instalar git para go mod download
RUN apk add --no-cache git

WORKDIR /build

# Copiar módulos Go
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copiar código fuente
COPY api/ .

# Compilación estática (amd64)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -ldflags="-w -s" -o server ./cmd/api

# Verificar que se construyó
RUN ls -la /build/server


# ===== STAGE 2: FINAL =====
FROM alpine:latest

# Instalar dependencias mínimas de runtime
RUN apk --no-cache add ca-certificates curl tzdata

# Crear directorio de la app
WORKDIR /app

# Copiar binario desde builder
COPY --from=builder /build/server /app/server

# Dar permisos de ejecución
RUN chmod +x /app/server

# Crear usuario no root (opcional)
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup && \
    chown appuser:appgroup /app/server
USER appuser

# Exponer puerto
EXPOSE 8080

# Ejecutar binario
CMD ["/app/server"]