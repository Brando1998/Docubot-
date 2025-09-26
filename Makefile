# Variables del proyecto
PROJECT_NAME := docubot
COMPOSE_FILE := docker-compose.yml

.PHONY: help up-local down-local build-all logs-api logs-vue logs-rasa logs-playwright logs-baileys clean clean-project clean-all

help: ## Mostrar ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

up-local: ## Levantar entorno local con docker-compose
	@echo "🚀 Levantando entorno local..."
	docker compose -f $(COMPOSE_FILE) up --build -d
	@echo "✅ Entorno levantado. Verificando servicios..."
	@echo "🎨 Vue Dashboard: http://localhost:3002"
	@echo "🔧 API Go: http://localhost:8080"
	@echo "📊 Rasa: http://localhost:5005"
	@echo "🎭 Playwright: http://localhost:3001"
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
	sleep 15
	@echo "3️⃣ Levantando Playwright..."
	docker compose -f $(COMPOSE_FILE) up -d playwright
	@echo "⏳ Esperando Playwright..."
	sleep 15
	@echo "4️⃣ Levantando API..."
	docker compose -f $(COMPOSE_FILE) up -d api
	@echo "⏳ Esperando API..."
	sleep 15
	@echo "5️⃣ Levantando Vue Dashboard..."
	docker compose -f $(COMPOSE_FILE) up -d vue
	@echo "⏳ Esperando Vue..."
	sleep 10
	@echo "6️⃣ Levantando Baileys..."
	docker compose -f $(COMPOSE_FILE) up -d baileys
	@echo "✅ Todos los servicios levantados!"
	
down-local: ## Detener entorno local
	@echo "🛑 Deteniendo entorno local..."
	docker compose -f $(COMPOSE_FILE) down

build-all: ## Construir todas las imágenes del proyecto
	@echo "🔨 Construyendo todas las imágenes..."
	docker compose -f $(COMPOSE_FILE) build

build-api: ## Construir solo imagen de API
	docker compose -f $(COMPOSE_FILE) build api

build-vue: ## Construir solo imagen de Vue
	docker compose -f $(COMPOSE_FILE) build vue

build-rasa: ## Construir solo imagen de Rasa
	docker compose -f $(COMPOSE_FILE) build rasa

# Logs por servicio
logs-api: ## Ver logs del API
	docker compose -f $(COMPOSE_FILE) logs -f api

logs-vue: ## Ver logs del Dashboard Vue
	docker compose -f $(COMPOSE_FILE) logs -f vue

logs-rasa: ## Ver logs de Rasa
	docker compose -f $(COMPOSE_FILE) logs -f rasa

logs-playwright: ## Ver logs de Playwright
	docker compose -f $(COMPOSE_FILE) logs -f playwright

logs-baileys: ## Ver logs de Baileys
	docker compose -f $(COMPOSE_FILE) logs -f baileys

logs-all: ## Ver logs de todos los servicios del proyecto
	docker compose -f $(COMPOSE_FILE) logs -f

logs-database: ## Ver logs de bases de datos
	@echo "🗄️  Logs de bases de datos..."
	docker compose logs -f postgres mongodb

logs-auth: ## Ver logs relacionados con autenticación
	@echo "🔐 Filtrando logs de autenticación..."
	docker compose logs -f api | grep -i "auth\|login\|admin\|user"


# Estado y salud de servicios
status: ## Verificar estado de servicios del proyecto
	@echo "📊 Estado de los servicios de $(PROJECT_NAME):"
	@docker compose -f $(COMPOSE_FILE) ps

