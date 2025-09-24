# Variables del proyecto
PROJECT_NAME := docubot
COMPOSE_FILE := docker-compose.yml

.PHONY: help up-local down-local build-all logs-api logs-rasa logs-playwright logs-baileys clean clean-project clean-all

help: ## Mostrar ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

up-local: ## Levantar entorno local con docker-compose
	@echo "🚀 Levantando entorno local..."
	docker compose -f $(COMPOSE_FILE) up --build -d
	@echo "✅ Entorno levantado. Verificando servicios..."
	@echo "📊 Rasa: http://localhost:5005"
	@echo "🎭 Playwright: http://localhost:3001"
	@echo "🔧 API: http://localhost:8080"
	@echo "💬 Baileys: http://localhost:3000"

up-sequential: ## Levantar servicios secuencialmente (recomendado)
	@echo "🚀 Levantando servicios en orden..."
	@echo "1️⃣ Levantando bases de datos..."
	docker compose -f $(COMPOSE_FILE) up -d postgres mongodb
	@echo "⏳ Esperando bases de datos..."
	sleep 10
	@echo "2️⃣ Levantando Rasa..."
	docker compose -f $(COMPOSE_FILE) up -d rasa
	@echo "⏳ Esperando Rasa..."
	sleep 30
	@echo "3️⃣ Levantando Playwright..."
	docker compose -f $(COMPOSE_FILE) up -d playwright
	@echo "⏳ Esperando Playwright..."
	sleep 15
	@echo "4️⃣ Levantando API..."
	docker compose -f $(COMPOSE_FILE) up -d api
	@echo "⏳ Esperando API..."
	sleep 15
	@echo "5️⃣ Levantando Baileys..."
	docker compose -f $(COMPOSE_FILE) up -d baileys
	@echo "✅ Todos los servicios levantados!"

down-local: ## Detener entorno local
	@echo "🛑 Deteniendo entorno local..."
	docker compose -f $(COMPOSE_FILE) down

build-all: ## Construir todas las imágenes del proyecto
	@echo "🔨 Construyendo todas las imágenes..."
	docker compose -f $(COMPOSE_FILE) build

logs-api: ## Ver logs del API
	docker compose -f $(COMPOSE_FILE) logs -f api

logs-rasa: ## Ver logs de Rasa
	docker compose -f $(COMPOSE_FILE) logs -f rasa

logs-playwright: ## Ver logs de Playwright
	docker compose -f $(COMPOSE_FILE) logs -f playwright

logs-baileys: ## Ver logs de Baileys
	docker compose -f $(COMPOSE_FILE) logs -f baileys

logs-all: ## Ver logs de todos los servicios del proyecto
	docker compose -f $(COMPOSE_FILE) logs -f

status: ## Verificar estado de servicios del proyecto
	@echo "📊 Estado de los servicios de $(PROJECT_NAME):"
	@docker compose -f $(COMPOSE_FILE) ps

