FROM rasa/rasa:3.6.19-full

# Cambiar al directorio de trabajo
WORKDIR /app

# Copiar requirements.txt si existe
COPY requirements.txt* ./

# Instalar dependencias adicionales si existen
USER root
RUN if [ -f requirements.txt ]; then pip install --no-cache-dir --default-timeout=100 -r requirements.txt; fi
USER 1001

# Copiar el proyecto Rasa
COPY . .

# Crear directorio para modelos si no existe
RUN mkdir -p models

# Verificar que los archivos necesarios existen
RUN ls -la
RUN if [ ! -f domain.yml ]; then echo "ERROR: domain.yml not found"; exit 1; fi
RUN if [ ! -f config.yml ]; then echo "ERROR: config.yml not found"; exit 1; fi
RUN if [ ! -d data ]; then echo "ERROR: data directory not found"; exit 1; fi

# Entrenar el modelo primero
RUN rasa train --fixed-model-name current-model

# Exponer puerto
EXPOSE 5005

# Script de inicio
COPY <<EOF /app/start.sh
#!/bin/bash
set -e

# Entrenar si no existe modelo o si han cambiado los archivos
if [ ! -f models/current-model.tar.gz ] || [ data -nt models/current-model.tar.gz ] || [ domain.yml -nt models/current-model.tar.gz ]; then
    echo "Training model..."
    rasa train --fixed-model-name current-model
fi

echo "Starting Rasa server..."
exec rasa run --enable-api --cors "*" --debug --model models/current-model.tar.gz
EOF

RUN chmod +x /app/start.sh

# Comando por defecto
CMD ["/app/start.sh"]