health-check: ## Verificar salud de servicios del proyecto
	@echo "🏥 Verificando salud de servicios..."
	@echo -n "Postgres: " && (pg_isready -h localhost -p 5432 -U postgres 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "MongoDB: " && (curl -f http://localhost:27017 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "Rasa: " && (curl -f http://localhost:5005/status 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "Playwright: " && (curl -f http://localhost:3001/health 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "Vue Dashboard: " && (curl -f http://localhost:3002/health 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "API: " && (curl -f http://localhost:8080/health 2>/dev/null && echo "✅" || echo "❌")
	@echo -n "Baileys: " && (curl -f http://localhost:3000/health 2>/dev/null && echo "✅" || echo "❌")

# Limpieza
clean: clean-project ## Limpiar SOLO los contenedores, imágenes y volúmenes de este proyecto

clean-project: ## Limpiar contenedores, imágenes y volúmenes específicos del proyecto
	@echo "🧹 Limpiando recursos del proyecto $(PROJECT_NAME)..."
	@echo "⏹️  Deteniendo contenedores del proyecto..."
	-docker compose -f $(COMPOSE_FILE) down -v --remove-orphans 2>/dev/null
	@echo "🗑️  Eliminando contenedores del proyecto..."
	-docker container rm -f $(PROJECT_NAME)-postgres $(PROJECT_NAME)-mongo $(PROJECT_NAME)-rasa $(PROJECT_NAME)-playwright $(PROJECT_NAME)-vue $(PROJECT_NAME)-api $(PROJECT_NAME)-baileys 2>/dev/null || true
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

clean-auth: ## Limpiar datos de autenticación (sessions, tokens)
	@echo "🧹 Limpiando datos de autenticación..."
	@echo "⚠️  Esto cerrará todas las sesiones activas"
	@read -p "¿Continuar? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 0
	docker exec -it docubot-postgres psql -U ${POSTGRES_USER:-postgres} -d ${POSTGRES_DB:-docubot_db} -c "TRUNCATE TABLE system_users CASCADE;"
	@echo "✅ Datos de autenticación limpiados. Reinicia la API para crear el admin por defecto."

reset-database: ## Resetear completamente las bases de datos
	@echo "💀 PELIGRO: Esto eliminará TODOS los datos"
	@read -p "¿Estás seguro? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 0
	docker compose down -v
	docker volume rm docubot_postgres_data docubot_mongo_data 2>/dev/null || true
	@echo "✅ Bases de datos reseteadas"

# Reinicio de servicios
restart: ## Reiniciar todos los servicios del proyecto
	@echo "🔄 Reiniciando servicios del proyecto..."
	docker compose -f $(COMPOSE_FILE) restart

restart-api: ## Reiniciar solo la API
	docker compose -f $(COMPOSE_FILE) restart api

restart-vue: ## Reiniciar solo Vue Dashboard
	docker compose -f $(COMPOSE_FILE) restart vue

restart-rasa: ## Reiniciar solo Rasa
	docker compose -f $(COMPOSE_FILE) restart rasa

restart-baileys: ## Reiniciar solo Baileys
	docker compose -f $(COMPOSE_FILE) restart baileys

# Comandos de desarrollo
dev-logs: ## Ver logs en tiempo real de todos los servicios
	docker compose -f $(COMPOSE_FILE) logs -f --tail=100

dev-shell-api: ## Abrir shell en el contenedor de la API
	docker compose -f $(COMPOSE_FILE) exec api /bin/sh

dev-shell-vue: ## Abrir shell en el contenedor de Vue
	docker compose -f $(COMPOSE_FILE) exec vue /bin/sh

dev-shell-rasa: ## Abrir shell en el contenedor de Rasa
	docker compose -f $(COMPOSE_FILE) exec rasa /bin/bash

# Comandos específicos de Vue
vue-dev: ## Ejecutar Vue en modo desarrollo (local)
	@echo "🎨 Iniciando Vue en modo desarrollo..."
	cd vue-dashboard && npm run dev

vue-build: ## Compilar Vue para producción (local)
	@echo "🔨 Compilando Vue para producción..."
	cd vue-dashboard && npm run build

vue-install: ## Instalar dependencias de Vue (local)
	@echo "📦 Instalando dependencias de Vue..."
	cd vue-dashboard && npm install

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

# Comandos útiles para desarrollo
open-urls: ## Abrir todas las URLs del proyecto en el navegador
	@echo "🌐 Abriendo URLs del proyecto..."
	@which open >/dev/null && (open http://localhost:3002 && open http://localhost:8080/health && open http://localhost:5005/status) || echo "Comando 'open' no disponible. URLs: http://localhost:3002 http://localhost:8080/health http://localhost:5005/status"

check-ports: ## Verificar qué puertos están en uso
	@echo "🔍 Verificando puertos del proyecto..."
	@echo "Puerto 3002 (Vue):" && (lsof -i :3002 2>/dev/null || echo "  Libre")
	@echo "Puerto 8080 (API):" && (lsof -i :8080 2>/dev/null || echo "  Libre")
	@echo "Puerto 5005 (Rasa):" && (lsof -i :5005 2>/dev/null || echo "  Libre")
	@echo "Puerto 3001 (Playwright):" && (lsof -i :3001 2>/dev/null || echo "  Libre")
	@echo "Puerto 3000 (Baileys):" && (lsof -i :3000 2>/dev/null || echo "  Libre")

# ===== GESTIÓN DE USUARIOS ADMIN =====
create-admin: ## Crear usuario administrador manualmente
	@echo "🔧 Ejecutando script de creación de admin..."
	docker exec -it docubot-api /app/scripts/create_admin.sh

reset-admin: ## Resetear contraseña del administrador
	@echo "🔄 Ejecutando reset de credenciales de admin..."
	docker exec -it docubot-api sh -c "cd /app && go run ./scripts/reset-admin.go"

list-admins: ## Listar usuarios administradores
	@echo "📋 Listando usuarios administradores..."
	docker exec -it docubot-api sh -c "cd /app && go run -c 'database.ConnectPostgres(); db := database.GetDB(); var users []models.SystemUser; db.Where(\"role = ?\", \"admin\").Find(&users); for _, u := range users { fmt.Printf(\"ID: %d | Username: %s | Email: %s | Active: %t\\n\", u.ID, u.Username, u.Email, u.IsActive) }'"

show-admin-credentials: ## Mostrar credenciales por defecto del admin
	@echo "🔑 Credenciales por defecto del administrador:"
	@echo "   Username: ${ADMIN_USERNAME:-admin}"
	@echo "   Email: ${ADMIN_EMAIL:-admin@docubot.local}"
	@echo "   Password: ${ADMIN_PASSWORD:-DocubotAdmin123!}"
	@echo ""
	@echo "⚠️  Estas credenciales se usan solo si no existe un usuario admin en la BD"
	@echo "💡 Configura ADMIN_USERNAME, ADMIN_EMAIL, ADMIN_PASSWORD en .env para personalizar"

# ===== DESARROLLO CON SHELL DE CONTENEDORES =====
dev-shell-api: ## Abrir shell en contenedor API para desarrollo
	@echo "🐚 Abriendo shell en el contenedor API..."
	@echo "💡 Comandos útiles:"
	@echo "   - go run ./scripts/reset-admin.go"
	@echo "   - go run ./cmd/create_user.go"
	@echo ""
	docker exec -it docubot-api bash

dev-shell-postgres: ## Conectar a PostgreSQL
	@echo "🗄️  Conectando a PostgreSQL..."
	docker exec -it docubot-postgres psql -U ${POSTGRES_USER:-postgres} -d ${POSTGRES_DB:-docubot_db}