health-check: ## Verificar salud de servicios del proyecto
	@echo "🏥 Verificando salud de servicios..."
	@echo -n "Postgres: " && (curl -f http://localhost:5433 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "MongoDB: " && (curl -f http://localhost:27018 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "Rasa: " && (curl -f http://localhost:5005/status 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "Playwright: " && (curl -f http://localhost:3001/health 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "API: " && (curl -f http://localhost:8080/health 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "Baileys: " && (curl -f http://localhost:3000/health 2>/dev/null && echo "✅" || echo "❌")

clean: clean-project ## Limpiar SOLO los contenedores, imágenes y volúmenes de este proyecto

clean-project: ## Limpiar contenedores, imágenes y volúmenes específicos del proyecto
	@echo "🧹 Limpiando recursos del proyecto $(PROJECT_NAME)..."
	@echo "⏹️  Deteniendo contenedores del proyecto..."
	-docker compose -f $(COMPOSE_FILE) down -v --remove-orphans 2>/dev/null
	@echo "🗑️  Eliminando contenedores del proyecto..."
	-docker container rm -f $(PROJECT_NAME)-postgres $(PROJECT_NAME)-mongo $(PROJECT_NAME)-rasa $(PROJECT_NAME)-playwright $(PROJECT_NAME)-api $(PROJECT_NAME)-baileys 2>/dev/null || true
	@echo "🖼️  Eliminando imágenes del proyecto..."
	-docker image rm -f $$(docker images --filter "reference=$(PROJECT_NAME)*" -q) 2>/dev/null || true
	-docker image rm -f $$(docker images --filter "reference=docubot*" -q) 2>/dev/null || true
	@echo "💾 Eliminando volúmenes del proyecto..."
	-docker volume rm -f $$(docker volume ls --filter "name=$(PROJECT_NAME)" -q) 2>/dev/null || true
	-docker volume rm -f $$(docker volume ls --filter "name=docubot" -q) 2>/dev/null || true
	@echo "📁 Limpiando volúmenes huérfanos relacionados..."
	-docker volume rm -f $$(docker volume ls -f "label=com.docker.compose.project=$(PROJECT_NAME)" -q) 2>/dev/null || true
	@echo "🌐 Eliminando red del proyecto..."
	-docker network rm $(PROJECT_NAME)-network 2>/dev/null || true
	-docker network rm docubot-network 2>/dev/null || true
	@echo "✅ Limpieza del proyecto $(PROJECT_NAME) completada"

clean-all: ## ⚠️  PELIGROSO: Limpiar TODO el sistema Docker (usar con cuidado)
	@echo "⚠️  ADVERTENCIA: Esto eliminará TODOS los contenedores, imágenes, volúmenes y redes del sistema"
	@echo "¿Estás seguro? Presiona Ctrl+C para cancelar, o Enter para continuar..."
	@read dummy
	@echo "🧹 Limpiando TODO el sistema Docker..."
	docker system prune -af --volumes
	@echo "✅ Limpieza completa del sistema"

restart: ## Reiniciar todos los servicios del proyecto
	@echo "🔄 Reiniciando servicios del proyecto..."
	docker compose -f $(COMPOSE_FILE) restart

restart-api: ## Reiniciar solo la API
	docker compose -f $(COMPOSE_FILE) restart api

restart-rasa: ## Reiniciar solo Rasa
	docker compose -f $(COMPOSE_FILE) restart rasa

restart-baileys: ## Reiniciar solo Baileys
	docker compose -f $(COMPOSE_FILE) restart baileys

# Comandos de desarrollo
dev-logs: ## Ver logs en tiempo real de todos los servicios
	docker compose -f $(COMPOSE_FILE) logs -f --tail=100

dev-shell-api: ## Abrir shell en el contenedor de la API
	docker compose -f $(COMPOSE_FILE) exec api /bin/sh

dev-shell-rasa: ## Abrir shell en el contenedor de Rasa
	docker compose -f $(COMPOSE_FILE) exec rasa /bin/bash

# Comandos Kubernetes
k8s-deploy: ## Desplegar en Kubernetes
	kubectl apply -f k8s/configmaps/
	kubectl apply -f k8s/secrets/
	kubectl apply -f k8s/deployments/
	kubectl apply -f k8s/services/

k8s-delete: ## Eliminar despliegue de Kubernetes
	kubectl delete -f k8s/services/
	kubectl delete -f k8s/deployments/
	kubectl delete -f k8s/secrets/
	kubectl delete -f k8s/configmaps/

# Comandos de información
show-containers: ## Mostrar solo los contenedores de este proyecto
	@echo "📦 Contenedores del proyecto $(PROJECT_NAME):"
	@docker ps -a --filter "name=$(PROJECT_NAME)" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

show-images: ## Mostrar solo las imágenes de este proyecto
	@echo "🖼️  Imágenes del proyecto $(PROJECT_NAME):"
	@docker images --filter "reference=$(PROJECT_NAME)*" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

show-volumes: ## Mostrar solo los volúmenes de este proyecto
	@echo "💾 Volúmenes del proyecto $(PROJECT_NAME):"
	@docker volume ls --filter "name=$(PROJECT_NAME)" --format "table {{.Name}}\t{{.Size}}"