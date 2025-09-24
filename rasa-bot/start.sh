#!/bin/bash
set -e

# ✅ Verificar permisos del script al inicio
echo "🔍 Checking script permissions..."
ls -la /app/start.sh

export RASA_HOME=/app/.rasa

echo "🚀 Starting Rasa server with pre-trained model..."

# Verificar que el modelo existe
if [ ! -f models/current-model.tar.gz ]; then
    echo "❌ ERROR: Pre-trained model not found!"
    echo "Available files in models/:"
    ls -la models/ || echo "No models directory found"
    
    # Intentar entrenar modelo si no existe
    echo "🤖 Training model as fallback..."
    if [ -f domain.yml ] && [ -f config.yml ] && [ -d data ]; then
        rasa train --fixed-model-name current-model
    else
        echo "❌ Cannot train model: missing configuration files"
        exit 1
    fi
fi

echo "✅ Model found: models/current-model.tar.gz"
echo "Starting Rasa server..."

# ✅ Usar exec para proper signal handling
exec rasa run --enable-api --cors "*" --model models/current-model.tar.gz --port 5005