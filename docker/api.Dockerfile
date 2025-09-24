FROM golang:1.23

WORKDIR /app

# Copiar go.mod y go.sum desde la carpeta api
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copiar todo el código de la API
COPY api/ .

# Compilar la aplicación
RUN go build -o server ./cmd/api

EXPOSE ${PORT:-8080}

CMD ["./server"]