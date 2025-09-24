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

# Variables de entorno para el build
ENV RASA_HOME=/app/.rasa
ENV SQLALCHEMY_WARN_20=1

# Verificar archivos necesarios
RUN test -f domain.yml || (echo "ERROR: domain.yml not found" && exit 1)
RUN test -f config.yml || (echo "ERROR: config.yml not found" && exit 1)
RUN test -d data || (echo "ERROR: data directory not found" && exit 1)

# ⭐ ENTRENAR EL MODELO DURANTE EL BUILD ⭐
RUN echo "🤖 Training Rasa model during build..." && \
    rasa train --fixed-model-name current-model && \
    echo "✅ Model training completed!" && \
    ls -la models/

# Crear script de inicio como ROOT
RUN echo '#!/bin/bash' > /app/start.sh && \
    echo 'set -e' >> /app/start.sh && \
    echo '' >> /app/start.sh && \
    echo 'export RASA_HOME=/app/.rasa' >> /app/start.sh && \
    echo '' >> /app/start.sh && \
    echo 'echo "🚀 Starting Rasa server with pre-trained model..."' >> /app/start.sh && \
    echo '' >> /app/start.sh && \
    echo '# Verificar que el modelo existe' >> /app/start.sh && \
    echo 'if [ ! -f models/current-model.tar.gz ]; then' >> /app/start.sh && \
    echo '    echo "❌ ERROR: Pre-trained model not found!"' >> /app/start.sh && \
    echo '    echo "Available files in models/:"' >> /app/start.sh && \
    echo '    ls -la models/ || echo "No models directory found"' >> /app/start.sh && \
    echo '    exit 1' >> /app/start.sh && \
    echo 'fi' >> /app/start.sh && \
    echo '' >> /app/start.sh && \
    echo 'echo "✅ Model found: models/current-model.tar.gz"' >> /app/start.sh && \
    echo 'echo "Starting Rasa server..."' >> /app/start.sh && \
    echo 'exec rasa run --enable-api --cors "*" --model models/current-model.tar.gz --port 5005' >> /app/start.sh

# Configurar permisos DESPUÉS de crear todo
RUN chown -R 1001:1001 /app && chmod -R 755 /app && chmod +x /app/start.sh

# Cambiar a usuario no root
USER 1001

# Exponer puerto
EXPOSE 5005

# Health check
HEALTHCHECK --interval=15s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:5005 || exit 1

# Usar ENTRYPOINT y CMD como configuraste
ENTRYPOINT ["/app/start.sh"]
CMD []