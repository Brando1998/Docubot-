FROM node:20-alpine AS builder

WORKDIR /app

# Copiar archivos de configuración
COPY package*.json ./
COPY vite.config.js ./
COPY tsconfig*.json ./
COPY index.html ./

# Instalar dependencias
RUN npm ci

# Copiar código fuente
COPY src/ ./src/
COPY public/ ./public/

# Build para producción
RUN npm run build

# Imagen final con servidor ligero
FROM nginx:alpine

# Copiar archivos compilados
COPY --from=builder /app/dist /usr/share/nginx/html

# Configuración personalizada de Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Health check endpoint para nginx
RUN echo '<!DOCTYPE html><html><head><title>Health</title></head><body>OK</body></html>' > /usr/share/nginx/html/health

# Crear usuario no-root
RUN addgroup -g 1001 -S nginx
RUN adduser -S nginx -u 1001

EXPOSE 3002

CMD ["nginx", "-g", "daemon off;"]