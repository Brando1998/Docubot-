.PHONY: help up-local down-local build-all logs-api logs-rasa logs-playwright logs-baileys clean

help: ## Mostrar ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

up-local: ## Levantar entorno local con docker-compose
	@echo "🚀 Levantando entorno local..."
	docker compose up --build -d
	@echo "✅ Entorno levantado. Verificando servicios..."
	@echo "📊 Rasa: http://localhost:5005"
	@echo "🎭 Playwright: http://localhost:3001"
	@echo "🔧 API: http://localhost:8080"
	@echo "💬 Baileys: http://localhost:3000"

up-sequential: ## Levantar servicios secuencialmente (recomendado)
	@echo "🚀 Levantando servicios en orden..."
	@echo "1️⃣ Levantando bases de datos..."
	docker compose up -d postgres mongodb
	@echo "⏳ Esperando bases de datos..."
	sleep 10
	@echo "2️⃣ Levantando Rasa..."
	docker compose up -d rasa
	@echo "⏳ Esperando Rasa..."
	sleep 30
	@echo "3️⃣ Levantando Playwright..."
	docker compose up -d playwright
	@echo "⏳ Esperando Playwright..."
	sleep 15
	@echo "4️⃣ Levantando API..."
	docker compose up -d api
	@echo "⏳ Esperando API..."
	sleep 15
	@echo "5️⃣ Levantando Baileys..."
	docker compose up -d baileys
	@echo "✅ Todos los servicios levantados!"

down-local: ## Detener entorno local
	@echo "🛑 Deteniendo entorno local..."
	docker compose down

build-all: ## Construir todas las imágenes
	@echo "🔨 Construyendo todas las imágenes..."
	docker compose build

logs-api: ## Ver logs del API
	docker compose logs -f api

logs-rasa: ## Ver logs de Rasa
	docker compose logs -f rasa

logs-playwright: ## Ver logs de Playwright
	docker compose logs -f playwright

logs-baileys: ## Ver logs de Baileys
	docker compose logs -f baileys

logs-all: ## Ver logs de todos los servicios
	docker compose logs -f

status: ## Verificar estado de servicios
	@echo "📊 Estado de los servicios:"
	@docker compose ps

health-check: ## Verificar salud de servicios
	@echo "🏥 Verificando salud de servicios..."
	@echo "Postgres:" && curl -f http://localhost:5432 2>/dev/null && echo "✅" || echo "❌"
	@echo "Rasa:" && curl -f http://localhost:5005/status 2>/dev/null && echo "✅" || echo "❌"
	@echo "Playwright:" && curl -f http://localhost:3001/health 2>/dev/null && echo "✅" || echo "❌"
	@echo "API:" && curl -f http://localhost:8080/health 2>/dev/null && echo "✅" || echo "❌"
	@echo "Baileys:" && curl -f http://localhost:3000/health 2>/dev/null && echo "✅" || echo "❌"

clean: ## Limpiar contenedores, imágenes y volúmenes
	@echo "🧹 Limpiando..."
	docker compose down -v --remove-orphans
	docker system prune -f
	@echo "✅ Limpieza completa"

restart: ## Reiniciar todos los servicios
	@echo "🔄 Reiniciando servicios..."
	docker compose restart

restart-api: ## Reiniciar solo la API
	docker compose restart api

restart-rasa: ## Reiniciar solo Rasa
	docker compose restart rasa

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