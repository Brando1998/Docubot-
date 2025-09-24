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