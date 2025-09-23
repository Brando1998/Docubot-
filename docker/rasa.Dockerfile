FROM rasa/rasa:3.6.19-full

# Cambiar al directorio de trabajo
WORKDIR /app

# Copiar archivos de configuración
COPY requirements.txt requirements.txt

# Instalar dependencias adicionales
USER root
RUN pip install --no-cache-dir --default-timeout=100 -r requirements.txt
USER 1001

# Copiar el proyecto Rasa
COPY . .

# Crear directorio para modelos si no existe
RUN mkdir -p models

# Exponer puerto
EXPOSE 5005

# Comando por defecto (se sobrescribe en docker-compose)
CMD ["rasa", "run", "--enable-api", "--cors", "*", "--debug"]