FROM rasa/rasa:3.6.19-full

# Cambiar al directorio de trabajo
WORKDIR /app

# Crear directorios necesarios
RUN mkdir -p .rasa models /tmp/cache_db

# Copiar requirements.txt si existe
COPY requirements.txt* ./

# Instalar dependencias adicionales si existen
USER root
RUN if [ -f requirements.txt ]; then pip install --no-cache-dir --default-timeout=100 -r requirements.txt; fi

# Copiar el proyecto Rasa
COPY . .

# Configurar permisos
RUN chown -R 1001:1001 /app && chmod -R 755 /app

# Cambiar a usuario no root
USER 1001

# Verificar archivos necesarios
RUN test -f domain.yml || (echo "ERROR: domain.yml not found" && exit 1)
RUN test -f config.yml || (echo "ERROR: config.yml not found" && exit 1)
RUN test -d data || (echo "ERROR: data directory not found" && exit 1)

# Variables de entorno para el build
ENV RASA_HOME=/app/.rasa
ENV SQLALCHEMY_WARN_20=1

# ⭐ ENTRENAR EL MODELO DURANTE EL BUILD ⭐
RUN echo "🤖 Training Rasa model during build..." && \
    rasa train --fixed-model-name current-model --no-cache && \
    echo "✅ Model training completed!" && \
    ls -la models/

# Crear script de inicio simplificado (sin entrenamiento)
RUN cat > start.sh << 'EOF'
#!/bin/bash
set -e

export RASA_HOME=/app/.rasa

echo "🚀 Starting Rasa server with pre-trained model..."

# Verificar que el modelo existe
if [ ! -f models/current-model.tar.gz ]; then
    echo "❌ ERROR: Pre-trained model not found!"
    echo "Available files in models/:"
    ls -la models/ || echo "No models directory found"
    exit 1
fi

echo "✅ Model found: models/current-model.tar.gz"
echo "Starting Rasa server..."
exec rasa run --enable-api --cors "*" --model models/current-model.tar.gz --port 5005
EOF

RUN chmod +x start.sh

# Exponer puerto
EXPOSE 5005

# Health check (ahora con menos tiempo de espera)
HEALTHCHECK --interval=15s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:5005 || exit 1

# Comando por defecto - ejecutar directamente con shell
CMD ["/bin/bash", "/app/start.sh"]