# Multi-stage build para resolver problemas de inconsistencia
FROM golang:1.23-alpine AS builder

# Instalar git para go mod download
RUN apk add --no-cache git

WORKDIR /build

# Copiar módulos Go
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copiar código fuente
COPY api/ .

# Build estático
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o server ./cmd/api

# Verificar que se construyó
RUN ls -la server && file server

# ===== STAGE FINAL =====
FROM alpine:latest

# Instalar dependencias de runtime
RUN apk --no-cache add ca-certificates curl tzdata

WORKDIR /app

# Copiar binario desde builder
COPY --from=builder /build/server ./server

# Hacer ejecutable (por si acaso)
RUN chmod +x ./server

# Verificar que está ahí
RUN ls -la ./server && file ./server

# Usuario no-root (opcional, puedes comentar si causa problemas)
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup && \
    chown appuser:appgroup ./server
USER appuser

EXPOSE 8080

CMD ["./server"]