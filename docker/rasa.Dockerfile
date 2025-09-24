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

# Exponer puerto
EXPOSE 5005

# Variables de entorno
ENV RASA_HOME=/app/.rasa
ENV SQLALCHEMY_WARN_20=1

# Crear script de inicio como archivo separado
RUN cat > start.sh << 'EOF'
#!/bin/bash
set -e

export RASA_HOME=/app/.rasa

echo "Checking if model needs training..."

if [ ! -f models/current-model.tar.gz ]; then
    echo "No model found, training..."
    rasa train --fixed-model-name current-model --no-cache || exit 1
elif [ data -nt models/current-model.tar.gz ] || [ domain.yml -nt models/current-model.tar.gz ] || [ config.yml -nt models/current-model.tar.gz ]; then
    echo "Files newer than model, retraining..."
    rasa train --fixed-model-name current-model --no-cache || exit 1
fi

echo "Starting Rasa server..."
exec rasa run --enable-api --cors "*" --model models/current-model.tar.gz
EOF

RUN chmod +x start.sh

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=180s --retries=5 \
    CMD curl -f http://localhost:5005 || exit 1

# Comando por defecto - ejecutar el script con bash
CMD ["bash", "start.sh"]