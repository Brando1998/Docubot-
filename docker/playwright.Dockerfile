FROM node:20-alpine

# Instalar dependencias del sistema para Playwright
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

# Configurar Chromium
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
ENV CHROMIUM_PATH=/usr/bin/chromium-browser

WORKDIR /app

# Copiar package.json
COPY package*.json ./

# Instalar dependencias (usar npm install si no hay package-lock.json)
RUN npm install --only=production

# Instalar Playwright
RUN npx playwright install chromium

# Copiar el código fuente
COPY . .

# Crear usuario no-root
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nodejs -u 1001

# Cambiar permisos
RUN chown -R nodejs:nodejs /app
USER nodejs

# Exponer puerto
EXPOSE 3001

# Comando de inicio
CMD ["npm", "start"]